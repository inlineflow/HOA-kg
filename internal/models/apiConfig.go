package models

type APIConfig struct {
	Homes       []Home
	Appartments []Appartment
}

func NewConfig() *APIConfig {
	return &APIConfig{
		Homes: []Home{
			{ID: "1", Address: "Боконбаева 7"},
			{ID: "2", Address: "Усенбаева 44"},
		},
		Appartments: []Appartment{
			// {ID: "1", FlatNumber: 1},
			// {ID: "2", FlatNumber: 2},
		},
	}
}
