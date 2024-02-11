package repository_test

import (
	"context"
	"server/errs"
	"server/model"
	"server/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitUserRepo() repository.UserRepository {
	client, _ := repository.InitMongoDB("mongodb://localhost:27017/trading-system")
	db := client.Database("trading-system")
	collection := db.Collection("user")
	userRepo := repository.NewUserRepositoryDB(collection)

	return userRepo
}

type CreateAccount = repository.CreateAccount
type OrderRequest = model.OrderRequest
type UserHistory = model.UserHistory

var userRepo = InitUserRepo()
var userIdTesting = "65c896695ec42b4f4f77af63"
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
	ErrFavoriteStock = errs.ErrFavoriteStock
) 

type mockCollection struct {
	mock.Mock
}

var ctx = context.Background()

func (m *mockCollection) InsertOne(document bson.M) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

type MockMongoCollection struct {
	*mongo.Collection
	*mockCollection
}

func TestInitMongoDB(t *testing.T) {
	t.Run("Error cannot connect database", func(t *testing.T) {
		_, err := repository.InitMongoDB("")
		assert.ErrorIs(t, err, errs.ErrDB)
	})
}

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

func TestDeposit(t *testing.T) {
	t.Run("Error invalid money", func(t *testing.T) {
		depositMoney := -1

		_, err := userRepo.Deposit(userIdTesting, float64(depositMoney))

		assert.ErrorIs(t, err, ErrMoney)
	})

	t.Run("Error invalid user", func(t *testing.T) {
		depositMoney := 1

		_, err := userRepo.Deposit("", float64(depositMoney))

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		depositMoney := 1

		_, err := userRepo.Deposit("teste", float64(depositMoney))

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error no documents in result", func(t *testing.T) {
		depositMoney := 1

		_, err := userRepo.Deposit("65c896695ec42b4f4f77af61", float64(depositMoney))

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Successfully deposited money", func(t *testing.T) {
		depositMoney := 5000

		actual, _ := userRepo.Deposit(userIdTesting, float64(depositMoney))
		expected := "Successfully deposited money"

		assert.Equal(t, expected, actual)		
	})
}

func TestWithdraw(t *testing.T) {
	t.Run("Error invalid money", func(t *testing.T) {
		withdrawMoney := -1

		_, err := userRepo.Withdraw(userIdTesting, float64(withdrawMoney))

		assert.ErrorIs(t, err, ErrMoney)
	})

	t.Run("Error invalid user", func(t *testing.T) {
		withdrawMoney := 1

		_, err := userRepo.Withdraw("", float64(withdrawMoney))

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		withdrawMoney := 1

		_, err := userRepo.Withdraw("teste", float64(withdrawMoney))

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error no documents in result", func(t *testing.T) {
		depositMoney := 1

		_, err := userRepo.Withdraw("65c896695ec42b4f4f77af61", float64(depositMoney))

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Error balance not enough", func(t *testing.T) {
		withdrawMoney := 1_000_000

		_, err := userRepo.Withdraw(userIdTesting, float64(withdrawMoney))

		assert.ErrorIs(t, err, ErrBalance)
	})

	t.Run("Successfully withdrawed money", func(t *testing.T) {
		withdrawMoney := 1000

		actual, _ := userRepo.Withdraw(userIdTesting, float64(withdrawMoney))
		expected := "Successfully withdrawed money"

		assert.Equal(t, expected, actual)		
	})
}

func TestBuy(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      "",
			Price:       60,
			Amount:      5,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		_, err := userRepo.Buy(orderRequest)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error invalid data", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     "",
			UserId:      stockIdTesting,
			Price:       -1,
			Amount:      -1,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		_, err := userRepo.Buy(orderRequest)

		assert.ErrorIs(t, err, ErrData)
	})
	
	t.Run("Error invalid order type", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       1,
			Amount:      1,
			OrderType:   "test",
			OrderMethod: "buy",
		}

		_, err := userRepo.Buy(orderRequest)

		assert.ErrorIs(t, err, ErrOrderType)
	})

	t.Run("Error invalid order method", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       1,
			Amount:      1,
			OrderType:   "auto",
			OrderMethod: "test",
		}

		_, err := userRepo.Buy(orderRequest)

		assert.ErrorIs(t, err, ErrOrderMethod)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      "test",
			Price:       60,
			Amount:      5,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		_, err := userRepo.Buy(orderRequest)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error no documents in result", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      "65c896695ec42b4f4f77af61",
			Price:       60,
			Amount:      5,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		_, err := userRepo.Buy(orderRequest)

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Error balance not enough", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       1_000_000,
			Amount:      1_000_000,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		_, err := userRepo.Buy(orderRequest)

		assert.ErrorIs(t, err, ErrBalance)
	})

	t.Run("Successfully bought stock", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       50,
			Amount:      10,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		actual, _ := userRepo.Buy(orderRequest)
		expected := "Successfully bought stock"

		assert.Equal(t, expected, actual)	
	})
}

