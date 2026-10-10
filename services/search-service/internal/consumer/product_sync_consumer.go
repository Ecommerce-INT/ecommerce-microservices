package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"com.ecommerce/search-service/internal/model"
	"com.ecommerce/search-service/internal/service"
	"github.com/segmentio/kafka-go"
)

type ProductSyncConsumer struct {
	svc     *service.SearchService
	brokers []string
	topic   string
	groupID string
}

func NewProductSyncConsumer(svc *service.SearchService, brokers []string, topic, groupID string) *ProductSyncConsumer {
	return &ProductSyncConsumer{
		svc:     svc,
		brokers: brokers,
		topic:   topic,
		groupID: groupID,
	}
}

func (c *ProductSyncConsumer) Start(ctx context.Context) {
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  c.brokers,
			Topic:    c.topic,
			GroupID:  c.groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
		defer r.Close()

		log.Printf("[Search Consumer] Subscribed to topic '%s' with group '%s'", c.topic, c.groupID)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, err := r.ReadMessage(ctx)
				if err != nil {
					time.Sleep(2 * time.Second)
					continue
				}

				var payload map[string]any
				if err := json.Unmarshal(m.Value, &payload); err != nil {
					continue
				}

				op, _ := payload["op"].(string)
				after, _ := payload["after"].(map[string]any)

				if op == "d" {
					before, _ := payload["before"].(map[string]any)
					if before != nil {
						if idVal, ok := before["id"].(float64); ok {
							_ = c.svc.DeleteProduct(ctx, int64(idVal))
						}
					}
				} else if after != nil {
					idVal, _ := after["id"].(float64)
					nameVal, _ := after["name"].(string)
					slugVal, _ := after["slug"].(string)
					priceVal, _ := after["price"].(float64)

					p := model.ProductDoc{
						ID:        int64(idVal),
						Name:      nameVal,
						Slug:      slugVal,
						Price:     priceVal,
						CreatedOn: time.Now(),
					}
					_ = c.svc.IndexProduct(ctx, p)
				}
			}
		}
	}()
}
