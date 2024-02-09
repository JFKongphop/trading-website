package repository

import "errors"

var (
	ErrDB = errors.New("cannot connect database")
	ErrData = errors.New("invalid data")
	ErrUser = errors.New("invalid user")
	ErrMoney = errors.New("invalid money")
	ErrStock = errors.New("invalid stock")
)