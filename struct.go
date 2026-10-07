package main

import "time"

type Schedule struct {
	ID            int
	UserID        int
	DrugName      string
	PeriodMinutes int
	CourseDays    *int
	StartDate     time.Time
}

type Dose struct {
	At      string `json:"at"`
	Minutes int    `json:"-"`
}
