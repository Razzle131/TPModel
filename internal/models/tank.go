package models

type Tank struct {
	Name      string
	MaxVolume int
	CurVolume int
}

func NewTank(name string, capacity, initCapacity int) *Tank {
	return &Tank{
		Name:      name,
		MaxVolume: capacity,
		CurVolume: initCapacity,
	}
}
