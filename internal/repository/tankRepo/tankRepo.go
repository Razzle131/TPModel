package tankRepo

import (
	"context"

	"github.com/Razzle131/TPModel/internal/models"
)

type TankRepo interface {
	GetTankByName(ctx context.Context, name string) (models.Tank, error)
	UpdateTank(ctx context.Context, tank models.Tank) (models.Tank, error)
}
