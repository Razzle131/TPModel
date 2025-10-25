package models

type Sensor struct {
	Name  string
	Value bool
}

func NewSensor(name string) *Sensor {
	return &Sensor{
		Name:  name,
		Value: false,
	}
}
