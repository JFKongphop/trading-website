package service

import (
	"server/repository"

	"github.com/redis/go-redis/v9"
)

type StockRepository = repository.StockRepository

type stockService struct {
	stockRepo   StockRepository
	redisClient *redis.Client
}

// CreateStockCollection(StockCollection) (string, error)
// CreateStockOrder(string, StockHistory) (string, error)
// GetAllStockCollections() ([]AllStock, error)
// GetTop10Stocks() ([]TopStock, error)
// GetStockCollection(string) (StockCollection, error)
// GetFavoriteStock([]string) ([]StockCollectionResponse, error)
// GetStockHistory(string) ([]StockHistoryResponse, error)
// SetStockPrice(string, float64) (string, error)
// EditStockName(string, string) (string, error)
// EditStockSign(string, string) (string, error)
// DeleteStockCollection(string) (string, error)

func NewStockService(stockRepo StockRepository, redisClient *redis.Client) StockService {
	return stockService{stockRepo, redisClient}
}

func (s stockService) CreateStockCollection(StockCollection) (string, error) {
	return "", nil
}

func (s stockService) CreateStockOrder(string, StockHistory) (string, error) {
	return "", nil
}

func (s stockService) GetAllStockCollections() ([]AllStock, error) {
	return []AllStock{}, nil
}

func (s stockService) GetTop10Stocks() ([]TopStock, error) {
	return []TopStock{}, nil
}

func (s stockService) GetStockCollection(string) (StockCollection, error) {
	return StockCollection{}, nil
}

func (s stockService) GetFavoriteStock([]string) ([]StockCollectionResponse, error) {
	return []StockCollectionResponse{}, nil
}

func (s stockService) GetStockHistory(string) ([]StockHistoryResponse, error) {
	return []StockHistoryResponse{}, nil
}

func (s stockService) SetStockPrice(string, float64) (string, error) {
	return "", nil
}

func (s stockService) EditStockName(string, string) (string, error) {
	return "", nil
}

func (s stockService) EditStockSign(string, string) (string, error) {
	return "", nil
}

func (s stockService) DeleteStockCollection(string) (string, error) {
	return "", nil
}
