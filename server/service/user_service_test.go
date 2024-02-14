package service_test

import (
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

func TestCreateUserAccount(t *testing.T) {
	t.Run("Create user account", func(t *testing.T) {
		account := CreateAccount{
			Name:         "kongphop",
			ProfileImage: "test",
			Email:        "test@gmail.com",
		}
		expected := "Successfully created account"

		userRepo.On("Create", account).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.CreateUserAccount(account)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestDepositBalance(t *testing.T){
	t.Run("Deposit balance", func(t *testing.T) {})
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


