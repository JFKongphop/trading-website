# Documentation
## User
* [Signup](#signup)
* [Deposit](#deposit)
* [Withdraw](#withdraw)
* [Buy](#buy)
* [Sale](#sale)
* [Set Favorite](#set-favorite)
* [Balance Transaction](#balance-transaction)
* [Balance](#balance)
* [Get Favorite](#get-favorite)
* [Signin](#signin)
* [Trade Transaction](#trade-transaction)
* [Stock Transaction](#stock-transaction)
* [Stock Amount](#stock-amount)
* [Delete Favorite](#delete-favorite)
* [Delete Account](#delete-account)

## Stock
* [Create Stock](#create-stock)
* [Create Order](#create-order)
* [Stock Collections](#collections)
* [Top Stocks](#top-stocks)
* [Stock Collection](#collection)
* [Transaction](#transaction)
* [Price](#price)
* [Graph](#graph)
* [Set Price](#set-price)
* [Edit Name](#edit-name)
* [Edit Sign](#edit-sign)
* [Delete Stock Collection](#delete-stock)

#

### Signup
singup for create new user.
```http
POST /api/v1/user/signup
```
#### Request
```javascript
{
  "uid": string
  "name": string,
  "profileImage": string,
  "email": string
}
```
#### Response
```javascript
{
  "message": "Successfully created account"
}
```
#

### Deposit
deposit money for trading stock.
```http
POST /api/v1/user/deposit
```
#### Request
```javascript
{
  "balance": int
}
```
#### Response
```javascript
{
  "message": "Successfully deposited money"
}
```
#

### Withdraw
withdraw money.
```http
POST /api/v1/user/withdraw
```
#### Request
```javascript
{
  "balance": int
}
```
#### Response
```javascript
{
  "message": "Successfully withdrawed money"
}
```
#

### Buy
Buy stock.
```http
POST /api/v1/user/buy
```
##### Available Order Type
- order
- auto
##### Available Order Method
- buy
- sale
#### Request
```javascript
{
	"stockId": string,
	"userId": string,
	"price": int,
	"amount": int
	"orderType": string,
	"orderMethod": string
}
```
#### Response
```javascript
{
  "message": "Successfully bought stock"
}
```
#

### Sale
Sale stock.
```http
POST /api/v1/user/sale
```
##### Available Order Type
- order
- auto
##### Available Order Method
- buy
- sale
#### Request
```javascript
{
	"stockId": string,
	"userId": string,
	"price": int,
	"amount": int
	"orderType": string,
	"orderMethod": string
}
```
#### Response
```javascript
{
  "message": "Successfully sold stock"
}
```
#

### Set Favorite
Set favorite stock from user.
```http
POST /api/v1/user/set-favorite
```
#### Request
```javascript
{
	"stockId": string,
}
```
#### Response
```javascript
{
  "message": "Successfully set favorite stock"
}
```
#

### Balance Transaction
Get balance transaction user.
##### Available Methods
- ALL
- DEPOSIT
- WITHDRAW
```http
GET /api/v1/user/balance-transaction?startPage=0&method=ALL
```
#### Response
```javascript
{
  "message": "Successfully set favorite stock"
}
```
#

### Balance
Get balance user.
```http
GET /api/v1/user/balance
```
#### Response
```javascript
{
  "balance": int,
  "message": "Successfully fetched user balance"
}
```
#

### Get Favorite
Get favorite stock.
```http
GET /api/v1/user/get-favorite
```
#### Response
```javascript
{
  "favorites": [
    {
      "id": string,
      "stockImage": string,
      "name": string,
      "sign": string,
      "price": int
    }
  ],
  "message": "Successfully fetched favorite stock"
}
```
#

### Signin
get user account by signin.
```http
GET /api/v1/user/signin
```
#### Response
```javascript
{
  "message": "Successfully fetched user profile",
  "user": {
    "name": string,
    "profileImage": string,
    "email": string
  }
}
```
#

### Trade Transaction
get all stock transactions.
```http
GET /api/v1/user/trade-transaction?startPage=0
```
#### Response
```javascript
{
  "message": "Successfully fetched all transactions history",
  "transactions": [
    {
      "timestamp": int,
      "stockId": string,
      "price": int,
      "amount": int,
      "status": string,
      "orderType": string,
      "orderMethod": string
    },
  ]
}
```
#

### Stock Transaction
get stock transaction by stockId.
```http
GET /api/v1/user/trade-transaction?stockId=65cc5fd45aa71b64fbb551a9&startPage=0
```
#### Response
```javascript
{
  "message": "Successfully fetched all transactions history",
  "transactions": [
    {
      "timestamp": int,
      "stockId": string,
      "price": int,
      "amount": int,
      "status": string,
      "orderType": string,
      "orderMethod": string
    },
  ]
}
```
#

### Stock Amount
get amount of user stock.
```http
GET /api/v1/user/stock-ratio?stockId=65bf707e040d36a26f4bf523
```
#### Response
```javascript
{
  "message": "Successfully fetched stock ratio",
  "stockRatio": {
    "stockId": string,
    "amount": int
  }
}
```
#

### Delete Favorite
delete stock favorite.
```http
DELETE /api/v1/user/delete-favorite?stockId=65bf707e040d36a26f4bf523
```
#### Response
```javascript
{
  "message": "Successfully deleted favorite stock"
}
```
#

### Delete Account
delete account.
```http
DELETE /api/v1/user/delete-account
```
#### Response
```javascript
{
  "message": "Successfully deleted account"
}
```
#