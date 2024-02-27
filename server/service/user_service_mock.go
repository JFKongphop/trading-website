package service

// import "github.com/stretchr/testify/mock"

// type userServiceMock struct {
// 	mock.Mock
// }

// func NewUserServiceMock() *userServiceMock {
// 	return &userServiceMock{}
// }


// func (m *userServiceMock) CreateUserAccount(userAccount CreateAccount) (string, error) {
// 	arge := m.Called()
// }

// func (m *userServiceMock) DepositBalance(string, float64) (string, error) {
// 	arge := m.Called()
// }

// func (m *userServiceMock) WithdrawBalance(string, float64) (string, error) {
// 	arge := m.Called()
// }

// func (m *userServiceMock) BuyStock(OrderRequest) (string, error) {
// 	arge := m.Called()
// }

// func (m *userServiceMock) SaleStock(OrderRequest) (string, error) {
// 	arge := m.Called()
// }

// func (m *userServiceMock) SetFavoriteStock(string, string) (string, error) {
// 	arge := m.Called()
// }

// func (m *userServiceMock) GetUserBalanceHistory(string, string, uint) ([]BalanceHistory, error) {
// 	arge := m.Called()
// }
// func (m *userServiceMock) GetUserBalance(string) (float64, error)
// func (m *userServiceMock) GetUserFavoriteStock(string) ([]string, error)
// func (m *userServiceMock) GetUserAccount(string) (UserResponse, error)
// func (m *userServiceMock) GetUserTradingHistories(string, uint) ([]ResponseUserHistory, error)
// func (m *userServiceMock) GetUserStockHistory(string, string, uint) ([]ResponseUserHistory, error)
// func (m *userServiceMock) GetUserStockAmount(string, string) (UserStock, error) 
// func (m *userServiceMock) DeleteFavoriteStock(string, string) (string, error)
// func (m *userServiceMock) DeleteUserAccount(string) (string, error)
