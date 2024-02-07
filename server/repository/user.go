package repository

import "server/model"

type CreateAccount = model.CreateAccount
type UserAccount = model.UserAccount
type UserHistory = model.UserHistory
type UserStock = model.UserStock
type OrderRequest = model.OrderRequest

type UserRepository interface {
	Create(CreateAccount) (string, error)
	Deposit(string, float64) (string, error)
	Withdraw(string, float64) (string, error)
	Buy(OrderRequest) (string, error)
	Sale(OrderRequest) (string, error)
	GetAccount(string) (UserAccount, error)
	GetAllHistories(string) ([]UserHistory, error)
	GetStockHistory(string, string) ([]UserHistory, error)
	GetStockAmount(string, string) (UserStock, error) 
	DeleteAccount(string) (string, error)
}
