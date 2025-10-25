package models

type Valve struct {
	Name         string
	Productivity int
	IsOpen       bool
}

func NewValve(name string, productivity int) *Valve {
	return &Valve{
		Name:         name,
		Productivity: productivity,
		IsOpen:       false,
	}
}
