package sensorRepo

import (
	"context"
	"sync"

	"github.com/Razzle131/TPModel/internal/consts"
	"github.com/Razzle131/TPModel/internal/models"
	"github.com/Razzle131/TPModel/internal/serverErrors"
)

type SensorRepoCache struct {
	mu      sync.RWMutex
	sensors map[string]models.Sensor
}

func NewCache() *SensorRepoCache {
	repo := SensorRepoCache{
		mu:      sync.RWMutex{},
		sensors: make(map[string]models.Sensor),
	}

	db1 := models.NewSensor(consts.DB1)
	db2 := models.NewSensor(consts.DB2)

	du0 := models.NewSensor(consts.DU0)
	du1 := models.NewSensor(consts.DU1)
	du2 := models.NewSensor(consts.DU2)

	drp := models.NewSensor(consts.DRP)

	repo.sensors[db1.Name] = *db1
	repo.sensors[db2.Name] = *db2

	repo.sensors[du0.Name] = *du0
	repo.sensors[du1.Name] = *du1
	repo.sensors[du2.Name] = *du2

	repo.sensors[drp.Name] = *drp

	return &repo
}

func (r *SensorRepoCache) GetSensorByName(ctx context.Context, name string) (models.Sensor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sensor, found := r.sensors[name]

	if !found {
		return models.Sensor{}, serverErrors.ErrNotFound
	}

	return sensor, nil
}

func (r *SensorRepoCache) UpdateSensor(ctx context.Context, sensor models.Sensor) (models.Sensor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, found := r.sensors[sensor.Name]
	if !found {
		return models.Sensor{}, serverErrors.ErrNotFound
	}

	r.sensors[sensor.Name] = sensor

	return r.sensors[sensor.Name], nil
}
