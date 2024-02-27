package handler_test

import (
	"server/handler"
	"server/service"
	"testing"
)


func TestSignUp(t *testing.T) {
	t.Run("s", func(t *testing.T) {
		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		_ = handler.NewUserHandler(userService, stockService)
	})
}