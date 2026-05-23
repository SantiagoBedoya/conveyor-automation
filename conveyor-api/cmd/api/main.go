package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/SantiagoBedoya/coveyor-api/internal/db"
	"github.com/SantiagoBedoya/coveyor-api/internal/handler"
	"github.com/SantiagoBedoya/coveyor-api/internal/repository"
	"github.com/SantiagoBedoya/coveyor-api/internal/ws"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://localhost:5432/conveyor?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}
	log.Println("connected to database")

	q := db.New(pool)

	hub := ws.NewHub()
	go hub.Run()

	h := handler.New(
		repository.NewReadingsRepo(q),
		repository.NewAlertsRepo(q),
		hub,
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Route("/api", func(r chi.Router) {
		r.Get("/ws", h.ServeWS)

		r.Route("/readings", func(r chi.Router) {
			r.Post("/", h.CreateReading)
			r.Get("/", h.ListReadings)
			r.Get("/{id}", h.GetReading)
		})

		r.Route("/alerts", func(r chi.Router) {
			r.Get("/", h.ListAlerts)
			r.Get("/active", h.ListActiveAlerts)
			r.Post("/{id}/resolve", h.ResolveAlert)
		})

		r.Get("/status", h.GetStatus)
	})

	fs := http.FileServer(http.Dir("./dashboard"))
	r.Handle("/*", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
