package errs

import "errors"

var (
	ErrSignin = errors.New("account is already existed")
	ErrDB = errors.New("cannot connect database")
	ErrData = errors.New("invalid data")
	ErrUser = errors.New("invalid user")
	ErrMoney = errors.New("invalid money")
	ErrInvalidStock = errors.New("invalid stock")
	ErrBalance = errors.New("balance not enough")
	ErrNotEnoughStock = errors.New("stock not enough")
	ErrOrderType = errors.New("invalid order type")
	ErrOrderMethod = errors.New("invalid order method")
	ErrFavoriteStock = errors.New("already set favorite stock")
)