package repository_test

import (
	"server/errs"
	"server/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StockCollection = repository.StockCollection
type StockHistory = repository.StockHistory

func InitStockRepo() repository.StockRepository {
	client, _ := repository.InitMongoDB("mongodb://localhost:27017/trading-system")
	db := client.Database("trading-system")
	collection := db.Collection("stock")
	userRepo := repository.NewStockRepositoryDB(collection)

	return userRepo
}

var stockRepo = InitStockRepo()
var (
	ErrPrice = errs.ErrPrice
	ErrName = errs.ErrName
	ErrSign = errs.ErrSign
) 

func TestCreateStock(t *testing.T) {
	t.Run("Error invalid data", func(t *testing.T) {
		stockCollection := StockCollection{
			StockImage: "",
			Name:       "",
			Sign:       "",
			Price:      0,
			History:    []StockHistory{},
		}

		_, err := stockRepo.CreateStock(stockCollection)

		assert.ErrorIs(t, err, ErrData)
	})

	t.Run("Create stock", func(t *testing.T) {
		stockCollection := StockCollection{
			StockImage: "test",
			Name:       "test",
			Sign:       "test",
			Price:      1,
			History:    []StockHistory{},
		}

		actual, _ := stockRepo.CreateStock(stockCollection)
		expected := "Successfullt created stock collection"

		assert.Equal(t, expected, actual)
	})
}

func TestCreateStockOrder(t *testing.T) {
	t.Run("Error invalid data", func(t *testing.T) {
		// 65c99e67b244d2f0231ed667
		stockOrder := StockHistory{}
		_, err := stockRepo.CreateStockOrder("", stockOrder)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error invalid data", func(t *testing.T) {
		stockOrder := StockHistory{}
		_, err := stockRepo.CreateStockOrder("65c99e67b244d2f0231ed667", stockOrder)

		assert.ErrorIs(t, err, ErrData)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		stockOrder := StockHistory{
			ID: "65c39b189f5c807c54a53030",
			Amount: 5,
			Price: 13,
		}

		_, err := stockRepo.CreateStockOrder("test", stockOrder)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error no documents in result", func(t *testing.T) {
		stockOrder := StockHistory{
			ID: "65c39b189f5c807c54a53030",
			Amount: 5,
			Price: 13,
		}
		_, err := stockRepo.CreateStockOrder("65c99e67b244d2f0231ed660", stockOrder)

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Create stock order", func(t *testing.T) {
		stockOrder := StockHistory{
			ID: "65c39b189f5c807c54a53030",
			Amount: 5,
			Price: 13,
		}

		actual, _ := stockRepo.CreateStockOrder("65c99e67b244d2f0231ed667", stockOrder)
		expected := "Successfully created stock order"
		assert.Equal(t, expected, actual)
	})
}

func TestGetAllStocks(t *testing.T) {
	t.Run("Get all stocka", func(t *testing.T) {
		actual, _ := stockRepo.GetAllStocks()
		expectedFirstStockName := "KONGPHOP"
		expectedSecondStockName := "TANAGUY"

		assert.Equal(t, expectedFirstStockName, actual[0].Name)
		assert.Equal(t, expectedSecondStockName, actual[1].Name)
	})
}

func TestGetTopStocks(t *testing.T) {
	t.Run("Get top stcok", func(t *testing.T) {
		actual, _ := stockRepo.GetTopStocks()
		expectedTopStockName := "TANAGUY"

		assert.Equal(t, expectedTopStockName, actual[0].Name)
	})
} 

func TestGetStock(t *testing.T) {
	t.Run("Error in valid stock", func(t *testing.T) {
		_, err := stockRepo.GetStock("")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := stockRepo.GetStock("test")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Get stock", func(t *testing.T) {
		actual, _ := stockRepo.GetStock("65c39a03dfb8060d99995934")
		expectedStockName := "KONGPHOP"

		assert.Equal(t, expectedStockName, actual.Name)
	})
}

func TestGetFavoriteStock(t *testing.T) {
	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := stockRepo.GetFavoriteStock([]string{""})

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := stockRepo.GetFavoriteStock([]string{"test"})

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Get favorite stock", func(t *testing.T) {
		actual, _ := stockRepo.GetFavoriteStock([]string{"65c39a03dfb8060d99995934"})
		expectedStockName := "KONGPHOP"

		assert.Equal(t, expectedStockName, actual[0].Name)
	})
}

func TestGetStockHistory(t *testing.T) {
	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := stockRepo.GetStockHistory("")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := stockRepo.GetStockHistory("test")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Get stock history", func(t *testing.T) {
		actual, _ := stockRepo.GetStockHistory("65c39a03dfb8060d99995934")
		expectedStockAmount := float64(25)
		expectedStockPrice := float64(11.11)

		assert.Equal(t, expectedStockAmount, actual[0].Amount)
		assert.Equal(t, expectedStockPrice, actual[0].Price)
	})
}

func TestSetPrice(t *testing.T) {
	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := stockRepo.SetPrice("", 0)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error invalid price", func(t *testing.T) {
		_, err := stockRepo.SetPrice("test", 0)

		assert.ErrorIs(t, err, ErrPrice)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := stockRepo.SetPrice("test", 1)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Set price", func(t *testing.T) {
		actual, _ := stockRepo.SetPrice("65c99e67b244d2f0231ed667", 2)
		expected := "Successfully set price"

		assert.Equal(t, expected, actual)
	})
}

func TestEditName(t *testing.T) {
// 65c99e67b244d2f0231ed667
	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := stockRepo.EditName("", "")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error invalid name", func(t *testing.T) {
		_, err := stockRepo.EditName("test", "")

		assert.ErrorIs(t, err, ErrName)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := stockRepo.EditName("test", "T")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Edit name", func(t *testing.T) {
		actual, _ := stockRepo.EditName("65c99e67b244d2f0231ed667", "T")
		expected := "Successfully updated name"

		assert.Equal(t, expected, actual)
	})
}

func TestEditSign(t *testing.T) {
	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := stockRepo.EditSign("", "")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error invalid sign", func(t *testing.T) {
		_, err := stockRepo.EditSign("test", "")

		assert.ErrorIs(t, err, ErrSign)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := stockRepo.EditSign("test", "T")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Edit sign", func(t *testing.T) {
		actual, _ := stockRepo.EditSign("65c99e67b244d2f0231ed667", "T")
		expected := "Successfully updated sign"

		assert.Equal(t, expected, actual)
	})
}

func TestDeleteStock(t *testing.T) {
	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := stockRepo.DeleteStock("")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := stockRepo.DeleteStock("test")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Delete stock", func(t *testing.T) {
		actual, _ := stockRepo.DeleteStock("65c99e6c02a43e12a634f777")
		expected := "Successfully deleted stock"

		assert.Equal(t, expected, actual)
	})
}