package main

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMultipleOf15(t *testing.T) {
	table := []struct {
		inputMinute, outputMinute int
	}{
		{0, 0},
		{1, 15},
		{15, 15},
		{480, 480},
		{487, 495},
		{510, 510},
		{1319, 1320},
		{1321, 1335},
	}

	for _, val := range table {
		multVal := multipleOf15(val.inputMinute)
		require.Equal(t, val.outputMinute, multVal)
	}
}

func TestFormatHHMM(t *testing.T) {
	table := []struct {
		inputMinute int
		corrTime    string
	}{
		{0, "00:00"},
		{15, "00:15"},
		{480, "08:00"},
		{570, "09:30"},
		{1319, "21:59"},
		{1320, "22:00"},
		{1321, "22:01"},
	}

	for _, val := range table {
		timeStr := formatHHMM(val.inputMinute)
		require.Equal(t, val.corrTime, timeStr)
	}
}

func TestGenerateDoses(t *testing.T) {
	table := []struct {
		name    string
		perMin  int
		perTime string
	}{
		{"раз в сутки", 1440, "08:00"},
		{"интервал больше допустимых часов приема", 1000, "08:00"},
		{"каждые 8 часов", 480, "08:00 16:00"},
		{"каждые 4 часа", 240, "08:00 12:00 16:00 20:00"},
		{"ежечасно", 60, "08:00 09:00 10:00 11:00 12:00 13:00 14:00 15:00 16:00 17:00 18:00 19:00 20:00 21:00"},
		{"период не кратен 15", 100, "08:00 09:45 11:30 13:00 14:45 16:30 18:00 19:45 21:30"},
	}
	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			doseTime := GenerateDoses(v.perMin)

			got := make([]string, len(doseTime))
			for i, val := range doseTime {
				got[i] = val.At
			}
			want := strings.Split(v.perTime, " ")
			require.Equal(t, want, got, "период %d", v.perMin)
		})
	}
}

func TestIsActive(t *testing.T) {
	days := 14
	start, err := time.Parse("2006-01-02", "2026-10-05")
	require.NoError(t, err)

	finish := Schedule{StartDate: start, CourseDays: &days}
	permanent := Schedule{StartDate: start, CourseDays: nil}

	table := []struct {
		name string
		s    Schedule
		day  string
		want bool
	}{
		{"до начала курса", finish, "2026-10-01", false},
		{"день начала", finish, "2026-10-05", true},
		{"один из дней курса", finish, "2026-10-10", true},
		{"конец курса", finish, "2026-10-18", true},
		{"после конца курса", finish, "2026-10-19", false},
		{"бессрочный прием, до начала", permanent, "2026-10-01", false},
		{"бессрочный прием, день начала", permanent, "2026-10-05", true},
		{"бессрочный прием, в процессе", permanent, "2026-11-11", true},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			day, err := time.Parse("2006-01-02", v.day)
			require.NoError(t, err)
			result := IsActive(v.s, day)
			require.Equal(t, v.want, result)
		})
	}
}
