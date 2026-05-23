package handler

import (
	"github.com/SantiagoBedoya/coveyor-api/internal/repository"
	"github.com/SantiagoBedoya/coveyor-api/internal/ws"
)

type Handler struct {
	ReadingsRepo repository.ReadingsRepository
	AlertsRepo   repository.AlertsRepository
	Hub          *ws.Hub
}

func New(readingsRepo repository.ReadingsRepository, alertsRepo repository.AlertsRepository, hub *ws.Hub) *Handler {
	return &Handler{
		ReadingsRepo: readingsRepo,
		AlertsRepo:   alertsRepo,
		Hub:          hub,
	}
}
