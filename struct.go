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

type CreateScheduleRequest struct {
	DrugName   string `json:"drug_name"`
	Period     string `json:"period"`
	CourseDays *int   `json:"course_days"`
}

type ListSchedule struct {
	Schedule
	Active    bool   `json:"active"`
	Intervals []Dose `json:"intervals"`
}

type NextTaking struct {
	ScheduleID int    `json:"schedule_id"`
	DrugName   string `json:"drug_name"`
	At         string `json:"at"`
}

type NextTakingResponse struct {
	From  string       `json:"from"`
	To    string       `json:"to"`
	Items []NextTaking `json:"items"`
}
