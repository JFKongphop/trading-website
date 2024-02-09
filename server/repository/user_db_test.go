package repository_test

import (
	"server/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	//"go.mongodb.org/mongo-driver/mongo"
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

var userRepo = InitUseRealrRepo()
var userIdTesting = "65c5da45d40c5e57a213b200"
var (
	ErrData = repository.ErrData
	ErrMoney = repository.ErrMoney
) 

func TestInitMongoDB(t *testing.T) {
	t.Run("Error cannot connect database", func(t *testing.T) {
		_, err := repository.InitMongoDB("")
		assert.ErrorIs(t, err, repository.ErrDB)
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

}