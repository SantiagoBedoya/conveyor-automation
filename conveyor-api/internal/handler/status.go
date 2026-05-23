package handler

import (
	"log"
	"net/http"

	"github.com/SantiagoBedoya/coveyor-api/internal/model"
)

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	latest, err := h.ReadingsRepo.GetLatest(r.Context())
	if err != nil {
		log.Printf("get status: %v", err)
		writeJSON(w, http.StatusOK, model.SystemStatus{
			ActiveAlerts: []model.Alert{},
		})
		return
	}

	alerts, err := h.AlertsRepo.ListActive(r.Context())
	if err != nil {
		log.Printf("list active alerts for status: %v", err)
		alerts = []model.Alert{}
	}
	if alerts == nil {
		alerts = []model.Alert{}
	}

	status := model.SystemStatus{
		LastSeen:     latest.Timestamp,
		BeltRunning:  latest.BeltRunning,
		FanOn:        latest.FanOn,
		BuzzerOn:     latest.BuzzerOn,
		DoorAngle:    latest.DoorAngle,
		ActiveAlerts: alerts,
	}
	writeJSON(w, http.StatusOK, status)
}
