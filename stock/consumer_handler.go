package main

import (
	"context"

	"example.com/oms/common/api/avro"
	pb "example.com/oms/common/api/protobuf"
	"example.com/oms/common/broker/consumer"
	"example.com/oms/common/broker/topic"
	"example.com/oms/common/schemaregistry/deserializer"
)

type ConsumerHandler struct {
	deserializer deserializer.Deserializer
	service      StockService
	Handlers     map[string]consumer.EventHandler
}

func NewConsumerHandler(deser deserializer.Deserializer, service StockService) *ConsumerHandler {
	handler := &ConsumerHandler{
		deserializer: deser,
		service:      service,
	}

	handler.Handlers = map[string]consumer.EventHandler{
		topic.MenuCreated: handler.AddStock,
	}

	return handler
}

func (h *ConsumerHandler) AddStock(msg consumer.EventMessage) error {
	event := avro.MenuCreatedEvent{}

	err := h.deserializer.Deserialize(msg.Topic, msg.Value, &event)
	if err != nil {
		return err
	}

	req := &pb.AddStockRequest{
		ItemId:   event.Id,
		Quantity: 0,
	}

	_, err = h.service.AddStock(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}
