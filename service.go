package main

import (
	"time"
)

type ScheduleService struct {
	storage ScheduleStorage
	cfg     Config
	now     func() time.Time
}

func normCourseDays(days *int) *int {
	if days == nil || *days <= 0 {
		return nil
	}
	return days
}

func NewScheduleService(storage ScheduleStorage, cfg Config, now func() time.Time) ScheduleService {
	return ScheduleService{
		storage: storage,
		cfg:     cfg,
		now:     now,
	}
}

func (s ScheduleService) Create(userID int, req CreateScheduleRequest) (int, error) {
	period, err := time.ParseDuration(req.Period)
	if err != nil {
		return 0, err
	}

	if period < time.Hour || period > 24*time.Hour {
		return 0, err
	}

	now := s.now().UTC()

	sch := Schedule{
		UserID:        userID,
		DrugName:      req.DrugName,
		PeriodMinutes: int(period.Minutes()),
		CourseDays:    normCourseDays(req.CourseDays),
		StartDate:     now,
	}

	return s.storage.Create(sch)
}

func (s ScheduleService) List(userID int) ([]int, error) {
	return s.storage.GetByUser(userID)
}

func (s ScheduleService) Get(userID, scheduleID int) (ListSchedule, error) {
	sch, err := s.storage.GetByUserAndID(userID, scheduleID)
	if err != nil {
		return ListSchedule{}, err
	}
	now := s.now().UTC()

	active := IsActive(sch, now)

	intervals := make([]Dose, 0)
	if active {
		intervals = GenerateDoses(sch.PeriodMinutes)
	}

	return ListSchedule{
		Schedule:  sch,
		Active:    active,
		Intervals: intervals,
	}, nil
}

func (s ScheduleService) NextDose(userID int) (NextTakingResponse, error) {
	now := s.now().UTC()

	nowMin := now.Hour()*60 + now.Minute()
	timer := s.cfg.NextTakingsMin

	schedules, err := s.storage.FindByUser(userID, now)
	if err != nil {
		return NextTakingResponse{}, err
	}

	items := make([]NextTaking, 0)
	for _, sch := range schedules {
		if !IsActive(sch, now) {
			continue
		}
		for _, dose := range GenerateDoses(sch.PeriodMinutes) {
			if dose.Minutes < nowMin || dose.Minutes >= nowMin+timer {
				continue
			}
			items = append(items, NextTaking{
				ScheduleID: sch.ID,
				DrugName:   sch.DrugName,
				At:         dose.At,
			})
		}
	}
	return NextTakingResponse{
		From:  now.Format(time.RFC3339),
		To:    now.Add(time.Duration(timer) * time.Minute).Format(time.RFC3339),
		Items: items,
	}, nil
}
