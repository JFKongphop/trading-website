package repository

import "server/model"


type StockCollection = model.StockCollection
type AllStock = model.AllStock
type TopStock = model.TopStock

type StockRepository interface {
	CreateStock(StockCollection) (string, error)
	GetAllStocks() ([]StockCollection, error)
	GetTopStocks() ([]StockCollection, error)
	GetStock(string) (StockCollection, error)
	EditStock(string) (string, error)
	DeleteStock(string) (string, error)
}
