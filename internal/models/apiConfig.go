package models

type APIConfig struct {
	Data []Home
}

func NewConfig() *APIConfig {
	return &APIConfig{
		Data: []Home{
			{ID: "1", Address: "Боконбаева 7"},
			{ID: "2", Address: "Усенбаева 44"},
		},
	}
}
