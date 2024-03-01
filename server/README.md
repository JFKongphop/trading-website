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
