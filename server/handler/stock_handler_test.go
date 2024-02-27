package handler_test

import (
	"server/handler"
	"server/service"
	"testing"
)

func TestCreateStockCollection(t *testing.T) {
	t.Run("s", func(t *testing.T) {
		stockService := service.NewStockServiceMock()
		stockHandler := handler.NewStockHandler(stockService)
	})
}