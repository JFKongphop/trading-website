package repository

import "github.com/stretchr/testify/mock"

type stockRepositoryDBMock struct {
	mock.Mock
}

func NewStockRepositoryDBMock() *stockRepositoryDBMock {
	return &stockRepositoryDBMock{}
}

func (m *stockRepositoryDBMock) CreateStock(stockCollection StockCollection) (string, error) {
	arge := m.Called(stockCollection)
	return arge.String(0), arge.Error(1)
}

func (m *stockRepositoryDBMock) CreateStockOrder(stockId string, stockOrder StockHistory) (string, error) {
	arge := m.Called(stockId, stockOrder)
	return arge.String(0), arge.Error(1)
}

func (m *stockRepositoryDBMock) GetAllStocks() ([]StockCollectionResponse, error) {
	arge := m.Called()
	return arge.Get(0).([]StockCollectionResponse), arge.Error(1)
}

func (m *stockRepositoryDBMock) GetTopStocks() ([]StockGroup, error) {
	arge := m.Called()
	return arge.Get(0).([]StockGroup), arge.Error(1)
}

func (m *stockRepositoryDBMock) GetStock(stockId string) (StockCollectionResponse, error) {
	arge := m.Called(stockId)
	return arge.Get(0).(StockCollectionResponse), arge.Error(1)
}

func (m *stockRepositoryDBMock) GetFavoriteStock(favoriteStockIds []string) ([]StockCollectionResponse, error) {
	arge := m.Called(favoriteStockIds)
	return arge.Get(0).([]StockCollectionResponse), arge.Error(1)
}

func (m *stockRepositoryDBMock) GetStockHistory(stockId string) ([]StockHistoryResponse, error) {
	arge := m.Called(stockId)
	return arge.Get(0).([]StockHistoryResponse), arge.Error(1)
}

func (m *stockRepositoryDBMock) SetPrice(stockId string, price float64) (string, error) {
	arge := m.Called(stockId, price)
	return arge.String(0), arge.Error(1)
}

func (m *stockRepositoryDBMock) EditName(stockId string, name string) (string, error) {
	arge := m.Called(stockId, name)
	return arge.String(0), arge.Error(1)
}

func (m *stockRepositoryDBMock) EditSign(stockId string, sign string) (string, error) {
	arge := m.Called(stockId, sign)
	return arge.String(0), arge.Error(1)
}

func (m *stockRepositoryDBMock) DeleteStock(stockId string) (string, error) {
	arge := m.Called(stockId)
	return arge.String(0), arge.Error(1)
}
