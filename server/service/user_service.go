package service

import "server/repository"

type UserRepository = repository.UserRepository

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return userService{userRepo}
}

func (s userService) CreateUserAccount(CreateAccount) (string, error) {

	return "", nil
}

func (s userService) DepositBalance(string, float64) (string, error) {
	return "", nil
}

func (s userService) WithdrawBalance(string, float64) (string, error) {
	return "", nil
}

func (s userService) BuyStock(OrderRequest) (string, error) {
	return "", nil
}

func (s userService) SaleStock(OrderRequest) (string, error) {
	return "", nil
}

func (s userService) SetFavoriteStock(string, string) (string, error) {
	return "", nil
}

func (s userService) GetUserBalanceHistory(string, string, uint) ([]BalanceHistory, error) {
	return []BalanceHistory{}, nil
}

func (s userService) GetUserBalance(string) (float64, error) {
	return 0, nil
}

func (s userService) GetFavoriteStock(string) ([]string, error) {
	return []string{}, nil
}

func (s userService) GetUserAccount(string) (UserResponse, error) {
	return UserResponse{}, nil
}

func (s userService) GetUserTradingHistories(string) ([]ResponseUserHistory, error) {
	return []ResponseUserHistory{}, nil
}

func (s userService) GetUserStockHistory(string) ([]ResponseUserHistory, error) {
	return []ResponseUserHistory{}, nil
}

func (s userService) GetUserStockAmount(string, string) (UserStock, error) {
	return UserStock{}, nil
}

func (s userService) DeleteFavoriteStock(string, string) (string, error) {
	return "", nil
}

func (s userService) DeleteUserAccount(string, string) (string, error) {
	return "", nil
}