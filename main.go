package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := NewFromEnv()
	if err != nil {
		log.Fatalf("ошибка конфигурации: %v", err)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("не удалось открыть БД %s: %v", cfg.DBPath, err)
	}
	defer db.Close()

	storage, err := NewScheduleStorage(db)
	if err != nil {
		log.Fatalf("не удалось создать схему: %v", err)
	}

	svc = NewScheduleService(storage, cfg, time.Now)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      NewRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Printf("сервер запущен, текущий адрес %s, связанная БД %s", server.Addr, cfg.DBPath)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("сервер остановился с ошибкой: %v", err)
	}
}
