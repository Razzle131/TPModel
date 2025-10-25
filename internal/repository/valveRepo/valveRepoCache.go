package valveRepo

import (
	"context"
	"sync"

	"github.com/Razzle131/TPModel/internal/consts"
	"github.com/Razzle131/TPModel/internal/models"
	"github.com/Razzle131/TPModel/internal/serverErrors"
)

type ValveRepoCache struct {
	mu     sync.RWMutex
	valves map[string]models.Valve
}

func NewCache() *ValveRepoCache {
	repo := ValveRepoCache{
		mu:     sync.RWMutex{},
		valves: make(map[string]models.Valve),
	}

	kl0 := models.NewValve(consts.Kl0, consts.ProdValve0)
	kl1 := models.NewValve(consts.Kl1, consts.ProdValve1)
	kl2 := models.NewValve(consts.Kl2, consts.ProdValve2)
	kl3 := models.NewValve(consts.Kl3, consts.ProdValve3)

	repo.valves[kl0.Name] = *kl0
	repo.valves[kl1.Name] = *kl1
	repo.valves[kl2.Name] = *kl2
	repo.valves[kl3.Name] = *kl3

	return &repo
}

func (r *ValveRepoCache) GetValveByName(ctx context.Context, name string) (models.Valve, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	valve, found := r.valves[name]

	if !found {
		return models.Valve{}, serverErrors.ErrNotFound
	}

	return valve, nil
}

func (r *ValveRepoCache) UpdateValve(ctx context.Context, valve models.Valve) (models.Valve, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, found := r.valves[valve.Name]
	if !found {
		return models.Valve{}, serverErrors.ErrNotFound
	}

	r.valves[valve.Name] = valve

	return r.valves[valve.Name], nil
}
