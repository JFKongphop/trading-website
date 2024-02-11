package repository_test

// func InitStockRepo() repository.UserRepository {
// 	client, _ := repository.InitMongoDB("mongodb://localhost:27017/trading-system")
// 	db := client.Database("trading-system")
// 	collection := db.Collection("stock")
// 	userRepo := repository.NewUserRepositoryDB(collection)

// 	return userRepo
// }

var stockRepo = InitStockRepo()

// func Test

// CreateStock(StockCollection) (string, error)
// CreateStockOrder(string, StockHistory) (string, error)
// GetAllStocks() ([]StockCollectionResponse, error)
// GetTopStocks() ([]StockGroup, error)
// GetStock(string) (StockCollectionResponse, error)
// GetFavoriteStock([]string) ([]StockCollectionResponse, error)
// GetStockHistory(string) ([]StockHistoryResponse, error) 
// SetPrice(string, float64) (string, error)
// EditName(string, string) (string, error)
// EditSign(string, string) (string, error)
// DeleteStock(string) (string, error)