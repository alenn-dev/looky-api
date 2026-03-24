package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

type OrderStatusEvent struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
}

func InitProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}

	producer = p
	return nil
}

func PublishOrderStatus(event OrderStatusEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "order_status",
		Value: sarama.StringEncoder(payload),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func CloseProducer() {
	if producer != nil {
		producer.Close()
	}
}
