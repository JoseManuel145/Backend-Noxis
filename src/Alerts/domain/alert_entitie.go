package domain

type Alert struct {
	Sensor string         `json:"name"`
	Data   map[string]any `json:"data"`
}
