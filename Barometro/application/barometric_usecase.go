// barometric_usecase.go
package application

import (
	"github.com/vicpoo/websocketBarometro/Barometro/domain"
	"github.com/vicpoo/websocketBarometro/Barometro/domain/entities"
)

type BarometricUseCase struct {
	repo domain.BarometricRepository
}

func NewBarometricUseCase(repo domain.BarometricRepository) *BarometricUseCase {
	return &BarometricUseCase{repo: repo}
}

func (uc *BarometricUseCase) SaveBarometricData(data entities.BarometricData) error {
	return uc.repo.Save(data)
}

func (uc *BarometricUseCase) GetAllBarometricData() ([]entities.BarometricData, error) {
	return uc.repo.GetAll()
}
