package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Razzle131/TPModel/internal/consts"
	"github.com/Razzle131/TPModel/internal/models"
	"github.com/Razzle131/TPModel/internal/repository/sensorRepo"
	"github.com/Razzle131/TPModel/internal/repository/tankRepo"
	"github.com/Razzle131/TPModel/internal/repository/valveRepo"
)

type Service struct {
	isTpActive bool
	tankRepo   tankRepo.TankRepo
	valveRepo  valveRepo.ValveRepo
	sensorRepo sensorRepo.SensorRepo
}

func NewService(tankRepo tankRepo.TankRepo, valveRepo valveRepo.ValveRepo, sensorRepo sensorRepo.SensorRepo) *Service {
	return &Service{
		isTpActive: false,
		tankRepo:   tankRepo,
		valveRepo:  valveRepo,
		sensorRepo: sensorRepo,
	}
}

func (s *Service) UpdateTpValues(ctx context.Context) {
	for {
		db1, _ := s.sensorRepo.GetSensorByName(ctx, consts.DB1)
		db2, _ := s.sensorRepo.GetSensorByName(ctx, consts.DB2)
		du0, _ := s.sensorRepo.GetSensorByName(ctx, consts.DU0)
		du1, _ := s.sensorRepo.GetSensorByName(ctx, consts.DU1)
		du2, _ := s.sensorRepo.GetSensorByName(ctx, consts.DU2)

		kl0, _ := s.valveRepo.GetValveByName(ctx, consts.Kl0)
		kl1, _ := s.valveRepo.GetValveByName(ctx, consts.Kl1)
		kl2, _ := s.valveRepo.GetValveByName(ctx, consts.Kl2)
		kl3, _ := s.valveRepo.GetValveByName(ctx, consts.Kl3)

		tank0, _ := s.tankRepo.GetTankByName(ctx, consts.Tank0)
		tank1, _ := s.tankRepo.GetTankByName(ctx, consts.Tank1)
		tank2, _ := s.tankRepo.GetTankByName(ctx, consts.Tank2)

		if kl0.IsOpen {
			tank0.CurVolume = max(tank0.CurVolume-kl0.Productivity, 0)
		}
		if kl1.IsOpen {
			transfer := 0
			if tank1.CurVolume-kl1.Productivity >= 0 {
				transfer = kl1.Productivity
			} else {
				transfer = tank1.CurVolume % kl1.Productivity
			}
			tank0.CurVolume = min(tank0.CurVolume+transfer, tank0.MaxVolume)
			tank1.CurVolume -= transfer
		}
		if kl2.IsOpen {
			transfer := 0
			if tank2.CurVolume-kl2.Productivity >= 0 {
				transfer = kl2.Productivity
			} else {
				transfer = tank2.CurVolume % kl2.Productivity
			}
			tank0.CurVolume = min(tank0.CurVolume+transfer, tank0.MaxVolume)
			tank2.CurVolume -= transfer
		}
		if kl3.IsOpen {
			tank0.CurVolume = max(tank0.CurVolume-kl3.Productivity, 0)
		}

		s.tankRepo.UpdateTank(ctx, tank0)
		s.tankRepo.UpdateTank(ctx, tank1)
		s.tankRepo.UpdateTank(ctx, tank2)

		du0.Value = tank0.CurVolume > 0
		du1.Value = float32(tank0.CurVolume)/float32(tank0.MaxVolume) >= 0.45
		du2.Value = float32(tank0.CurVolume)/float32(tank0.MaxVolume) >= 0.9
		db1.Value = tank1.CurVolume > 0
		db2.Value = tank2.CurVolume > 0

		s.sensorRepo.UpdateSensor(ctx, du0)
		s.sensorRepo.UpdateSensor(ctx, du1)
		s.sensorRepo.UpdateSensor(ctx, du2)
		s.sensorRepo.UpdateSensor(ctx, db1)
		s.sensorRepo.UpdateSensor(ctx, db2)

		select {
		case <-ctx.Done():
			return
		default:
			//slog.Debug("updating values...")
		}

		// slog.Debug(fmt.Sprint(kl1))
		// slog.Debug(fmt.Sprint(tank1))
		// slog.Debug(fmt.Sprint(tank0))

		time.Sleep(time.Millisecond * consts.Tick)
	}
}

func (s *Service) OpenValve(ctx context.Context, valveName string) error {
	valve, err := s.valveRepo.GetValveByName(ctx, valveName)
	if err != nil {
		slog.Error(fmt.Sprintf("service open valve: %s", err.Error()))
		return fmt.Errorf("service open valve: %s", err.Error())
	}

	valve.IsOpen = true

	_, err = s.valveRepo.UpdateValve(ctx, valve)
	if err != nil {
		slog.Error(fmt.Sprintf("service open valve: %s", err.Error()))
		return fmt.Errorf("service open valve: %s", err.Error())
	}

	return nil
}

func (s *Service) CloseValve(ctx context.Context, valveName string) error {
	valve, err := s.valveRepo.GetValveByName(ctx, valveName)
	if err != nil {
		slog.Error(fmt.Sprintf("service close valve: %s", err.Error()))
		return fmt.Errorf("service close valve: %s", err.Error())
	}

	valve.IsOpen = false

	_, err = s.valveRepo.UpdateValve(ctx, valve)
	if err != nil {
		slog.Error(fmt.Sprintf("service close valve: %s", err.Error()))
		return fmt.Errorf("service close valve: %s", err.Error())
	}

	return nil
}

func (s *Service) ToggleDRP(ctx context.Context) error {
	drp, err := s.sensorRepo.GetSensorByName(ctx, consts.DRP)
	if err != nil {
		slog.Error(fmt.Sprintf("service toggle DRP: %s", err.Error()))
		return fmt.Errorf("service toggle DRP: %s", err.Error())
	}

	drp.Value = !drp.Value

	_, err = s.sensorRepo.UpdateSensor(ctx, drp)
	if err != nil {
		slog.Error(fmt.Sprintf("service toggle DRP: %s", err.Error()))
		return fmt.Errorf("service toggle DRP: %s", err.Error())
	}

	return nil
}

func (s *Service) GetSensorByName(ctx context.Context, sensorName string) (models.Sensor, error) {
	sensor, err := s.sensorRepo.GetSensorByName(ctx, sensorName)
	if err != nil {
		slog.Error(fmt.Sprintf("service get sensor: %s", err.Error()))
		return models.Sensor{}, fmt.Errorf("service get sensor: %s", err.Error())
	}

	return sensor, nil
}

// only for visualisation
func (s *Service) GetTankByName(ctx context.Context, tankName string) (models.Tank, error) {
	tank, err := s.tankRepo.GetTankByName(ctx, tankName)
	if err != nil {
		slog.Error(fmt.Sprintf("service get tank: %s", err.Error()))
		return models.Tank{}, fmt.Errorf("service get tank: %s", err.Error())
	}

	return tank, nil
}
