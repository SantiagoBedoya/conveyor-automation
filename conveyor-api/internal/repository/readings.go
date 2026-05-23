package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/SantiagoBedoya/coveyor-api/internal/db"
	"github.com/SantiagoBedoya/coveyor-api/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReadingsRepository interface {
	Create(ctx context.Context, r model.SensorReading) (*model.SensorReading, error)
	GetByID(ctx context.Context, id string) (*model.SensorReading, error)
	List(ctx context.Context, limit, offset int) ([]model.SensorReading, error)
	ListSince(ctx context.Context, since time.Time) ([]model.SensorReading, error)
	GetLatest(ctx context.Context) (*model.SensorReading, error)
}

type readingsRepo struct {
	q *db.Queries
}

func NewReadingsRepo(q *db.Queries) ReadingsRepository {
	return &readingsRepo{q: q}
}

func (r *readingsRepo) Create(ctx context.Context, sr model.SensorReading) (*model.SensorReading, error) {
	result, err := r.q.CreateReading(ctx, db.CreateReadingParams{
		GasValue:      int32(sr.GasValue),
		HumidityValue: int32(sr.HumidityValue),
		DistanceCm:    sr.DistanceCm,
		ObjectCount:   int32(sr.ObjectCount),
		BeltRunning:   sr.BeltRunning,
		FanOn:         sr.FanOn,
		BuzzerOn:      sr.BuzzerOn,
		DoorAngle:     int32(sr.DoorAngle),
	})
	if err != nil {
		return nil, fmt.Errorf("create reading: %w", err)
	}
	return readingToDomain(result), nil
}

func (r *readingsRepo) GetByID(ctx context.Context, id string) (*model.SensorReading, error) {
	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := r.q.GetReadingByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get reading: %w", err)
	}
	return readingToDomain(result), nil
}

func (r *readingsRepo) List(ctx context.Context, limit, offset int) ([]model.SensorReading, error) {
	results, err := r.q.ListReadings(ctx, db.ListReadingsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("list readings: %w", err)
	}
	return readingsToDomain(results), nil
}

func (r *readingsRepo) ListSince(ctx context.Context, since time.Time) ([]model.SensorReading, error) {
	var ts pgtype.Timestamptz
	if err := ts.Scan(since); err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}
	results, err := r.q.ListReadingsSince(ctx, ts)
	if err != nil {
		return nil, fmt.Errorf("list readings since: %w", err)
	}
	return readingsToDomain(results), nil
}

func (r *readingsRepo) GetLatest(ctx context.Context) (*model.SensorReading, error) {
	result, err := r.q.GetLatestReading(ctx)
	if err != nil {
		return nil, fmt.Errorf("get latest reading: %w", err)
	}
	return readingToDomain(result), nil
}

func readingToDomain(r db.SensorReading) *model.SensorReading {
	return &model.SensorReading{
		ID:            r.ID.String(),
		Timestamp:     r.Timestamp.Time,
		GasValue:      int(r.GasValue),
		HumidityValue: int(r.HumidityValue),
		DistanceCm:    r.DistanceCm,
		ObjectCount:   int(r.ObjectCount),
		BeltRunning:   r.BeltRunning,
		FanOn:         r.FanOn,
		BuzzerOn:      r.BuzzerOn,
		DoorAngle:     int(r.DoorAngle),
	}
}

func readingsToDomain(rows []db.SensorReading) []model.SensorReading {
	out := make([]model.SensorReading, len(rows))
	for i, r := range rows {
		out[i] = *readingToDomain(r)
	}
	return out
}
