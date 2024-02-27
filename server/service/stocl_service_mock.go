package service

import "github.com/stretchr/testify/mock"

type stockServiceMock struct {
	mock.Mock
}

func NewStockServiceMock() *stockServiceMock {
	return &stockServiceMock{}
}

func (m *stockServiceMock) CreateStockCollection(StockCollectionRequest) (string, error)
func (m *stockServiceMock) CreateStockOrder(string, StockHistory) (string, error)
func (m *stockServiceMock) GetAllStockCollections() ([]StockCollectionResponse, error)
func (m *stockServiceMock) GetTop10Stocks() ([]TopStock, error)
func (m *stockServiceMock) GetStockCollection(string) (StockCollectionResponse, error)
func (m *stockServiceMock) GetFavoriteStock([]string) ([]StockCollectionResponse, error)
func (m *stockServiceMock) GetStockHistory(string) ([]StockHistoryResponse, error) 
func (m *stockServiceMock) GetStockPrice(string) (float64, error)
func (m *stockServiceMock) GetStockGraph(string) ([]Graph, error)
func (m *stockServiceMock) SetStockPrice(string, float64) (string, error)
func (m *stockServiceMock) EditStockName(string, string) (string, error)
func (m *stockServiceMock) EditStockSign(string, string) (string, error)
func (m *stockServiceMock) DeleteStockCollection(string) (string, error)