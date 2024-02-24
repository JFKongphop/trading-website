package repository

import "server/model"


type StockCollection = model.StockCollection
type AllStock = model.AllStock
type TopStock = model.TopStock
type StockHistory = model.StockHistory
type StockCollectionResponse = model.StockCollectionResponse
type StockHistoryResponse = model.StockHistoryResponse

// first his 1707319626

type StockRepository interface {
	CreateStock(StockCollection) (string, error)
	CreateStockOrder(string, StockHistory) (string, error)
	GetAllStocks() ([]StockCollectionResponse, error)
	GetTopStocks() ([]StockGroup, error)
	GetStock(string) (StockCollectionResponse, error)
	GetFavoriteStock([]string) ([]StockCollectionResponse, error)
	GetStockHistory(string) ([]StockHistoryResponse, error) 
	GetPrice(string) (float64, error)
	SetPrice(string, float64) (string, error)
	EditName(string, string) (string, error)
	EditSign(string, string) (string, error)
	DeleteStock(string) (string, error)
}
