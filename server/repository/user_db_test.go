package repository_test

import (
	"server/errs"
	"server/model"
	"server/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/mongo/integration/mtest"
	//"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// Create(CreateAccount) (string, error)
// Deposit(string, float64) (string, error)
// Withdraw(string, float64) (string, error)
// Buy(OrderRequest) (string, error)
// Sale(OrderRequest) (string, error)
// SetFavorite(string, string) (string, error)
// GetBalanceHistory(string, string) ([]BalanceHistory, error)
// GetBalance(string) (float64, error)
// GetFavorite(string) ([]string, error)
// GetAccount(string) (UserAccount, error)
// GetAllHistories(string) ([]UserHistory, error)
// GetStockHistory(string, string) ([]UserHistory, error)
// GetStockAmount(string, string) (UserStock, error)
// DeleteFavorite(string, string) (string, error)
// DeleteAccount(string) (string, error)

func InitUseRealrRepo() repository.UserRepository {
	client, _ := repository.InitMongoDB("mongodb://localhost:27017/trading-system")
	db := client.Database("trading-system")
	collection := db.Collection("user")
	userRepo := repository.NewUserRepositoryDB(collection)

	return userRepo
}
 

type CreateAccount = repository.CreateAccount
type OrderRequest = model.OrderRequest

var userRepo = InitUseRealrRepo()
var userIdTesting = "65c5e8fcdc5093167512fb31"
var stockIdTesting = "65c39a12c4e3672bcbf15b0f"
var (
	ErrData = errs.ErrData
	ErrUser = errs.ErrUser
	ErrMoney = errs.ErrMoney
	ErrBalance = errs.ErrBalance
	ErrOrderType = errs.ErrOrderType
	ErrOrderMethod = errs.ErrOrderMethod	
	ErrNotEnoughStock = errs.ErrNotEnoughStock
	ErrInvalidStock = errs.ErrInvalidStock
) 

// func TestInitMongoDB(t *testing.T) {
// 	t.Run("Error cannot connect database", func(t *testing.T) {
// 		_, err := repository.InitMongoDB("")
// 		assert.ErrorIs(t, err, errs.ErrDB)
// 	})
// }

// func TestCreate(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	
// 	client := mt.Client
// 	db := client.Database("trading-system")

// 	repo := repository.NewUserRepositoryDB(db.Collection("users"))

//   data := CreateAccount{
//     Name:           "John Doe",
//     ProfileImage:   "https://example.com/profile.jpg",
//     Email:          "johndoe@example.com",
//   }

// 	result, err := repo.Create(data)

//   assert.NoError(t, err)
//   assert.Equal(t, "Successfully created account", result)
// }

func TestCreate(t *testing.T) {
	t.Run("Error invalid data", func(t *testing.T) {
		data := CreateAccount{
			Name: "",
			ProfileImage: "",
			Email: "",
		}

		_, err := userRepo.Create(data)
		assert.ErrorIs(t, err, ErrData)
	})

	t.Run("Create user account", func(t *testing.T) {
		account := CreateAccount{
			Name:         "test",
			ProfileImage: "test",
			Email:        "test@gmail.com",
		}

		expected := "Successfully created account"
		actual, _ := userRepo.Create(account)
		
		assert.Equal(t, expected, actual)
	})
}

// func TestDeposit(t *testing.T) {
// 	t.Run("Error invalid money", func(t *testing.T) {
// 		depositMoney := -1

// 		_, err := userRepo.Deposit(userIdTesting, float64(depositMoney))

// 		assert.ErrorIs(t, err, ErrMoney)
// 	})

// 	t.Run("Error invalid user", func(t *testing.T) {
// 		depositMoney := 1

// 		_, err := userRepo.Deposit("", float64(depositMoney))

// 		assert.ErrorIs(t, err, ErrUser)
// 	})

// 	t.Run("Error convert userId to objectId", func(t *testing.T) {
// 		depositMoney := 1

// 		_, err := userRepo.Deposit("teste", float64(depositMoney))

// 		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
// 	})

// 	t.Run("Successfully deposited money", func(t *testing.T) {
// 		depositMoney := 5000

// 		actual, _ := userRepo.Deposit(userIdTesting, float64(depositMoney))
// 		expected := "Successfully deposited money"

// 		assert.Equal(t, expected, actual)		
// 	})
// }

// func TestWithdraw(t *testing.T) {
// 	t.Run("Error invalid money", func(t *testing.T) {
// 		withdrawMoney := -1

// 		_, err := userRepo.Withdraw(userIdTesting, float64(withdrawMoney))

// 		assert.ErrorIs(t, err, ErrMoney)
// 	})

// 	t.Run("Error invalid user", func(t *testing.T) {
// 		withdrawMoney := 1

// 		_, err := userRepo.Withdraw("", float64(withdrawMoney))

// 		assert.ErrorIs(t, err, ErrUser)
// 	})

// 	t.Run("Error convert userId to objectId", func(t *testing.T) {
// 		withdrawMoney := 1

// 		_, err := userRepo.Withdraw("teste", float64(withdrawMoney))

// 		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
// 	})

// 	t.Run("Error balance not enough", func(t *testing.T) {
// 		withdrawMoney := 1_000_000

// 		_, err := userRepo.Withdraw(userIdTesting, float64(withdrawMoney))

// 		assert.ErrorIs(t, err, ErrBalance)
// 	})

// 	t.Run("Successfully withdrawed money", func(t *testing.T) {
// 		withdrawMoney := 5000

// 		actual, _ := userRepo.Withdraw(userIdTesting, float64(withdrawMoney))
// 		expected := "Successfully withdrawed money"

// 		assert.Equal(t, expected, actual)		
// 	})
// }

// func TestBuy(t *testing.T) {
// 	t.Run("Error invalid user", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      "",
// 			Price:       60,
// 			Amount:      5,
// 			OrderType:   "auto",
// 			OrderMethod: "buy",
// 		}

// 		_, err := userRepo.Buy(orderRequest)

// 		assert.ErrorIs(t, err, ErrUser)
// 	})

// 	t.Run("Error invalid data", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     "",
// 			UserId:      stockIdTesting,
// 			Price:       -1,
// 			Amount:      -1,
// 			OrderType:   "auto",
// 			OrderMethod: "buy",
// 		}

// 		_, err := userRepo.Buy(orderRequest)

// 		assert.ErrorIs(t, err, ErrData)
// 	})
	
// 	t.Run("Error invalid order type", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       1,
// 			Amount:      1,
// 			OrderType:   "test",
// 			OrderMethod: "buy",
// 		}

// 		_, err := userRepo.Buy(orderRequest)

// 		assert.ErrorIs(t, err, ErrOrderType)
// 	})

// 	t.Run("Error invalid order method", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       1,
// 			Amount:      1,
// 			OrderType:   "auto",
// 			OrderMethod: "test",
// 		}

// 		_, err := userRepo.Buy(orderRequest)

// 		assert.ErrorIs(t, err, ErrOrderMethod)
// 	})

// 	t.Run("Error convert userId to objectId", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      "test",
// 			Price:       60,
// 			Amount:      5,
// 			OrderType:   "auto",
// 			OrderMethod: "buy",
// 		}

// 		_, err := userRepo.Buy(orderRequest)

// 		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
// 	})

// 	t.Run("Error balance not enough", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       1_000_000,
// 			Amount:      1_000_000,
// 			OrderType:   "auto",
// 			OrderMethod: "buy",
// 		}

// 		_, err := userRepo.Buy(orderRequest)

// 		assert.ErrorIs(t, err, ErrBalance)
// 	})

// 	t.Run("Successfully bought stock", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       50,
// 			Amount:      10,
// 			OrderType:   "auto",
// 			OrderMethod: "buy",
// 		}

// 		actual, _ := userRepo.Buy(orderRequest)
// 		expected := "Successfully bought stock"

// 		assert.Equal(t, expected, actual)	
// 	})
// }

// func TestSale(t *testing.T) {
// 	t.Run("Error invalid user", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      "",
// 			Price:       60,
// 			Amount:      5,
// 			OrderType:   "auto",
// 			OrderMethod: "sale",
// 		}

// 		_, err := userRepo.Sale(orderRequest)

// 		assert.ErrorIs(t, err, ErrUser)
// 	})
	
// 	t.Run("Error invalid data", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     "",
// 			UserId:      stockIdTesting,
// 			Price:       -1,
// 			Amount:      -1,
// 			OrderType:   "auto",
// 			OrderMethod: "sale",
// 		}

// 		_, err := userRepo.Sale(orderRequest)

// 		assert.ErrorIs(t, err, ErrData)
// 	})
	
// 	t.Run("Error invalid order type", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       1,
// 			Amount:      1,
// 			OrderType:   "test",
// 			OrderMethod: "sale",
// 		}

// 		_, err := userRepo.Sale(orderRequest)

// 		assert.ErrorIs(t, err, ErrOrderType)
// 	})

// 	t.Run("Error invalid order method", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       1,
// 			Amount:      1,
// 			OrderType:   "auto",
// 			OrderMethod: "test",
// 		}

// 		_, err := userRepo.Sale(orderRequest)

// 		assert.ErrorIs(t, err, ErrOrderMethod)
// 	})

// 	t.Run("Error convert userId to objectId", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      "test",
// 			Price:       60,
// 			Amount:      5,
// 			OrderType:   "auto",
// 			OrderMethod: "sale",
// 		}

// 		_, err := userRepo.Sale(orderRequest)

// 		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
// 	})

// 	t.Run("Error balance not enough", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       1_000_000,
// 			Amount:      1_000_000,
// 			OrderType:   "auto",
// 			OrderMethod: "sale",
// 		}

// 		_, err := userRepo.Sale(orderRequest)

// 		assert.ErrorIs(t, err, ErrNotEnoughStock)
// 	})

// 	t.Run("Successfully sold stock", func(t *testing.T) {
// 		orderRequest := OrderRequest{
// 			StockId:     stockIdTesting,
// 			UserId:      userIdTesting,
// 			Price:       1,
// 			Amount:      1,
// 			OrderType:   "auto",
// 			OrderMethod: "sale",
// 		}

// 		actual, _ := userRepo.Sale(orderRequest)
// 		expected := "Successfully sold stock"

// 		assert.Equal(t, expected, actual)
// 	})
// }

// func TestSetFavorite(t *testing.T) {
// 	t.Run("Error invalid user", func(t *testing.T) {
// 		_, err := userRepo.SetFavorite("", stockIdTesting)

// 		assert.ErrorIs(t, err, ErrUser)
// 	})

// 	t.Run("Error invalid stock", func(t *testing.T) {
// 		_, err := userRepo.SetFavorite(userIdTesting, "")

// 		assert.ErrorIs(t, err, ErrInvalidStock)
// 	})

// 	t.Run("Error convert userId to objectId", func(t *testing.T) {
// 		_, err := userRepo.SetFavorite("test", stockIdTesting)

// 		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
// 	})

// 	t.Run("Successfully set favorite stock", func(t *testing.T) {
// 		actual, _ := userRepo.SetFavorite(userIdTesting, stockIdTesting)
// 		expected := "Successfully set favorite stock"

// 		assert.Equal(t, expected, actual)
// 	})
// }

// func TestGetBalanceHistory(t *testing.T) {
// 	t.Run("Error invalid user", func(t *testing.T) {
// 		_, err := userRepo.GetBalanceHistory("", "")

// 		assert.ErrorIs(t, err, ErrUser)
// 	})

// 	t.Run("Error invalid method", func(t *testing.T) {
// 		_, err := userRepo.GetBalanceHistory(userIdTesting, "")

// 		assert.ErrorIs(t, err, ErrOrderMethod)
// 	})

// 	t.Run("Error convert userId to objectId", func(t *testing.T) {
// 		_, err := userRepo.GetBalanceHistory("test", "ALL")

// 		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
// 	})
// }