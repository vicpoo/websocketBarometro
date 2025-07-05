// controller.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/websocketBarometro/Barometro/application"
)

type BarometricController struct {
	barometricUseCase *application.BarometricUseCase
}

func NewBarometricController(barometricUseCase *application.BarometricUseCase) *BarometricController {
	return &BarometricController{
		barometricUseCase: barometricUseCase,
	}
}

func (bc *BarometricController) GetAllBarometricData(c *gin.Context) {
	data, err := bc.barometricUseCase.GetAllBarometricData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
