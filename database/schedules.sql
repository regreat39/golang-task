CREATE TABLE IF NOT EXISTS schedules(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    drug_name VARCHAR(128) NOT NULL DEFAULT "",
    period_minutes INTEGER NOT NULL CHECK (period_minutes BETWEEN 60 AND 1440),
    course_days INTEGER,
    start_date CHAR(10) NOT NULL DEFAULT ""
);