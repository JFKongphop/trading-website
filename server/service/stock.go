package service

import "server/model"

type StockCollection = model.StockCollection
type AllStock = model.AllStock
type TopStock = model.TopStock

type StockService interface {
	CreateStockCollection(StockCollection) (string, error)
	GetAllStockCollections() ([]AllStock, error)
	GetTop10Stocks() ([]TopStock, error)
	GetStockCollection(string) (StockCollection, error)
	EditStockCollection(string) (string, error)
	DeleteStockCollection(string) (string, error)
}