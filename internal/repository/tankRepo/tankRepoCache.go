package tankRepo

import (
	"context"
	"sync"

	"github.com/Razzle131/TPModel/internal/consts"
	"github.com/Razzle131/TPModel/internal/models"
	"github.com/Razzle131/TPModel/internal/serverErrors"
)

type TankRepoCache struct {
	mu    sync.RWMutex
	tanks map[string]models.Tank
}

func NewCache() *TankRepoCache {
	repo := TankRepoCache{
		mu:    sync.RWMutex{},
		tanks: make(map[string]models.Tank),
	}

	tank0 := models.NewTank(consts.Tank0, consts.CapTank0, 0)
	tank1 := models.NewTank(consts.Tank1, consts.CapTank1, consts.CapTank1)
	tank2 := models.NewTank(consts.Tank2, consts.CapTank2, consts.CapTank2)

	repo.tanks[tank0.Name] = *tank0
	repo.tanks[tank1.Name] = *tank1
	repo.tanks[tank2.Name] = *tank2

	return &repo
}

func (r *TankRepoCache) GetTankByName(ctx context.Context, name string) (models.Tank, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tank, found := r.tanks[name]

	if !found {
		return models.Tank{}, serverErrors.ErrNotFound
	}

	return tank, nil
}

func (r *TankRepoCache) UpdateTank(ctx context.Context, tank models.Tank) (models.Tank, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, found := r.tanks[tank.Name]
	if !found {
		return models.Tank{}, serverErrors.ErrNotFound
	}

	r.tanks[tank.Name] = tank

	return r.tanks[tank.Name], nil
}
