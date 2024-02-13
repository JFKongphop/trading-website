package service

import "server/model"

type CreateAccount = model.CreateAccount
type UserResponse model.UserResponse
type ResponseUserHistory = model.UserHistory
type OrderRequest = model.OrderRequest
type BalanceHistory = model.BalanceHistory
type UserStock = model.UserStock

type UserService interface {
	CreateUserAccount(CreateAccount) (string, error)
	DepositBalance(string, float64) (string, error)
	WithdrawBalance(string, float64) (string, error)
	BuyStock(OrderRequest) (string, error)
	SaleStock(OrderRequest) (string, error)
	SetFavoriteStock(string, string) (string, error)
	GetUserBalanceHistory(string, string, uint) ([]BalanceHistory, error)
	GetUserBalance(string) (float64, error)
	GetFavoriteStock(string) ([]string, error)
	GetUserAccount(string) (UserResponse, error)
	GetUserTradingHistories(string, uint) ([]ResponseUserHistory, error)
	GetUserStockHistory(string, string, uint) ([]ResponseUserHistory, error)
	GetUserStockAmount(string, string) (UserStock, error) 
	DeleteFavoriteStock(string, string) (string, error)
	DeleteUserAccount(string) (string, error)
}