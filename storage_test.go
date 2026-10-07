package main

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func getTestSchedule(userID int) Schedule {
	days := 14
	return Schedule{
		UserID:        userID,
		DrugName:      "Аспирин",
		PeriodMinutes: 120,
		CourseDays:    &days,
		StartDate:     time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC),
	}
}

func TestNewScheduleStorage(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	storage, err := NewScheduleStorage(db)
	require.NoError(t, err)

	idSch, err := storage.GetByUser(1)
	require.NoError(t, err)
	require.NotNil(t, idSch)

	_, err = NewScheduleStorage(db)
	require.NoError(t, err)
}

func TestCreateReturnID(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	storage, err := NewScheduleStorage(db)
	require.NoError(t, err)

	schedule := getTestSchedule(1)

	id, err := storage.Create(schedule)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestCreateWithoutCourseDays(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.go")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	storage, err := NewScheduleStorage(db)
	require.NoError(t, err)

	schedule := getTestSchedule(1)
	schedule.CourseDays = nil

	id, err := storage.Create(schedule)
	require.NoError(t, err)

	got, err := storage.GetByUserAndID(1, id)
	require.NoError(t, err)
	require.Nil(t, got.CourseDays)
}

func TestCreateIncorrectPeriod(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.go")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	storage, err := NewScheduleStorage(db)
	require.NoError(t, err)

	schedule := getTestSchedule(1)
	schedule.PeriodMinutes = 30

	_, err = storage.Create(schedule)
	require.Error(t, err)

	idSch, err := storage.GetByUser(1)
	require.NoError(t, err)
	require.Empty(t, idSch)
}

func TestGetByUserAndIDForeignUser(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.go")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	storage, err := NewScheduleStorage(db)
	require.NoError(t, err)

	id, err := storage.Create(getTestSchedule(1))
	require.NoError(t, err)

	_, err = storage.GetByUserAndID(2, id)
	require.ErrorIs(t, err, sql.ErrNoRows)
}