func TestSale(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      "",
			Price:       60,
			Amount:      5,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		_, err := userRepo.Sale(orderRequest)

		assert.ErrorIs(t, err, ErrUser)
	})
	
	t.Run("Error invalid data", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     "",
			UserId:      stockIdTesting,
			Price:       -1,
			Amount:      -1,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		_, err := userRepo.Sale(orderRequest)

		assert.ErrorIs(t, err, ErrData)
	})
	
	t.Run("Error invalid order type", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       1,
			Amount:      1,
			OrderType:   "test",
			OrderMethod: "sale",
		}

		_, err := userRepo.Sale(orderRequest)

		assert.ErrorIs(t, err, ErrOrderType)
	})

	t.Run("Error invalid order method", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       1,
			Amount:      1,
			OrderType:   "auto",
			OrderMethod: "test",
		}

		_, err := userRepo.Sale(orderRequest)

		assert.ErrorIs(t, err, ErrOrderMethod)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      "test",
			Price:       60,
			Amount:      5,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		_, err := userRepo.Sale(orderRequest)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error no documents in result", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      "65c896695ec42b4f4f77af61",
			Price:       60,
			Amount:      5,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		_, err := userRepo.Sale(orderRequest)

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Error balance not enough", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       1_000_000,
			Amount:      1_000_000,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		_, err := userRepo.Sale(orderRequest)

		assert.ErrorIs(t, err, ErrNotEnoughStock)
	})

	t.Run("Successfully sold stock", func(t *testing.T) {
		orderRequest := OrderRequest{
			StockId:     stockIdTesting,
			UserId:      userIdTesting,
			Price:       1,
			Amount:      1,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		actual, _ := userRepo.Sale(orderRequest)
		expected := "Successfully sold stock"

		assert.Equal(t, expected, actual)
	})
}

func TestSetFavorite(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.SetFavorite("", stockIdTesting)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := userRepo.SetFavorite(userIdTesting, "")

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.SetFavorite("test", stockIdTesting)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error alreay set favorite stock", func(t *testing.T) {
		_, err := userRepo.SetFavorite("65c896695ec42b4f4f77af61", stockIdTesting)

		assert.ErrorIs(t, err, ErrFavoriteStock)	
	})

	t.Run("Successfully set favorite stock", func(t *testing.T) {
		actual, _ := userRepo.SetFavorite(userIdTesting, stockIdTesting)
		expected := "Successfully set favorite stock"

		assert.Equal(t, expected, actual)
	})
}

func TestGetBalanceHistory(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.GetBalanceHistory("", "", 0)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error invalid method", func(t *testing.T) {
		_, err := userRepo.GetBalanceHistory(userIdTesting, "", 0)

		assert.ErrorIs(t, err, ErrOrderMethod)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.GetBalanceHistory("test", "ALL", 0)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})
}

func TestBalance(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.GetBalance("")

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.GetBalance("test")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error no documents in result", func(t *testing.T) {
		_, err := userRepo.GetBalance("65c8993c48096b5150cee5d1")

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Get balance user", func(t *testing.T) {
		actual, _ := userRepo.GetBalance("65c8993c48096b5150cee5d6")
		expected := 3220

		assert.Equal(t, float64(expected), actual)
	})
}

