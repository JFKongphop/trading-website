package repository_test

import (
	"fmt"
	"server/repository"
	"testing"
)

func InitStockRepo() repository.StockRepository {
	client, _ := repository.InitMongoDB("mongodb://localhost:27017/trading-system")
	db := client.Database("trading-system")
	collection := db.Collection("stock")
	userRepo := repository.NewStockRepositoryDB(collection)

	return userRepo
}

var stockRepo = InitStockRepo()

func TestCreateStock(t *testing.T) {

}

func TestCreateStockOrder(t *testing.T) {

}

func TestGetAllStocks(t *testing.T) {
	t.Run("test stock", func(t *testing.T) {
		actual, _ := stockRepo.GetAllStocks()

		fmt.Println(actual)
	})
}

func TestGetTopStocks(t *testing.T) {

} 

func TestGetStock(t *testing.T) {

}

func TestGetFavoriteStock(t *testing.T) {

}

func TestGetStockHistory(t *testing.T) {
	
}

func TestSetPrice(t *testing.T) {

}

func TestEditName(t *testing.T) {

}

func TestEditSign(t *testing.T) {

}

func TestDeleteStock(t *testing.T) {

}