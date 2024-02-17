package service_test

import (
	"server/errs"
	"server/model"
	"server/repository"
	"server/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StockHistory = model.StockHistory
type StockCollection = model.StockCollection
type StockCollectionResponse = model.StockCollectionResponse
type StockGroup = model.StockGroup
type TopStock = model.TopStock
type StockHistoryResponse = model.StockHistoryResponse

var stockRepo = repository.NewStockRepositoryDBMock()
var (
	ErrPrice = errs.ErrPrice
	ErrName = errs.ErrName
	ErrSign = errs.ErrSign
)


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

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid data", func(t *testing.T) {
		stockCollection := StockCollection{
			StockImage: "",
			Name:       "",
			Sign:       "",
			Price:      20,
			History:    []StockHistory{},
		}

		stockRepo.On(
			"CreateStock",
			stockCollection,
		).Return(expected, ErrData)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.CreateStockCollection(stockCollection)

		assert.ErrorIs(t, err, ErrData)
	})
}

func TestCreateStockOrder(t *testing.T) {
	expected := "Successfully created stock order"

	t.Run("Create stock order", func(t *testing.T) {
		stockOrder := StockHistory{
			ID: "65c8993c48096b5150cee5d6",
			Timestamp: int64(1),
			Amount: float64(1),
			Price: float64(1),
		}

		stockRepo.On(
			"CreateStockOrder",
			"65cc5fd45aa71b64fbb551a9",
			stockOrder,
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.CreateStockOrder(
			"65cc5fd45aa71b64fbb551a9",
			stockOrder,
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)		
	})

	t.Run("Error invalid data", func(t *testing.T) {
		stockOrder := StockHistory{
			ID: "",
			Timestamp: int64(1),
			Amount: float64(1),
			Price: float64(1),
		}

		stockRepo.On(
			"CreateStockOrder",
			"65cc5fd45aa71b64fbb551a9",
			stockOrder,
		).Return(expected, ErrData)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.CreateStockOrder(
			"65cc5fd45aa71b64fbb551a9",
			stockOrder,
		)

		assert.ErrorIs(t, err, ErrData)
	})
}

func TestGetAllStockCollections(t *testing.T) {
	expected := []StockCollectionResponse{{
		ID: "1",
		StockImage: "test",
		Name: "test",
		Sign: "t",
		Price: 1,
	}}
	
	t.Run("Get all stock collections", func(t *testing.T) {
		stockRepo.On(
			"GetAllStocks",
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.GetAllStockCollections()

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetTop10Stocks(t *testing.T) {
	expected := []TopStock{{
		ID: "1",
		Sign: "t",
		Price: 1,
	}}

	t.Run("Get top 10 stocks", func(t *testing.T) {
		stockRepo.On(
			"GetTopStocks",
		).Return([]StockGroup{{
			ID: "1",
			StockImage: "test",
			Name: "test",
			Sign: "t",
			Price: 1,
			Volume: 1,
		}}, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.GetTop10Stocks()
		 
		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetStockCollection(t *testing.T) {
	expected := StockCollectionResponse{
		ID: "1",
		StockImage: "test",
		Name: "test",
		Sign: "t",
		Price: 1,
	}

	t.Run("Get stock collection", func(t *testing.T) {
		stockRepo.On(
			"GetStock",
			"65cc5fd45aa71b64fbb551a9",
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.GetStockCollection("65cc5fd45aa71b64fbb551a9")

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		stockRepo.On(
			"GetStock",
			"",
		).Return(expected, ErrInvalidStock)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.GetStockCollection("")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})
}

func TestGetFavoriteStock(t *testing.T) {
	expected := []StockCollectionResponse{{
		ID: "65cc5fd45aa71b64fbb551a9",
		StockImage: "test",
		Name: "test",
		Sign: "t",
		Price: 1,
	}}
	

	t.Run("Get favorite stock", func(t *testing.T) {
		stockIds := []string{"65cc5fd45aa71b64fbb551a9"}
		stockRepo.On(
			"GetFavoriteStock",
			stockIds,
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.GetFavoriteStock(stockIds)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		stockRepo.On(
			"GetFavoriteStock",
			[]string{""},
		).Return(expected, ErrInvalidStock)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.GetFavoriteStock([]string{""})

		assert.ErrorIs(t, err, ErrInvalidStock)
	})
}

func TestGetStockHistory(t *testing.T) {
	expected := []StockHistoryResponse{{
		Amount: float64(1),
		Price: float64(1),
	}}

	t.Run("Get stock history", func(t *testing.T) {
		stockRepo.On(
			"GetStockHistory", 
			"65cc5fd45aa71b64fbb551a9",
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.GetStockHistory("65cc5fd45aa71b64fbb551a9")

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		stockRepo.On(
			"GetStockHistory", 
			"",
		).Return(expected, ErrInvalidStock)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.GetStockHistory("")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})
}

func TestSetStockPrice(t *testing.T) {
	expected := "Successfully set stock price"

	t.Run("Set stock price", func(t *testing.T) {
		stockRepo.On(
			"SetPrice",
			"65cc5fd45aa71b64fbb551a9",
			float64(1),
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.SetStockPrice(
			"65cc5fd45aa71b64fbb551a9", 
			float64(1),
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid price", func(t *testing.T) {
		stockRepo.On(
			"SetPrice",
			"65cc5fd45aa71b64fbb551a9",
			float64(0),
		).Return(expected, ErrPrice)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.SetStockPrice(
			"65cc5fd45aa71b64fbb551a9", 
			float64(0),
		)

		assert.ErrorIs(t, err, ErrPrice)
	})
}

func TestEditStockName(t *testing.T) {
	expected := "Successfully edit stock name"

	t.Run("Edit stock name", func(t *testing.T) {
		stockRepo.On(
			"EditName",
			"65cc5fd45aa71b64fbb551a9",
			"T",
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.EditStockName(
			"65cc5fd45aa71b64fbb551a9", 
			"T",
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid name", func(t *testing.T) {
		stockRepo.On(
			"EditName",
			"65cc5fd45aa71b64fbb551a9",
			"",
		).Return(expected, ErrName)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.EditStockName(
			"65cc5fd45aa71b64fbb551a9", 
			"",
		)

		assert.ErrorIs(t, err, ErrName)
	})
}

func TestEditStockSign(t *testing.T) {
	expected := "Successfully edit stock sign"

	t.Run("Edit stock sign", func(t *testing.T) {
		stockRepo.On(
			"EditSign",
			"65cc5fd45aa71b64fbb551a9",
			"T",
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.EditStockSign(
			"65cc5fd45aa71b64fbb551a9", 
			"T",
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid sign", func(t *testing.T) {
		stockRepo.On(
			"EditSign",
			"65cc5fd45aa71b64fbb551a9",
			"",
		).Return(expected, ErrName)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.EditStockSign(
			"65cc5fd45aa71b64fbb551a9", 
			"",
		)

		assert.ErrorIs(t, err, ErrName)
	})
}

func TestDeleteStockCollection(t *testing.T) {
	expected := "Successfully deleted stock"

	t.Run("Delete stock collection", func(t *testing.T) {
		stockRepo.On(
			"DeleteStock",
			"65cc5fd45aa71b64fbb551a9",
		).Return(expected, nil)
		stockService := service.NewStockService(stockRepo, redisClient)

		actual, err := stockService.DeleteStockCollection("65cc5fd45aa71b64fbb551a9")

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		stockRepo.On(
			"DeleteStock",
			"",
		).Return(expected, ErrInvalidStock)
		stockService := service.NewStockService(stockRepo, redisClient)

		_, err := stockService.DeleteStockCollection("")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})
}

