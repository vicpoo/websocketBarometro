// routes.go
package infrastructure

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, hub *Hub) {
	// Ruta WebSocket
	r.GET("/ws", hub.HandleWebSocket)
}
