// messaging_service.go
package infrastructure

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/vicpoo/websocketBarometro/Barometro/application"
	"github.com/vicpoo/websocketBarometro/Barometro/domain/entities"
	"github.com/vicpoo/websocketBarometro/repository"
)

type MessagingService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	hub  *Hub
}

func NewMessagingService(hub *Hub) *MessagingService {
	conn, err := amqp.Dial("amqp://reyhades:reyhades@44.219.123.4:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return nil
	}

	err = ch.ExchangeDeclare(
		"amq.topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
		return nil
	}

	return &MessagingService{
		conn: conn,
		ch:   ch,
		hub:  hub,
	}
}

func (ms *MessagingService) ConsumeBarometricMessages() error {
	q, err := ms.ch.QueueDeclare(
		"sensor_barometro",
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	err = ms.ch.QueueBind(
		q.Name,
		"sensor_baro",
		"amq.topic",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ms.ch.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	// Repositorio y usecase
	repo := repository.NewBarometricRepositoryMySQL()
	useCase := application.NewBarometricUseCase(repo)

	go func() {
		for msg := range msgs {
			log.Printf("Mensaje crudo recibido: %s", string(msg.Body))

			var payload struct {
				Sensor          string  `json:"sensor"`
				Temperatura     float64 `json:"temperatura"`
				Presion         float64 `json:"presion"`
				Altitud         float64 `json:"altitud"`
				UnidadTemp      string  `json:"unidad_temperatura"`
				UnidadPresion   string  `json:"unidad_presion"`
				UnidadAltitud   string  `json:"unidad_altitud"`
				Timestamp       int64   `json:"timestamp"`
				Ubicacion       string  `json:"ubicacion"`
			}

			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("Error al parsear JSON: %v", err)
				msg.Nack(false, false)
				continue
			}

			data := entities.NewBarometricData(
				payload.Sensor,
				payload.Temperatura,
				payload.Presion,
				payload.Altitud,
				payload.UnidadTemp,
				payload.UnidadPresion,
				payload.UnidadAltitud,
				payload.Timestamp,
				payload.Ubicacion,
			)

			if err := useCase.SaveBarometricData(*data); err != nil {
				log.Printf("Error al guardar en BD: %v", err)
			} else {
				log.Println("Datos barométricos guardados correctamente.")
			}

			ms.hub.broadcast <- msg.Body
			msg.Ack(false)
		}
	}()

	return nil
}

func (ms *MessagingService) Close() {
	if ms.ch != nil {
		ms.ch.Close()
	}
	if ms.conn != nil {
		ms.conn.Close()
	}
}
