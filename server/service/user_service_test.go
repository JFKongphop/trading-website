package service_test

import (
	"server/errs"
	"server/model"
	"server/redis"
	"server/repository"
	"server/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CreateAccount = model.CreateAccount

var redisClient = redis.InitRedis()
var userRepo = repository.NewUserRepositoryDBMock()

var (
	ErrData = errs.ErrData
	ErrMoney = errs.ErrMoney
)

func TestCreateUserAccount(t *testing.T) {
	expected := "Successfully created account"

	t.Run("Error invalid data", func(t *testing.T) {
		userRepo.On("Create", CreateAccount{}).Return(expected, ErrData)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.CreateUserAccount(CreateAccount{})

		assert.ErrorIs(t, err, ErrData)
	})

	t.Run("Create user account", func(t *testing.T) {
		account := CreateAccount{
			Name:         "kongphop",
			ProfileImage: "test",
			Email:        "test@gmail.com",
		}

		userRepo.On("Create", account).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.CreateUserAccount(account)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestDepositBalance(t *testing.T){
	expected := "Successfully deposited money"

	t.Run("Error invalid money", func(t *testing.T) {
		userRepo.
		On(
			"Deposit",
			"65c8993c48096b5150cee5d6",
			float64(0),
		).
		Return(expected, ErrMoney)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.DepositBalance(
			"65c8993c48096b5150cee5d6",
			float64(0),
		)

		assert.Error(t, err, ErrMoney)
	})

	t.Run("Deposit balance", func(t *testing.T) {
		userRepo.
			On(
				"Deposit",
				"65c8993c48096b5150cee5d6",
				float64(1),
			).
			Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual ,err := userService.DepositBalance(
			"65c8993c48096b5150cee5d6",
			1,
		)
 
		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestWithdrawBalance(t *testing.T){}

func TestBuyStock(t *testing.T){}

func TestSaleStock(t *testing.T){}

func TestSetFavoriteStock(t *testing.T){}

func TestGetUserBalanceHistory(t *testing.T){}

func TestGetUserBalance(t *testing.T){}

func TestGetFavoriteStock(t *testing.T){}

func TestGetUserAccount(t *testing.T){}

func TestGetUserTradingHistories(t *testing.T){}

func TestGetUserStockHistory(t *testing.T){}

func TestGetUserStockAmount(t *testing.T){}

func TestDeleteFavoriteStock(t *testing.T){}

func TestDeleteUserAccount(t *testing.T){}


