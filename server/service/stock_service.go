package service

import "server/repository"

type StockRepository = repository.StockRepository

type stockService struct {
	stockRepo StockRepository
}

func NewStockService(stockRepo StockRepository) StockService {
	return stockService{stockRepo}
}

func (s stockService) CreateStockCollection(StockCollection) (string, error) {
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

func (s stockService) EditStockCollection(string) (string, error) {
	return "", nil
}

func (s stockService) DeleteStockCollection(string) (string, error) {
	return "", nil
}