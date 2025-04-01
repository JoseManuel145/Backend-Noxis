package application

import "fmt"

type Threshold struct {
	Min float64
	Max float64
}

var sensorThresholds = map[string]map[string]struct{ Min, Max float64 }{
	"Calidad Aire MQ-135": {
		"calidad_aire": {Min: 0, Max: 100},
	},
	"Carbono CJMCU-811": {
		"co2": {Min: 400, Max: 1500},
	},
	"Carbono MQ-7": {
		"monoxido": {Min: 0, Max: 50},
	},
	"Flama KY-026": {
		"flama": {Min: 0, Max: 1}, // Se puede detectar o no (0 o 1)
	},
	"Gas Natural MQ-5": {
		"gas_natural": {Min: 0, Max: 300},
	},
	"Hidrogeno MQ-136": {
		"sulfuro_hidrogeno": {Min: 0, Max: 10},
	},
	"Metano MQ-4": {
		"metano": {Min: 0, Max: 200},
	},
	"BME-680": {
		"temperatura": {Min: -10, Max: 50},
		"presion":     {Min: 900, Max: 1100},
		"humedad":     {Min: 10, Max: 90},
	},
}

func EsPeligroso(sensor string, data map[string]any) bool {
	thresholds, exists := sensorThresholds[sensor]
	if !exists {
		fmt.Println("Sensor desconocido, se ignora:", sensor)
		return false
	}
	for key, value := range data {
		val, ok := value.(float64)
		if !ok {
			continue
		}
		if t, ok := thresholds[key]; ok {
			if val < t.Min || val > t.Max {
				return true
			}
		}
	}
	return false
}
