package service

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventProducer struct {
	writer *kafka.Writer
}

func NewEventProducer(brokers []string) *EventProducer {
	if len(brokers) == 0 || brokers[0] == "" {
		return &EventProducer{writer: nil}
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 3 * time.Second,
		RequiredAcks: kafka.RequireOne,
	}

	return &EventProducer{writer: w}
}

func (p *EventProducer) Send(ctx context.Context, topic, message string) error {
	if p.writer == nil {
		log.Printf("[EventProducer] Kafka disabled or unavailable, skipping send to %s: %s", topic, message)
		return nil
	}

	msg := kafka.Message{
		Topic: topic,
		Value: []byte(message),
		Time:  time.Now(),
	}

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("[EventProducer] Warning: Failed to send Kafka message to topic %s: %v", topic, err)
		return err
	}
	log.Printf("[EventProducer] Sent Kafka message to topic %s successfully", topic)
	return nil
}

func (p *EventProducer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
