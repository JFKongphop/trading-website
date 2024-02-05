package service

import "server/model"

type CreateAccount = model.CreateAccount
type UserResponse model.UserResponse
type ResponseUserHistory = model.UserHistory

type UserService interface {
	CreateUserAccount(CreateAccount) (string, error)
	GetUserAccount(string) (UserResponse, error)
	GetUserTradingHistories(string) ([]ResponseUserHistory, error)
	GetUserStockHistory(string) ([]ResponseUserHistory, error)
	DeleteUserAccount(string, string) (string, error)
}