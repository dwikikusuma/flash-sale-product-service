package msgBroker

import (
	"context"
	"encoding/json"
	"product-catalog-service/infrastructure/log"
	"product-catalog-service/internal/entity"
	"product-catalog-service/internal/service"
	"strings"

	"github.com/segmentio/kafka-go"
)

type MsgConsumer struct {
	productSvc service.ProductService
}

func NewMsgConsumer(productSvc service.ProductService) *MsgConsumer {
	return &MsgConsumer{
		productSvc: productSvc,
	}
}

func (c *MsgConsumer) StartConsumer(brokers []string, topic string, groupID string) {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		ctx := context.Background()
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to read message from Kafka")
			continue
		}

		log.Logger.Info().Str("message", string(m.Value)).Msg("Received message from Kafka")
		c.processMessage(ctx, m)
	}
}

func (c *MsgConsumer) processMessage(ctx context.Context, msg kafka.Message) {
	var order *entity.Order
	err := json.Unmarshal(msg.Value, &order)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to unmarshal Kafka message")
		return
	}

	key := string(msg.Key)
	listKey := strings.Split(key, ".")
	event := listKey[1]

	switch event {
	case "created":
		for _, orderReq := range order.ProductRequests {
			isAvail, resvErr := c.productSvc.ReserveProductStock(ctx, orderReq.ProductID, int(orderReq.Quantity))
			if resvErr != nil || !isAvail {
				log.Logger.Error().Err(resvErr).Int64("productID", orderReq.ProductID).Msg("Failed to reserve product stock")
			} else {
				log.Logger.Info().Int64("productID", orderReq.ProductID).Msg("Successfully reserved product stock")
			}
		}

	case "cancelled":
		for _, orderReq := range order.ProductRequests {
			isReleased, relErr := c.productSvc.ReleaseProductStock(ctx, orderReq.ProductID, int(orderReq.Quantity))
			if relErr != nil || !isReleased {
				log.Logger.Error().Err(relErr).Int64("productID", orderReq.ProductID).Msg("Failed to release product stock")
			} else {
				log.Logger.Info().Int64("productID", orderReq.ProductID).Msg("Successfully released product stock")
			}
		}
	default:
		log.Logger.Warn().Str("event", event).Msg("Unknown event type")
	}
}
