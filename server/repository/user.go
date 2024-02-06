package repository

import "server/model"

type CreateAccount = model.CreateAccount
type UserAccount = model.UserAccount
type UserHistory = model.UserHistory
type OrderRequest = model.OrderRequest

type UserRepository interface {
	Create(CreateAccount) (string, error)
	Buy(OrderRequest) (string, error)
	Sale(OrderRequest) (string, error)
	GetAccount(string) (UserAccount, error)
	GetAllHistories(string) ([]UserHistory, error)
	GetStockHistory(string) ([]UserHistory, error)
	DeleteAccount(string, string) (string, error)
}
