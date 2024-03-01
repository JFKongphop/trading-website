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
* [Get Price](#get-price)
* [Get Graph](#get-graph)
* [Set Price](#set-price)
* [Edit Name](#edit-name)
* [Edit Sign](#edit-sign)
* [Delete Stock Collection](#delete-stock)

#

## User

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

## Stock

### Create Stock
create stock collection.
```http
POST /api/v1/stock/create-stock
```
#### Request
| Key          | Value         |
|--------------|---------------|
| stockImage   | filename.jpg  |
| price        | 1.1           |
| name         | test          |
| sign         | t             |

#### Response
```javascript
{
  "message": "Successfully created stock collection"
}
```
#

### Create Order
create stock order.
```http
POST /api/v1/stock/create-order/:stockId
```
#### Request
```javascript
{
  "amount": 2,
  "price": 10
}
```
#### Response
```javascript
{
  "message": "Successfully created stock collection"
}
```
#

### Stock Collections
get collections.
```http
GET /api/v1/stock/collections
```
#### Response
```javascript
{
  "message": "Successfully fetched all stocks",
  "stocks": [
    {
      "id": string,
      "stockImage": string,
      "name": string,
      "sign": string,
      "price": int
    },
  ]
}
```
#

### Top Stock
get top 10 stock collections.
```http
GET /api/v1/stock/top-stocks
```
#### Response
```javascript
{
  "message": "Successfully fetched top volume stock",
  "topStocks": [
    {
      "id": string,
      "sign": string,
      "price": int
    }
  ]
}
```
#

### Stock Collection
get collection by stockId.
```http
GET /api/v1/stock/collection/:stockId
```
#### Response
```javascript
{
  "message": "Successfully fetched stock",
  "stock": {
    "id": string,
    "stockImage": string,
    "name": string,
    "sign": string,
    "price": int
  }
}
```
#

### Transaction
get transaction stock
```http
GET /api/v1/stock/transaction/:stockId
```
#### Response
```javascript
{
  "message": "Successfully fetched transactions",
  "transactions": [
    {
      "amount": int,
      "price": int
    },
    {
      "amount": int,
      "price": int
    }
  ]
}
```
#

### Get Price
get price stock
```http
GET /api/v1/stock/price/:stockId
```
#### Response
```javascript
{
  "message": "Successfully fetched transactions",
  "price": int
}
```
#

### Get Graph
get price stock
```http
GET /api/v1/stock/garph/:stockId
```
#### Response
```javascript
{
  "graph": [
    {
      "x": int,
      "y": [int, int,int, int]
    },
    {
      "x": int,
      "y": [int, int,int, int]
    },
  ]
  "message": "Successfully fetched stock graph"
}
```
#

### Set Price
set price stock
```http
POST /api/v1/stock/set-price/:stockId
```
#### Request
```javascript
{
  "price": 11.1
}
```
#### Response
```javascript
{
  "message": "Successfully set price"
}
```
#

### Edit Name
edit name stock
```http
POST /api/v1/stock/edit-name/:stockId
```
#### Request
```javascript
{
  "name": "test"
}
```
#### Response
```javascript
{
  "message": "Successfully updated name"
}
```
#

### Edit Sign
edit sign stock
```http
POST /api/v1/stock/edit-ign/:stockId
```
#### Request
```javascript
{
  "sign": "t"
}
```
#### Response
```javascript
{
  "message": "Successfully updated sign"
}
```
#

### Delete Stock Collection
delete stock colllection
```http
DELETE /api/v1/stock/delete/:stockId
```
#### Response
```javascript
{
  "message": "Successfully deleted stock"
}
```

