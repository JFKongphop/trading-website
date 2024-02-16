package service_test

import (
	"fmt"
	"server/errs"
	"server/model"
	"server/redis"
	"server/repository"
	"server/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAccount = model.CreateAccount
type OrderRequest = model.OrderRequest
type BalanceHistory = model.BalanceHistory
type UserAccount = model.UserAccount
type UserHistory = model.UserHistory
type UserStock = model.UserStock

var redisClient = redis.InitRedis()
var userRepo = repository.NewUserRepositoryDBMock()

var (
	ErrData         = errs.ErrData
	ErrMoney        = errs.ErrMoney
	ErrUser         = errs.ErrUser
	ErrInvalidStock = errs.ErrInvalidStock
	ErrOrderMethod  = errs.ErrOrderMethod
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

func TestDepositBalance(t *testing.T) {
	expected := "Successfully deposited money"

	t.Run("Error invalid money", func(t *testing.T) {
		userRepo.On(
			"Deposit",
			"65c8993c48096b5150cee5d6",
			float64(0),
		).Return(expected, ErrMoney)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.DepositBalance(
			"65c8993c48096b5150cee5d6",
			float64(0),
		)

		assert.Error(t, err, ErrMoney)
	})

	t.Run("Deposit balance", func(t *testing.T) {
		userRepo.On(
			"Deposit",
			"65c8993c48096b5150cee5d6",
			float64(1),
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.DepositBalance(
			"65c8993c48096b5150cee5d6",
			1,
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestWithdrawBalance(t *testing.T) {
	expected := "Successfully withdrawed money"

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"Withdraw",
			"",
			float64(1),
		).Return(expected, ErrUser)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.WithdrawBalance(
			"",
			1,
		)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Withdraw balance", func(t *testing.T) {
		userRepo.On(
			"Withdraw",
			"65c8993c48096b5150cee5d6",
			float64(1),
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.WithdrawBalance(
			"65c8993c48096b5150cee5d6",
			1,
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestBuyStock(t *testing.T) {
	expected := "Successfully bought stock"

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"Buy",
			OrderRequest{},
		).Return(expected, ErrUser)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.BuyStock(OrderRequest{})

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Buy stock", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     "65c39a03dfb8060d99995934",
			UserId:      "65c8993c48096b5150cee5d6",
			Price:       60,
			Amount:      8,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		userRepo.On(
			"Buy",
			orderRequest,
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.BuyStock(orderRequest)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestSaleStock(t *testing.T) {
	expected := "Successfully sold stock"

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"Sale",
			OrderRequest{},
		).Return(expected, ErrUser)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.SaleStock(OrderRequest{})

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Sale stock", func(t *testing.T) {
		orderRequest := model.OrderRequest{
			StockId:     "65bf707e040d36a26f4bf523",
			UserId:      "65c30de7b654c0e7bf938081",
			Price:       10,
			Amount:      100,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		userRepo.On(
			"Sale",
			orderRequest,
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.SaleStock(orderRequest)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestSetFavoriteStock(t *testing.T) {
	expected := "Successfully set favorite stock"
	t.Run("Error invalid stock", func(t *testing.T) {
		userRepo.On(
			"SetFavorite",
			"65c30de7b654c0e7bf938081",
			"",
		).Return(expected, ErrInvalidStock)
		userService := service.NewUserService(userRepo, redisClient)
		_, err := userService.SetFavoriteStock(
			"65c30de7b654c0e7bf938081",
			"",
		)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Set favorite stock", func(t *testing.T) {
		userRepo.On(
			"SetFavorite",
			"65c30de7b654c0e7bf938081",
			"65bf707e040d36a26f4bf523",
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.SetFavoriteStock(
			"65c30de7b654c0e7bf938081",
			"65bf707e040d36a26f4bf523",
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetUserBalanceHistory(t *testing.T) {
	expected := []BalanceHistory{{
		Timestamp: 1,
		Balance:   1,
		Method:    "DEPOSIT",
	}}

	t.Run("Error invalid method", func(t *testing.T) {
		userRepo.On(
			"GetBalanceHistory",
			"65c30de7b654c0e7bf938081",
			"",
			uint(1),
		).Return(expected, ErrOrderMethod)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.GetUserBalanceHistory(
			"65c30de7b654c0e7bf938081",
			"",
			1,
		)

		assert.ErrorIs(t, err, ErrOrderMethod)
	})

	t.Run("Get balance history", func(t *testing.T) {
		userRepo.
			On(
				"GetBalanceHistory",
				"65c30de7b654c0e7bf938081",
				"DEPOSIT",
				uint(1),
			).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.GetUserBalanceHistory(
			"65c30de7b654c0e7bf938081",
			"DEPOSIT",
			1,
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestGetUserBalance(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"GetBalance",
			"65c30de7b654c0e7bf938081",
		).Return(int(1), ErrUser)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.GetUserBalance("65c30de7b654c0e7bf938081")

		fmt.Println(err)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Get balance", func(t *testing.T) {
		userRepo.On(
			"GetBalance",
			"65c30de7b654c0e7bf938081",
		).Return(int(1), nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.GetUserBalance("65c30de7b654c0e7bf938081")

		fmt.Println(err)
		// assert.Empty(t, err)
		assert.Equal(t, float64(0), actual)
	})
}

func TestGetFavoriteStock(t *testing.T) {
	expected := []string{"65c30de7b654c0e7bf938081"}
	t.Run("Get favorite stock", func(t *testing.T) {
		userRepo.On(
			"GetFavorite",
			"65c30de7b654c0e7bf938081",
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.GetFavoriteStock("65c30de7b654c0e7bf938081")

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"GetFavorite",
			"",
		).Return(expected, ErrUser)

		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.GetFavoriteStock("")
		assert.ErrorIs(t, err, ErrUser)
	})
}

func TestGetUserAccount(t *testing.T) {
	expetced := UserAccount{
		ID:             primitive.ObjectID{},
		Name:           "kongphop",
		ProfileImage:   "test",
		Email:          "test",
		RegisterDate:   primitive.Timestamp{},
		Balance:        0,
		BalanceHistory: []BalanceHistory{},
		Favorite:       []string{},
		History:        []UserHistory{},
		Stock:          []UserStock{},
	}

	t.Run("Get user account", func(t *testing.T) {
		userRepo.On(
			"GetAccount",
			"65c30de7b654c0e7bf938081",
		).Return(expetced, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.GetUserAccount("65c30de7b654c0e7bf938081")

		assert.Empty(t, err)
		assert.Equal(t, expetced.Name, actual.Name)
	})

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"GetAccount",
			"",
		).Return(expetced, ErrUser)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.GetUserAccount("")

		assert.ErrorIs(t, err, ErrUser)
	})
}

func TestGetUserTradingHistories(t *testing.T) {
	expected := []UserHistory{{
		Timestamp:   0,
		StockId:     "1",
		Price:       1,
		Amount:      1,
		Status:      "pending",
		OrderType:   "auto",
		OrderMethod: "buy",
	}}

	t.Run("Get user trading histories", func(t *testing.T) {
		userRepo.On(
			"GetAllHistories",
			"65c30de7b654c0e7bf938081",
			uint(0),
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.GetUserTradingHistories(
			"65c30de7b654c0e7bf938081",
			uint(0),
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"GetAllHistories",
			"",
			uint(0),
		).Return(expected, ErrUser)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.GetUserTradingHistories(
			"",
			uint(0),
		)

		assert.ErrorIs(t, err, ErrUser)
	})
}

func TestGetUserStockHistory(t *testing.T) {
	expected := []UserHistory{{
		Timestamp:   0,
		StockId:     "1",
		Price:       1,
		Amount:      1,
		Status:      "pending",
		OrderType:   "auto",
		OrderMethod: "buy",
	}}

	t.Run("Get user stock history", func(t *testing.T) {
		userRepo.On(
			"GetUserStockHistory",
			"65c30de7b654c0e7bf938081",
			"65c30de7b654c0e7bf938081",
			uint(0),
		).Return(expected, nil)
		userRepo := service.NewUserService(userRepo, redisClient)

		actual, err := userRepo.GetUserStockHistory(
			"65c30de7b654c0e7bf938081",
			"65c30de7b654c0e7bf938081",
			uint(0),
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"GetUserStockHistory",
			"65c30de7b654c0e7bf938081",
			"",
			uint(0),
		).Return(expected, ErrInvalidStock)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.GetUserStockHistory(
			"65c30de7b654c0e7bf938081",
			"",
			uint(0),
		)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})
}

func TestGetUserStockAmount(t *testing.T) {
	expected := UserStock{
		StockId: "65c30de7b654c0e7bf938081",
		Amount: 1,
	}
	t.Run("Get user stock amount", func(t *testing.T) {
		userRepo.On(
			"GetStockAmount",
			"65c30de7b654c0e7bf938081",
			"65c30de7b654c0e7bf938081",
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.GetUserStockAmount(
			"65c30de7b654c0e7bf938081",
			"65c30de7b654c0e7bf938081",
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		userRepo.On(
			"GetStockAmount",
			"65c30de7b654c0e7bf938081",
			"",
		).Return(expected, ErrInvalidStock)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.GetUserStockAmount(
			"65c30de7b654c0e7bf938081",
			"",
		)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})
}

func TestDeleteFavoriteStock(t *testing.T) {
	expected := "Successfully deleted favorite stock"

	t.Run("Delete favorite stock", func(t *testing.T) {
		userRepo.On(
			"DeleteFavorite",
			"65c30de7b654c0e7bf938081",
			"65c30de7b654c0e7bf938081",
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.DeleteFavoriteStock(
			"65c30de7b654c0e7bf938081",
			"65c30de7b654c0e7bf938081",
		)

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		userRepo.On(
			"DeleteFavorite",
			"65c30de7b654c0e7bf938081",
			"",
		).Return(expected, ErrInvalidStock)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.DeleteFavoriteStock(
			"65c30de7b654c0e7bf938081",
			"",
		)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})
}

func TestDeleteUserAccount(t *testing.T) {
	expected := "Successfully deleted account"

	t.Run("Delete user account", func(t *testing.T) {
		userRepo.On(
			"DeleteAccount",
			"65c30de7b654c0e7bf938081",
		).Return(expected, nil)
		userService := service.NewUserService(userRepo, redisClient)

		actual, err := userService.DeleteUserAccount("65c30de7b654c0e7bf938081")

		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Error invalid user", func(t *testing.T) {
		userRepo.On(
			"DeleteAccount",
			"",
		).Return(expected, ErrUser)
		userService := service.NewUserService(userRepo, redisClient)

		_, err := userService.DeleteUserAccount("")

		assert.ErrorIs(t, err, ErrUser)
	})
}
