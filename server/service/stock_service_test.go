package service_test

import (
	"fmt"
	"server/model"
	"server/repository"
	"server/service"
	"testing"
)

type StockHistory = model.StockHistory
type StockCollection = model.StockCollection

var stockRepo = repository.NewStockRepositoryDBMock()

func TestCreateStockCollection(t *testing.T) {
	expected := "Successfully created stock collection"

	t.Run("Create stock collection", func(t *testing.T) {
		stockCollection := StockCollection{
			StockImage: "test-image",
			Name:       "test",
			Sign:       "test",
			Price:      20,
			History:    []StockHistory{},
		}

		stockRepo.On(
			"CreateStock",
			stockCollection,
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.CreateStockCollection(stockCollection)

		fmt.Println(actual, err)
	})
}

func TestCreateStockOrder(t *testing.T) {}

func TestGetAllStockCollections(t *testing.T) {}

func TestGetTop10Stocks(t *testing.T) {}

func TestGetStockCollection(t *testing.T) {}

func TestGetFavoriteStock(t *testing.T) {}

func TestGetStockHistory(t *testing.T) {}

func TestSetStockPrice(t *testing.T) {}

func TestEditStockName(t *testing.T) {}

func TestEditStockSign(t *testing.T) {}

func TestDeleteStockCollection(t *testing.T) {}

