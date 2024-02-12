package service

import (
	"server/repository"

	"github.com/redis/go-redis/v9"
)

type UserRepository = repository.UserRepository

type userService struct {
	userRepo UserRepository
	redisClient *redis.Client
}

func NewUserService(userRepo UserRepository, redisClient *redis.Client) UserService {
	return userService{userRepo, redisClient}
}

func (s userService) CreateUserAccount(userAccount CreateAccount) (string, error) {
	result, err := s.userRepo.Create(userAccount)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s userService) DepositBalance(userId string, depositMoney float64) (string, error) {
	result, err := s.userRepo.Deposit(userId, depositMoney)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s userService) WithdrawBalance(userId string, withdrawMoney float64) (string, error) {
	result, err := s.userRepo.Withdraw(userId, withdrawMoney)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s userService) BuyStock(orderRequest OrderRequest) (string, error) {
	result, err := s.userRepo.Buy(orderRequest)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s userService) SaleStock(orderRequest OrderRequest) (string, error) {
	result, err := s.userRepo.Sale(orderRequest)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s userService) SetFavoriteStock(userId string, stockId string) (string, error) {
	result, err := s.userRepo.SetFavorite(userId, stockId)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s userService) GetUserBalanceHistory(userId string, method string, perpage uint) ([]BalanceHistory, error) {
	
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