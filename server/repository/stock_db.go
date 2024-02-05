package repository

import "go.mongodb.org/mongo-driver/mongo"

type stockRepositoryDB struct {
	db *mongo.Collection
}

func NewStockRepositoryDB(db *mongo.Collection) StockRepository {
	return stockRepositoryDB{db}
}

func (r stockRepositoryDB) CreateStock(StockCollection) (string, error) {
	return "", nil
}

func (r stockRepositoryDB) GetAllStocks() ([]StockCollection, error) {
	return []StockCollection{}, nil
}

func (r stockRepositoryDB) GetTopStocks() ([]StockCollection, error) {
	return []StockCollection{}, nil
}

func (r stockRepositoryDB) GetStock(string) (StockCollection, error) {
	return StockCollection{}, nil
}

func (r stockRepositoryDB) EditStock(string) (string, error) {
	return "", nil
}

func (r stockRepositoryDB) 	DeleteStock(string) (string, error) {
	return "", nil
}