package service

import "server/model"

type StockCollection = model.StockCollection
type AllStock = model.AllStock
type StockHistory = model.StockHistory
type TopStock = model.TopStock
type StockCollectionResponse = model.StockCollectionResponse
type StockHistoryResponse = model.StockHistoryResponse
type StockCollectionRequest = model.StockCollectionRequest


type StockService interface {
	CreateStockCollection(StockCollectionRequest) (string, error)
	CreateStockOrder(string, StockHistory) (string, error)
	GetAllStockCollections() ([]StockCollectionResponse, error)
	GetTop10Stocks() ([]TopStock, error)
	GetStockCollection(string) (StockCollectionResponse, error)
	GetFavoriteStock([]string) ([]StockCollectionResponse, error)
	GetStockHistory(string) ([]StockHistoryResponse, error) 
	SetStockPrice(string, float64) (string, error)
	EditStockName(string, string) (string, error)
	EditStockSign(string, string) (string, error)
	DeleteStockCollection(string) (string, error)
}