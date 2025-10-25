package valveRepo

import (
	"context"

	"github.com/Razzle131/TPModel/internal/models"
)

type ValveRepo interface {
	GetValveByName(ctx context.Context, name string) (models.Valve, error)
	UpdateValve(ctx context.Context, valve models.Valve) (models.Valve, error)
}
