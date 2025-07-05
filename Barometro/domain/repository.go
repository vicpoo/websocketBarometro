// repository.go
package domain

import "github.com/vicpoo/websocketBarometro/Barometro/domain/entities"

type BarometricRepository interface {
	Save(data entities.BarometricData) error
	GetAll() ([]entities.BarometricData, error)
}
