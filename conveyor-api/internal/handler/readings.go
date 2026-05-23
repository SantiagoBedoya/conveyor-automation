package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SantiagoBedoya/coveyor-api/internal/model"
)

type createReadingRequest struct {
	GasValue      int     `json:"gas_value"`
	HumidityValue int     `json:"humidity_value"`
	DistanceCm    float64 `json:"distance_cm"`
	ObjectCount   int     `json:"object_count"`
	BeltRunning   bool    `json:"belt_running"`
	FanOn         bool    `json:"fan_on"`
	BuzzerOn      bool    `json:"buzzer_on"`
	DoorAngle     int     `json:"door_angle"`
}

func (h *Handler) CreateReading(w http.ResponseWriter, r *http.Request) {
	var req createReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	sr := model.SensorReading{
		GasValue:      req.GasValue,
		HumidityValue: req.HumidityValue,
		DistanceCm:    req.DistanceCm,
		ObjectCount:   req.ObjectCount,
		BeltRunning:   req.BeltRunning,
		FanOn:         req.FanOn,
		BuzzerOn:      req.BuzzerOn,
		DoorAngle:     req.DoorAngle,
	}

	created, err := h.ReadingsRepo.Create(r.Context(), sr)
	if err != nil {
		log.Printf("create reading: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	alerts := evaluateThresholds(req)
	for _, a := range alerts {
		alert, err := h.AlertsRepo.Create(r.Context(), a.Type, a.TriggerValue, a.Threshold)
		if err != nil {
			log.Printf("create alert: %v", err)
			continue
		}
		msg, _ := json.Marshal(map[string]string{"type": "alert", "data": alert.ID})
		h.Hub.Broadcast(msg)
	}

	msg, _ := json.Marshal(map[string]interface{}{"type": "reading", "data": created})
	h.Hub.Broadcast(msg)

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetReading(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	reading, err := h.ReadingsRepo.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "reading not found"})
		return
	}
	writeJSON(w, http.StatusOK, reading)
}

func (h *Handler) ListReadings(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}
	readings, err := h.ReadingsRepo.List(r.Context(), limit, offset)
	if err != nil {
		log.Printf("list readings: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if readings == nil {
		readings = []model.SensorReading{}
	}
	writeJSON(w, http.StatusOK, readings)
}

func evaluateThresholds(req createReadingRequest) []createAlertRequest {
	var alerts []createAlertRequest

	if req.GasValue > 500 {
		alerts = append(alerts, createAlertRequest{
			Type:         model.AlertGas,
			TriggerValue: req.GasValue,
			Threshold:    500,
		})
	}
	if req.HumidityValue < 3000 {
		alerts = append(alerts, createAlertRequest{
			Type:         model.AlertHumidity,
			TriggerValue: req.HumidityValue,
			Threshold:    3000,
		})
	}
	return alerts
}
