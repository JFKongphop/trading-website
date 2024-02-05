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

func (s userService) GetUserAccount(string) (UserResponse, error) {
	return UserResponse{}, nil
}

func (s userService) GetUserTradingHistories(string) ([]ResponseUserHistory, error) {
	return []ResponseUserHistory{}, nil
}

func (s userService) GetUserStockHistory(string) ([]ResponseUserHistory, error) {
	return []ResponseUserHistory{}, nil
}

func (s userService) DeleteUserAccount(string, string) (string, error) {
	return "", nil
}