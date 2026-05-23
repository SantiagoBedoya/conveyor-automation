package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SantiagoBedoya/coveyor-api/internal/model"
)

type createAlertRequest struct {
	Type         model.AlertType `json:"type"`
	TriggerValue int             `json:"trigger_value"`
	Threshold    int             `json:"threshold"`
}

func (h *Handler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}
	alerts, err := h.AlertsRepo.List(r.Context(), limit, offset)
	if err != nil {
		log.Printf("list alerts: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if alerts == nil {
		alerts = []model.Alert{}
	}
	writeJSON(w, http.StatusOK, alerts)
}

func (h *Handler) ListActiveAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.AlertsRepo.ListActive(r.Context())
	if err != nil {
		log.Printf("list active alerts: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if alerts == nil {
		alerts = []model.Alert{}
	}
	writeJSON(w, http.StatusOK, alerts)
}

func (h *Handler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	alert, err := h.AlertsRepo.Resolve(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "alert not found"})
		return
	}
	msg, _ := json.Marshal(map[string]interface{}{"type": "alert_resolved", "data": alert})
	h.Hub.Broadcast(msg)
	writeJSON(w, http.StatusOK, alert)
}
