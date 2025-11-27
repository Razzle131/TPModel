package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Razzle131/TPModel/internal/repository/sensorRepo"
	"github.com/Razzle131/TPModel/internal/repository/tankRepo"
	"github.com/Razzle131/TPModel/internal/repository/valveRepo"
	"github.com/Razzle131/TPModel/internal/service"
)

const (
	valveNameKey  = "valve"
	sensorNameKey = "sensor"
	tankNameKey   = "tank"
)

type Server struct {
	service service.Service
}

func NewServer() *Server {
	tankRepo := tankRepo.NewCache()
	valveRepo := valveRepo.NewCache()
	sensorRepo := sensorRepo.NewCache()

	srv := Server{
		service: *service.NewService(tankRepo, valveRepo, sensorRepo),
	}

	go srv.service.UpdateTpValues(context.Background())

	return &srv
}

func (s *Server) PostOpenValve(w http.ResponseWriter, r *http.Request) {
	valveName := r.PathValue(valveNameKey)

	err := s.service.OpenValve(r.Context(), valveName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostCloseValve(w http.ResponseWriter, r *http.Request) {
	valveName := r.PathValue(valveNameKey)

	err := s.service.CloseValve(r.Context(), valveName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostToggleDRP(w http.ResponseWriter, r *http.Request) {
	err := s.service.ToggleDRP(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetSensor(w http.ResponseWriter, r *http.Request) {
	sensorName := r.PathValue(sensorNameKey)

	sensor, err := s.service.GetSensorByName(r.Context(), sensorName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, _ := json.Marshal(sensor)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *Server) GetTank(w http.ResponseWriter, r *http.Request) {
	tankName := r.PathValue(tankNameKey)

	tank, err := s.service.GetTankByName(r.Context(), tankName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, _ := json.Marshal(tank)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *Server) GetValve(w http.ResponseWriter, r *http.Request) {
	valveName := r.PathValue(valveNameKey)

	valve, err := s.service.GetValveByName(r.Context(), valveName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, _ := json.Marshal(valve)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
