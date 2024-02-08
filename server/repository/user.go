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
	SetFavorite(string) (string, error)
	GetBalanceHistory(string) ([]BalanceHistory, error)
	GetBalance(string) (float64, error)
	GetFavorite(string) ([]StockCollectionResponse, error)
	GetAccount(string) (UserAccount, error)
	GetAllHistories(string) ([]UserHistory, error)
	GetStockHistory(string, string) ([]UserHistory, error)
	GetStockAmount(string, string) (UserStock, error) 
	DeleteFavorite(string) (string, error)
	DeleteAccount(string) (string, error)
}


// {
//   "_id": {
//     "$oid": "65c39b189f5c807c54a53030"
//   },
//   "name": "kongphop",
//   "profileImage": "test-image",
//   "email": "test@gmail.com",
//   "registerDate": {
//     "$timestamp": {
//       "t": 1707318040,
//       "i": 1
//     }
//   },
//   "balance": 4400,
//   "userHistory": [
//     {
//       "timestamp": {
//         "$numberLong": "1707389401"
//       },
//       "stockId": "65c39a12c4e3672bcbf15b0f",
//       "price": 60,
//       "amount": 5,
//       "status": "pending",
//       "orderType": "auto",
//       "orderMethod": "buy"
//     },
//     {
//       "timestamp": {
//         "$numberLong": "1707389516"
//       },
//       "stockId": "65c39a12c4e3672bcbf15b0f",
//       "price": 60,
//       "amount": 5,
//       "status": "pending",
//       "orderType": "auto",
//       "orderMethod": "buy"
//     }
//   ],
//   "userStock": [
//     {
//       "stockId": "65c39a12c4e3672bcbf15b0f",
//       "amount": 10
//     }
//   ]
// }