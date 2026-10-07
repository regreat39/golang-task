package main

import (
	"fmt"
	"time"
)

const (
	dayStart  = 60 * 8
	dayFinish = 60 * 22
)

func multipleOf15(minute int) int {
	k := minute / 15
	rem := minute % 15
	if rem != 0 {
		return (k + 1) * 15
	}
	return k * 15
}

func formatHHMM(minute int) string {
	return fmt.Sprintf("%02d:%02d", minute/60, minute%60)
}

func GenerateDoses(perMin int) []Dose {
	doses := make([]Dose, 0)
	for n := 0; ; n++ {
		at := dayStart + n*perMin
		mulAt := multipleOf15(at)
		if mulAt >= dayFinish {
			break
		}
		doses = append(doses, Dose{
			At:      formatHHMM(mulAt),
			Minutes: mulAt,
		})
	}
	return doses
}

func IsActive(s Schedule, day time.Time) bool {
	date := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	if date.Before(s.StartDate) {
		return false
	}
	if s.CourseDays == nil {
		return true
	}
	durCourseDays := *s.CourseDays
	endDate := s.StartDate.AddDate(0, 0, durCourseDays)
	return date.Before(endDate)
}
