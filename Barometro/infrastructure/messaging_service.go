// messaging_service.go
package infrastructure

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
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

	// Declarar el exchange si no existe
	err = ch.ExchangeDeclare(
		"amq.topic", // nombre del exchange
		"topic",     // tipo
		true,        // durable
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
		"sensor_barometro", // nombre de la cola
		true,               // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ms.ch.QueueBind(
		q.Name,
		"sensor_baro",  // routing key
		"amq.topic",    // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ms.ch.Consume(
		q.Name,
		"",
		false, // autoAck
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Mensaje crudo recibido: %s", string(msg.Body))

			var data map[string]interface{}
			if err := json.Unmarshal(msg.Body, &data); err == nil {
				log.Printf("Datos barométricos parseados: %+v", data)
			} else {
				log.Printf("Error al parsear JSON: %v", err)
			}

			// Enviar a todos los clientes WebSocket conectados
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
