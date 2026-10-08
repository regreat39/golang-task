package main

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testMoment = time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)

func newTestStorage(t *testing.T) (ScheduleStorage, *sql.DB) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)

	storage, err := NewScheduleStorage(db)
	require.NoError(t, err)

	return storage, db
}

func TestServiceCreate(t *testing.T) {
	storage, db := newTestStorage(t)
	defer db.Close()

	cfg := Config{NextTakingsMin: 60}
	svc := NewScheduleService(storage, cfg, func() time.Time {
		return testMoment
	})

	zeroDays := 0
	req := CreateScheduleRequest{
		UserID:     1,
		DrugName:   "Аспирин",
		Period:     "24h",
		CourseDays: &zeroDays,
	}

	id, err := svc.Create(req)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	got, err := svc.Get(1, id)
	require.NoError(t, err)

	require.True(t, got.Active, "срок 0 = бессрочный период")
	require.NotEmpty(t, got.Intervals)
}

func TestServiceNextDose(t *testing.T) {
	storage, db := newTestStorage(t)
	defer db.Close()

	cfg := Config{NextTakingsMin: 60}
	svc := NewScheduleService(storage, cfg, func() time.Time {
		return testMoment
	})

	days := 14
	req := CreateScheduleRequest{
		UserID:     1,
		DrugName:   "Аспирин",
		Period:     "2h",
		CourseDays: &days,
	}

	_, err := svc.Create(req)
	require.NoError(t, err)

	got, err := svc.NextDose(1)
	require.NoError(t, err)
	require.Len(t, got.Items, 0)
	require.NotNil(t, got.Items)

	svc.cfg.NextTakingsMin = 120

	got, err = svc.NextDose(1)
	require.NoError(t, err)
	require.Len(t, got.Items, 1)
	require.Equal(t, "10:00", got.Items[0].At)
}
