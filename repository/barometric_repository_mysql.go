// barometric_repository_mysql.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/vicpoo/websocketBarometro/Barometro/domain"
	"github.com/vicpoo/websocketBarometro/Barometro/domain/entities"
	"github.com/vicpoo/websocketBarometro/core"
)

type barometricRepositoryMySQL struct {
	db *sql.DB
}

func NewBarometricRepositoryMySQL() domain.BarometricRepository {
	return &barometricRepositoryMySQL{
		db: core.GetBD(),
	}
}

func (r *barometricRepositoryMySQL) Save(data entities.BarometricData) error {
	var sensorID int
	err := r.db.QueryRow("SELECT id FROM sensors WHERE name = ?", data.Sensor).Scan(&sensorID)
	if err != nil {
		return fmt.Errorf("no se encontró el sensor '%s': %v", data.Sensor, err)
	}

	_, err = r.db.Exec(`
		INSERT INTO sensor_readings (
			sensor_id, temperature, pressure, recorded_at
		) VALUES (?, ?, ?, FROM_UNIXTIME(?))`,
		sensorID, data.Temperature, data.Pressure, data.Timestamp)

	if err != nil {
		return fmt.Errorf("error al insertar en sensor_readings: %v", err)
	}

	return nil
}

func (r *barometricRepositoryMySQL) GetAll() ([]entities.BarometricData, error) {
	// Implementar si lo necesitas
	return nil, nil
}
