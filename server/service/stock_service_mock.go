package service

import "github.com/stretchr/testify/mock"

type stockServiceMock struct {
	mock.Mock
}

func NewStockServiceMock() *stockServiceMock {
	return &stockServiceMock{}
}

func (m *stockServiceMock) CreateStockCollection(stockCollection StockCollectionRequest) (string, error) {
	arge := m.Called(stockCollection)
	return arge.String(0), arge.Error(1)
}

func (m *stockServiceMock) CreateStockOrder(stockId string, stockOrder StockHistory) (string, error) {
	arge := m.Called(stockId, stockOrder)
	return arge.String(0), arge.Error(1)
}

func (m *stockServiceMock) GetAllStockCollections() ([]StockCollectionResponse, error) {
	arge := m.Called()
	return arge.Get(0).([]StockCollectionResponse), arge.Error(1)
}

func (m *stockServiceMock) GetTop10Stocks() ([]TopStock, error) {
	arge := m.Called()
	return arge.Get(0).([]TopStock), arge.Error(1)
}

func (m *stockServiceMock) GetStockCollection(stockId string) (StockCollectionResponse, error) {
	arge := m.Called(stockId)
	return arge.Get(0).(StockCollectionResponse), arge.Error(1)
}

func (m *stockServiceMock) GetFavoriteStock(favoriteStockIds []string) ([]StockCollectionResponse, error) {
	arge := m.Called(favoriteStockIds)
	return arge.Get(0).([]StockCollectionResponse), arge.Error(1)
}

func (m *stockServiceMock) GetStockHistory(stockHistory string) ([]StockHistoryResponse, error) {
	arge := m.Called(stockHistory)
	return arge.Get(0).([]StockHistoryResponse), arge.Error(1)
}

func (m *stockServiceMock) GetStockPrice(stockId string) (float64, error) {
	arge := m.Called(stockId)
	return float64(arge.Int(0)), arge.Error(1)
}

func (m *stockServiceMock) GetStockGraph(stockId string) ([]Graph, error) {
	arge := m.Called(stockId)
	return arge.Get(0).([]Graph), arge.Error(1)
}

func (m *stockServiceMock) SetStockPrice(stockId string, price float64) (string, error) {
	arge := m.Called(stockId, price)
	return arge.String(0), arge.Error(1)
}

func (m *stockServiceMock) EditStockName(stockId string, name string) (string, error) {
	arge := m.Called(stockId, name)
	return arge.String(0), arge.Error(1)
}

func (m *stockServiceMock) EditStockSign(stockId string, sign string) (string, error) {
	arge := m.Called(stockId, sign)
	return arge.String(0), arge.Error(1)
}

func (m *stockServiceMock) DeleteStockCollection(stockId string) (string, error) {
	arge := m.Called(stockId)
	return arge.String(0), arge.Error(1)
}