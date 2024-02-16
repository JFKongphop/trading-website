package repository

import "github.com/stretchr/testify/mock"

type userRepositoryDBMock struct {
	mock.Mock
}

func NewUserRepositoryDBMock() *userRepositoryDBMock {
	return &userRepositoryDBMock{}
}

func (m *userRepositoryDBMock) Create(data CreateAccount) (string, error) {
	arge := m.Called(data)
	return arge.String(0), arge.Error(1)
}

func (m *userRepositoryDBMock) Deposit(userId string, depositMoney float64) (string, error) {
	arge := m.Called(userId, depositMoney)
	return arge.String(0), arge.Error(1)
}

func (m *userRepositoryDBMock) Withdraw(userId string, withdrawMoney float64) (string, error) {
	arge := m.Called(userId, withdrawMoney)
	return arge.String(0), arge.Error(1)
}

func (m *userRepositoryDBMock) Buy(orderRequest OrderRequest) (string, error) {
	arge := m.Called(orderRequest)
	return arge.String(0), arge.Error(1)
}

func (m *userRepositoryDBMock) Sale(orderRequest OrderRequest) (string, error) {
	arge := m.Called(orderRequest)
	return arge.String(0), arge.Error(1)
}

func (m *userRepositoryDBMock) SetFavorite(userId string, stockId string) (string, error) {
	arge := m.Called(userId, stockId)
	return arge.String(0), arge.Error(1)
}

func (m *userRepositoryDBMock) GetBalanceHistory(userId string, method string, skip uint) ([]BalanceHistory, error) {
	arge := m.Called(userId, method, skip)
	return arge.Get(0).([]BalanceHistory), arge.Error(1)
}

func (m *userRepositoryDBMock) GetBalance(userId string) (float64, error) {
	arge := m.Called(userId)
	return float64(arge.Int(0)), arge.Error(1)
}

func (m *userRepositoryDBMock) GetFavorite(userId string) ([]string, error) {
	arge := m.Called(userId)
	return arge.Get(0).([]string), arge.Error(1)
}

func (m *userRepositoryDBMock) GetAccount(userId string) (UserAccount, error) {
	arge := m.Called(userId)
	return arge.Get(0).(UserAccount), arge.Error(1)
}

func (m *userRepositoryDBMock) GetAllHistories(userId string, start uint) ([]UserHistory, error) {
	arge := m.Called(userId, start)
	return arge.Get(0).([]UserHistory), arge.Error(1)
}

func (m *userRepositoryDBMock) GetUserStockHistory(userId string, stockId string, skip uint) ([]UserHistory, error) {
	arge := m.Called(userId, stockId, skip)
	return arge.Get(0).([]UserHistory), arge.Error(1)
}

func (m *userRepositoryDBMock) GetStockAmount(userId string, stockId string) (UserStock, error) {
	arge := m.Called(userId, stockId)
	return arge.Get(0).(UserStock), arge.Error(1)
}

func (m *userRepositoryDBMock) DeleteFavorite(userId string, stockId string) (string, error) {
	arge := m.Called(userId, stockId)
	return arge.String(0), arge.Error(1)
}

func (m *userRepositoryDBMock) DeleteAccount(userId string) (string, error) {
	arge := m.Called(userId)
	return arge.String(0), arge.Error(1)
}