func TestGetFavorite(t *testing.T) {
	t.Run("Error in valid user", func(t *testing.T) {
		_, err := userRepo.GetFavorite("")

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.GetFavorite("test")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error no documents in result", func(t *testing.T) {
		_, err := userRepo.GetFavorite("65c8993c48096b5150cee5d1")

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})
	
	t.Run("Get favorite stock", func(t *testing.T) {
		actual, _ := userRepo.GetFavorite("65c8993c48096b5150cee5d6")
		expected := []string{"65c39a03dfb8060d99995934"}

		assert.Equal(t, expected, actual)
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("Error invalud user", func(t *testing.T) {
		_, err := userRepo.GetAccount("")

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.GetAccount("test")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error not document in result", func(t *testing.T) {
		_, err := userRepo.GetAccount("65c8993c48096b5150cee5d1")

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Get account", func(t *testing.T) {
		actual, _ := userRepo.GetAccount("65c8993c48096b5150cee5d6")
		expected := "kongphop"

		assert.Equal(t, expected, actual.Name)
	})
}

func TestGetAllHistories(t *testing.T)	{
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.GetAllHistories("", 0)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.GetAllHistories("test", 0)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error not document in result", func(t *testing.T) {
		_, err := userRepo.GetAllHistories("65c8993c48096b5150cee5d1", 0)

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Get account", func(t *testing.T) {
		actual, _ := userRepo.GetAllHistories("65c8993c48096b5150cee5d6", 0)
		expectedTimestamp := int64(1707645450)

		assert.Equal(t, expectedTimestamp, actual[len(actual) - 1].Timestamp)
	})
}

func TestGetStockHistory(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.GetStockHistory(
			"", 
			stockIdTesting, 
			0,
		)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := userRepo.GetStockHistory(
			"65c8993c48096b5150cee5d6", 
			"", 
			0,
		)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.GetStockHistory(
			"test", 
			stockIdTesting, 
			0,
		)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Get stock history", func(t *testing.T) {
		actual, _ := userRepo.GetStockHistory(
			"65c8993c48096b5150cee5d6", 
			stockIdTesting, 
			0,
		)
		expectedTimestamp := int64(1707645423)

		assert.Equal(t, expectedTimestamp, actual[len(actual) - 1].Timestamp)
	})
}

func TestGetStockAmount(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.GetStockAmount(
			"", 
			stockIdTesting,
		)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := userRepo.GetStockAmount(
			"65c8993c48096b5150cee5d6", 
			"",
		)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.GetStockAmount(
			"test", 
			stockIdTesting, 
		)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Error not document in result", func(t *testing.T) {
		_, err := userRepo.GetStockAmount("65c8993c48096b5150cee5d1", stockIdTesting)

		assert.Equal(t, err.Error(), "mongo: no documents in result")
	})

	t.Run("Get stock amount", func(t *testing.T) {
		actual, _ := userRepo.GetStockAmount("65c8993c48096b5150cee5d6", stockIdTesting)
		expectedStockId := stockIdTesting
		expectedStockAmount := float64(13)

		assert.Equal(t, expectedStockId, actual.StockId)
		assert.Equal(t, expectedStockAmount, expectedStockAmount)
	})
}

func TestDeleteFavorite(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.DeleteFavorite(
			"", 
			stockIdTesting,
		)

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error invalid stock", func(t *testing.T) {
		_, err := userRepo.DeleteFavorite(
			"65c8993c48096b5150cee5d6", 
			"",
		)

		assert.ErrorIs(t, err, ErrInvalidStock)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.DeleteFavorite(
			"test", 
			stockIdTesting, 
		)

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Delete Favorite", func(t *testing.T) {
		actual, _ := userRepo.DeleteFavorite(
			"65c8993c48096b5150cee5d6",
			"65c39a03dfb8060d99995934",
		)
		expected := "Successfully deleted favorite stock"

		assert.Equal(t, expected, actual)
	})
}

func TestDeleteAccount(t *testing.T) {
	t.Run("Error invalid user", func(t *testing.T) {
		_, err := userRepo.DeleteAccount("")

		assert.ErrorIs(t, err, ErrUser)
	})

	t.Run("Error convert userId to objectId", func(t *testing.T) {
		_, err := userRepo.DeleteAccount("test")

		assert.Equal(t, err.Error(), "the provided hex string is not a valid ObjectID")
	})

	t.Run("Delete account", func(t *testing.T) {
		actual, _ := userRepo.DeleteAccount("65c89a4c3a0dd78a473a5638")
		expected := "Successfully deleted account"

		assert.Equal(t, expected, actual)
	})
}