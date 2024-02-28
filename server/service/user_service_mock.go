package service

import "github.com/stretchr/testify/mock"

type userServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *userServiceMock {
	return &userServiceMock{}
}

func (m *userServiceMock) CreateUserAccount(userAccount CreateAccount) (string, error) {
	arge := m.Called(userAccount)
	return arge.String(0), arge.Error(1)
}

func (m *userServiceMock) DepositBalance(userId string, depositMoney float64) (string, error) {
	arge := m.Called(userId, depositMoney)
	return arge.String(0), arge.Error(1)
}

func (m *userServiceMock) WithdrawBalance(userId string, withdrawMoney float64) (string, error) {
	arge := m.Called(userId, withdrawMoney)
	return arge.String(0), arge.Error(1)
}

func (m *userServiceMock) BuyStock(orderRequest OrderRequest) (string, error) {
	arge := m.Called(orderRequest)
	return arge.String(0), arge.Error(1)
}

func (m *userServiceMock) SaleStock(orderRequest OrderRequest) (string, error) {
	arge := m.Called(orderRequest)
	return arge.String(0), arge.Error(1)
}

func (m *userServiceMock) SetFavoriteStock(userId string, stockId string) (string, error) {
	arge := m.Called(userId, stockId)
	return arge.String(0), arge.Error(1)
}

func (m *userServiceMock) GetUserBalanceHistory(userId string, method string, perPage uint) ([]BalanceHistory, error) {
	arge := m.Called(userId, method, perPage)
	return arge.Get(0).([]BalanceHistory), arge.Error(1)
}

func (m *userServiceMock) GetUserBalance(userId string) (float64, error) {
	arge := m.Called(userId)
	return float64(arge.Int(0)), arge.Error(1)
}

func (m *userServiceMock) GetUserFavoriteStock(userId string) ([]string, error) {
	arge := m.Called(userId)
	return arge.Get(0).([]string), arge.Error(1)
}

func (m *userServiceMock) GetUserAccount(userId string) (UserResponse, error) {
	arge := m.Called(userId)
	return arge.Get(0).(UserResponse), arge.Error(1)
}

func (m *userServiceMock) GetUserTradingHistories(userId string, startPage uint) ([]ResponseUserHistory, error) {
	arge := m.Called(userId, startPage)
	return arge.Get(0).([]ResponseUserHistory), arge.Error(1)
}

func (m *userServiceMock) GetUserStockHistory(userId string, stockId string, startPage uint) ([]ResponseUserHistory, error) {
	arge := m.Called(userId, stockId, startPage)
	return arge.Get(0).([]ResponseUserHistory), arge.Error(1)
}

func (m *userServiceMock) GetUserStockAmount(userId string, stockId string) (UserStock, error) {
	arge := m.Called(userId, stockId)
	return arge.Get(0).(UserStock), arge.Error(1)
}

func (m *userServiceMock) DeleteFavoriteStock(userId string, stockId string) (string, error) {
	arge := m.Called(userId, stockId)
	return arge.String(0), arge.Error(1)
}

func (m *userServiceMock) DeleteUserAccount(userId string) (string, error) {
	arge := m.Called(userId)
	return arge.String(0), arge.Error(1)
}
