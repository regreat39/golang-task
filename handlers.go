package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var svc ScheduleService

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/schedule", createSchedule)
	r.Get("/schedules", listSchedules)
	r.Get("/schedule", getSchedule)
	r.Get("/next_takings", nextTakings)

	return r
}

func createSchedule(w http.ResponseWriter, req *http.Request) {
	var r CreateScheduleRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		http.Error(w, "неверный JSON", http.StatusBadRequest)
		return
	}

	if r.UserID == 0 {
		http.Error(w, "user_id обязателен", http.StatusBadRequest)
		return
	}

	if r.DrugName == "" {
		http.Error(w, "drug_name обязателен", http.StatusBadRequest)
		return
	}
	id, err := svc.Create(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func listSchedules(w http.ResponseWriter, req *http.Request) {
	userID, ok := getUserID(w, req)
	if !ok {
		return
	}
	ids, err := svc.List(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string][]int{"ids": ids})
}

func getSchedule(w http.ResponseWriter, req *http.Request) {
	userID, ok := getUserID(w, req)
	if !ok {
		return
	}
	scheduleID, err := strconv.Atoi(req.URL.Query().Get("schedule_id"))
	if err != nil {
		http.Error(w, "некорректный schedule_id", http.StatusBadRequest)
		return
	}

	res, err := svc.Get(userID, scheduleID)
	if err != nil {
		http.Error(w, "расписание не найдено", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func nextTakings(w http.ResponseWriter, req *http.Request) {
	userID, ok := getUserID(w, req)
	if !ok {
		return
	}
	res, err := svc.NextDose(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func getUserID(w http.ResponseWriter, req *http.Request) (int, bool) {
	userID, err := strconv.Atoi(req.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "некорректный user_id", http.StatusBadRequest)
		return 0, false
	}
	return userID, true
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
