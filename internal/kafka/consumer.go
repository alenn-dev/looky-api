package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"looky/internal/ws"

	"github.com/IBM/sarama"
)

func StartConsumer(brokers []string) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal("couldn't start kafka consumer:", err)
	}

	partition, err := consumer.ConsumePartition("order_status", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("couldn't consume partition:", err)
	}

	go func() {
		for msg := range partition.Messages() {
			var event OrderStatusEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("error parsing kafka message:", err)
				continue
			}

			notification, _ := json.Marshal(map[string]string{
				"order_id": event.OrderID,
				"status":   event.Status,
			})

			ws.H.Send(event.CustomerID, notification)
			fmt.Printf("notified customer %s: order %s is now %s\n", event.CustomerID, event.OrderID, event.Status)
		}
	}()

}
