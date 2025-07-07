// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vicpoo/websocketBarometro/core"
	"github.com/vicpoo/websocketBarometro/Barometro/infrastructure"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// ✅ Inicializar conexión a la base de datos
	core.InitDB()

	r := gin.Default()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Inicializar WebSocket hub
	hub := infrastructure.NewHub()
	go hub.Run()

	// Inicializar servicio de mensajería (RabbitMQ)
	messagingService := infrastructure.NewMessagingService(hub)
	defer messagingService.Close()

	// Rutas para WebSocket
	infrastructure.SetupRoutes(r, hub)

	// Iniciar consumidor de RabbitMQ
	if err := messagingService.ConsumeBarometricMessages(); err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}

	// Señal para cerrar app
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":8002"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started on port 8002")
	log.Println("Barometric consumer started")

	<-sigChan
	log.Println("Shutting down server...")
}
