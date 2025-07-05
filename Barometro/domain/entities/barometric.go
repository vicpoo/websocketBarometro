// barometric.go
package entities

import "time"

type BarometricData struct {
	ID                int       `json:"id"`
	Sensor            string    `json:"sensor"`
	Temperature       float64   `json:"temperatura"`
	Pressure          float64   `json:"presion"`
	Altitude          float64   `json:"altitud"`
	TemperatureUnit   string    `json:"unidad_temperatura"`
	PressureUnit      string    `json:"unidad_presion"`
	AltitudeUnit      string    `json:"unidad_altitud"`
	Timestamp         int64     `json:"timestamp"`
	Location          string    `json:"ubicacion"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewBarometricData(
	sensor string,
	temperature float64,
	pressure float64,
	altitude float64,
	tempUnit string,
	pressureUnit string,
	altitudeUnit string,
	timestamp int64,
	location string,
) *BarometricData {
	return &BarometricData{
		Sensor:           sensor,
		Temperature:      temperature,
		Pressure:         pressure,
		Altitude:         altitude,
		TemperatureUnit:  tempUnit,
		PressureUnit:     pressureUnit,
		AltitudeUnit:     altitudeUnit,
		Timestamp:        timestamp,
		Location:         location,
		CreatedAt:        time.Now(),
	}
}