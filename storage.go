package main

import (
	"database/sql"
	_ "embed"
	"time"

	_ "modernc.org/sqlite"
)

type ScheduleStorage struct {
	db *sql.DB
}

//go:embed database/schedules.sql
var schemaSQL string

func NewScheduleStorage(db *sql.DB) (ScheduleStorage, error) {
	if _, err := db.Exec(schemaSQL); err != nil {
		return ScheduleStorage{}, err
	}
	return ScheduleStorage{db: db}, nil
}

func (s ScheduleStorage) Create(sch Schedule) (int, error) {
	courseDays := sql.NullInt64{}
	if sch.CourseDays != nil {
		courseDays = sql.NullInt64{Int64: int64(*sch.CourseDays), Valid: true}
	}

	startDate := sch.StartDate.Format(time.DateOnly)

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	strAdd, err := tx.Exec(`INSERT INTO schedules (user_id, drug_name, period_minutes, course_days, start_date)
		VALUES (:user_id, :drug_name, :period_minutes, :course_days, :start_date)`,
		sql.Named("user_id", sch.UserID),
		sql.Named("drug_name", sch.DrugName),
		sql.Named("period_minutes", sch.PeriodMinutes),
		sql.Named("course_days", courseDays),
		sql.Named("start_date", startDate))
	if err != nil {
		return 0, err
	}

	id, err := strAdd.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ScheduleStorage) GetByUser(userID int) ([]int, error) {
	rows, err := s.db.Query(`SELECT id FROM schedules WHERE user_id = :user_id ORDER BY id`,
		sql.Named("user_id", userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]int, 0)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ScheduleStorage) GetByUserAndID(userID, scheduleID int) (Schedule, error) {
	row := s.db.QueryRow(`SELECT id, user_id, drug_name, period_minutes, course_days, start_date FROM schedules
		WHERE user_id = :user_id AND id = :schedule_id`,
		sql.Named("user_id", userID),
		sql.Named("schedule_id", scheduleID))

	sch := Schedule{}
	courseDays := sql.NullInt64{}
	var stDate string

	err := row.Scan(&sch.ID, &sch.UserID, &sch.DrugName, &sch.PeriodMinutes, &courseDays, &stDate)
	if err != nil {
		return sch, err
	}

	if courseDays.Valid {
		days := int(courseDays.Int64)
		sch.CourseDays = &days
	}

	startDate, err := time.Parse(time.DateOnly, stDate)
	if err != nil {
		return sch, err
	}
	sch.StartDate = startDate
	return sch, nil
}

func (s ScheduleStorage) FindByUser(userID int, day time.Time) ([]Schedule, error) {
	rows, err := s.db.Query(`SELECT id, user_id, drug_name, period_minutes, course_days, start_date
		FROM schedules WHERE user_id = :user_id AND start_date <= :day`,
		sql.Named("user_id", userID),
		sql.Named("day", day.Format(time.DateOnly)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]Schedule, 0)
	for rows.Next() {
		sch := Schedule{}
		courseDays := sql.NullInt64{}
		var stDate string

		err := rows.Scan(&sch.ID, &sch.UserID, &sch.DrugName, &sch.PeriodMinutes, &courseDays, &stDate)
		if err != nil {
			return nil, err
		}

		if courseDays.Valid {
			days := int(courseDays.Int64)
			sch.CourseDays = &days
		}

		startDate, err := time.Parse(time.DateOnly, stDate)
		if err != nil {
			return nil, err
		}
		sch.StartDate = startDate
		schedules = append(schedules, sch)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}
