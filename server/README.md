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
#### Request Body
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
#### Request body
```javascript
{
  "balance": int
}
```
### Response
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
#### Request body
```javascript
{
  "balance": int
}
```
### Response
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
#### Request body
```javascript
{
	"stockId": string,
	"userId": string,
	"price": int,
	"amount": int
	"orderType": string, // order | auto
	"orderMethod": string // buy | sale
}
```
### Response
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
#### Request body
```javascript
{
	"stockId": string,
	"userId": string,
	"price": int,
	"amount": int
	"orderType": string, // order | auto
	"orderMethod": string // buy | sale
}
```
### Response
```javascript
{
  "message": "Successfully sold stock"
}
```
#
