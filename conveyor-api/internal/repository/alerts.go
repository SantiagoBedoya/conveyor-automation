package repository

import (
	"context"
	"fmt"

	"github.com/SantiagoBedoya/coveyor-api/internal/db"
	"github.com/SantiagoBedoya/coveyor-api/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type AlertsRepository interface {
	Create(ctx context.Context, alertType model.AlertType, triggerValue, threshold int) (*model.Alert, error)
	GetByID(ctx context.Context, id string) (*model.Alert, error)
	List(ctx context.Context, limit, offset int) ([]model.Alert, error)
	ListActive(ctx context.Context) ([]model.Alert, error)
	Resolve(ctx context.Context, id string) (*model.Alert, error)
}

type alertsRepo struct {
	q *db.Queries
}

func NewAlertsRepo(q *db.Queries) AlertsRepository {
	return &alertsRepo{q: q}
}

func (r *alertsRepo) Create(ctx context.Context, alertType model.AlertType, triggerValue, threshold int) (*model.Alert, error) {
	result, err := r.q.CreateAlert(ctx, db.CreateAlertParams{
		Type:         string(alertType),
		TriggerValue: int32(triggerValue),
		Threshold:    int32(threshold),
	})
	if err != nil {
		return nil, fmt.Errorf("create alert: %w", err)
	}
	return alertToDomain(result), nil
}

func (r *alertsRepo) GetByID(ctx context.Context, id string) (*model.Alert, error) {
	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := r.q.GetAlertByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get alert: %w", err)
	}
	return alertToDomain(result), nil
}

func (r *alertsRepo) List(ctx context.Context, limit, offset int) ([]model.Alert, error) {
	results, err := r.q.ListAlerts(ctx, db.ListAlertsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("list alerts: %w", err)
	}
	return alertsToDomain(results), nil
}

func (r *alertsRepo) ListActive(ctx context.Context) ([]model.Alert, error) {
	results, err := r.q.ListActiveAlerts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active alerts: %w", err)
	}
	return alertsToDomain(results), nil
}

func (r *alertsRepo) Resolve(ctx context.Context, id string) (*model.Alert, error) {
	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := r.q.ResolveAlert(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("resolve alert: %w", err)
	}
	return alertToDomain(result), nil
}

func alertToDomain(a db.Alert) *model.Alert {
	alert := &model.Alert{
		ID:           a.ID.String(),
		Timestamp:    a.Timestamp.Time,
		Type:         model.AlertType(a.Type),
		TriggerValue: int(a.TriggerValue),
		Threshold:    int(a.Threshold),
		Active:       a.Active,
	}
	if a.ResolvedAt.Valid {
		alert.ResolvedAt = &a.ResolvedAt.Time
	}
	return alert
}

func alertsToDomain(rows []db.Alert) []model.Alert {
	out := make([]model.Alert, len(rows))
	for i, a := range rows {
		out[i] = *alertToDomain(a)
	}
	return out
}
