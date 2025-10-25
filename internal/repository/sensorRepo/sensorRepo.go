package sensorRepo

import (
	"context"

	"github.com/Razzle131/TPModel/internal/models"
)

type SensorRepo interface {
	GetSensorByName(ctx context.Context, name string) (models.Sensor, error)
	UpdateSensor(ctx context.Context, sensor models.Sensor) (models.Sensor, error)
}
