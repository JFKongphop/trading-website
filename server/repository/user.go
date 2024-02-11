package repository

import "server/model"

type CreateAccount = model.CreateAccount
type UserAccount = model.UserAccount
type UserHistory = model.UserHistory
type UserStock = model.UserStock
type OrderRequest = model.OrderRequest
type BalanceHistory = model.BalanceHistory

type UserRepository interface {
	Create(CreateAccount) (string, error)
	Deposit(string, float64) (string, error)
	Withdraw(string, float64) (string, error)
	Buy(OrderRequest) (string, error)
	Sale(OrderRequest) (string, error)
	SetFavorite(string, string) (string, error)
	GetBalanceHistory(string, string, uint) ([]BalanceHistory, error)
	GetBalance(string) (float64, error)
	// GetStockBalance(string)
	// GetStockValueRatio(string)
	// GetStockAmountRatio(string) ()
	GetFavorite(string) ([]string, error)
	GetAccount(string) (UserAccount, error)
	GetAllHistories(string, uint) ([]UserHistory, error)
	GetUserStockHistory(string, string, uint) ([]UserHistory, error)
	GetStockAmount(string, string) (UserStock, error) 
	DeleteFavorite(string, string) (string, error)
	DeleteAccount(string) (string, error)
}