# Conveyor API

IoT API para monitoreo de banda transportadora con sensores de gas, humedad, distancia y conteo de objetos.

## Stack

- **Router:** [chi](https://github.com/go-chi/chi)
- **DB:** PostgreSQL via [pgx](https://github.com/jackc/pgx) + [sqlc](https://sqlc.dev)
- **WebSocket:** [gorilla/websocket](https://github.com/gorilla/websocket)
- **Migraciones:** [golang-migrate](https://github.com/golang-migrate/migrate)

## Requisitos

- Go 1.23+
- PostgreSQL
- [sqlc](https://sqlc.dev) (para regenerar código)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI

## Setup

```bash
cp .env-example .env   # editar credenciales
createdb conveyor
migrate -path migrations -database "$DATABASE_URL" up
go run ./cmd/api
```

## Variables de entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `DATABASE_URL` | Conexión a PostgreSQL | — |
| `API_KEY` | API key para ESP32 | — |
| `PORT` | Puerto del servidor | `8080` |

## Endpoints

### Readings

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/api/readings` | Crear lectura |
| `GET` | `/api/readings` | Listar (`?limit=50&offset=0`) |
| `GET` | `/api/readings/{id}` | Obtener por ID |

### Alertas

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/alerts` | Listar (`?limit=50&offset=0`) |
| `GET` | `/api/alerts/active` | Alertas activas |
| `POST` | `/api/alerts/{id}/resolve` | Resolver alerta |

### Sistema

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/status` | Estado actual del sistema |
| `GET` | `/api/ws` | WebSocket (broadcast en tiempo real) |

## Autenticación

Los endpoints del ESP32 deben incluir el header `X-API-Key`.

## Regenerar código sqlc

```bash
sqlc generate
```

## WebSocket

El hub emite mensajes JSON con el formato `{"type": "reading" | "alert" | "alert_resolved", "data": ...}`.
