package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"server/model"
	"server/repository"
	"server/service"
	"server/redis"
	"time"

	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var ctx = context.Background()

type Person struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
	City string             `bson:"city"`
}

// type UserHistory struct {
// 	Timestamp uint   `bson:"timestamp"`
// 	Name      string `bson:"name"`
// 	Price     uint   `bson:"price"`
// 	Amount    uint   `bson:"amount"`
// 	Status    string `bson:"status"`
// }

// type UserStock struct {
// 	Name   string `bson:"name"`
// 	Amount uint   `bson:"amount"`
// }

// type User struct {
// 	ID           primitive.ObjectID `bson:"_id,omitempty"`
// 	Name         string             `bson:"name"`
// 	ProfileImage string             `bson:"profileImage"`
// 	Email        string             `bson:"email"`
// 	RegisterDate uint               `bson:"registerDate"`
// 	History      []UserHistory      `bson:"userHistory"`
// 	Stock        []UserStock        `bson:"userStock"`
// }

type Stock struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	StockImage string             `bson:"stockImage"`
	Name       string             `bson:"name"`
	Sign       string             `bson:"sign"`
	Price      uint               `bson:"price"`
	History    []StockHistory     `bson:"stockHistory"`
}

type User = model.UserAccount
type UserHistory = model.UserHistory
type UserStock = model.UserStock
type OrderRequest model.OrderRequest
type StockCollection = model.StockCollection
type StockHistory = model.StockHistory

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := "trading-system"

	client := InitMongoDB()
	redisClient := redis.InitRedis()

	// Get the database and collection
	db := client.Database(dbName)

	userRepositoryDB := repository.NewUserRepositoryDB(db.Collection("user"))
	stockRepositoryDB := repository.NewStockRepositoryDB(db.Collection("stock"))

	_ = service.NewUserService(userRepositoryDB, redisClient)
	_ = service.NewStockService(stockRepositoryDB)
	// objectId, err := primitive.ObjectIDFromHex("65c30de7b654c0e7bf938081")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// filter := bson.M{
	// 	"_id": objectId,
	// 	// "userStock.stockId": "65bf707e040d36a26f4bf523",
	//   // "userStock.amount": 201,
	// }

	// CREATE
	// account := model.CreateAccount{
	// 	Name:         "kongphop",
	// 	ProfileImage: "test",
	// 	Email:        "test@gmail.com",
	// }

	// result, err := userRepositoryDB.Create(account)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)

	// // DEPOSIT
	// result, err := userRepositoryDB.Deposit(
	// 	"65c8993c48096b5150cee5d6", 
	// 	5000,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// // WITHDRAW
	// result, err = userRepositoryDB.Withdraw(
	// 	"65c8993c48096b5150cee5d6", 
	// 	1000,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// SET FAVORITE
	// result, err := userRepositoryDB.SetFavorite(
	// 	"65c8993c48096b5150cee5d6",
	// 	"65c39a12c4e3672bcbf15b0f",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
	// fmt.Println(result)

	// GET BALANCE HISTORY
	// result, err := userRepositoryDB.GetBalanceHistory(
	// 	"65c4fa33835f044a5c8ed063", 
	// 	"WITHDRAW",
	// )
	// if err != nil {
	// 	log.Fatal()
	// }

	// fmt.Println(result)

	// BUY editd id
	// orderRequest := model.OrderRequest{
	// 	StockId:     "65c39a12c4e3672bcbf15b0f",
	// 	UserId:      "65c8993c48096b5150cee5d6",
	// 	Price:       60,
	// 	Amount:      8,
	// 	OrderType:   "auto",
	// 	OrderMethod: "buy",
	// }

	// result, err := userRepositoryDB.Buy(orderRequest)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)

	// SALE
	// orderRequest := model.OrderRequest{
	// 	StockId: "65bf707e040d36a26f4bf523",
	// 	UserId: "65c30de7b654c0e7bf938081",
	// 	Price: 10,
	// 	Amount: 100,
	// 	OrderType: "auto",
	// 	OrderMethod: "sale",
	// }

	// result, err := userRepositoryDB.Sale(orderRequest)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)

	// GET ACCOUNT
	// result, err := userRepositoryDB.GetAccount("65c8993c48096b5150cee5d6")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result.Favorite)

	// GET HISTORIES
	// result, err := userRepositoryDB.GetAllHistories("65c30de7b654c0e7bf938081")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET STOCK HISTORIES
	// stockId 65c35c9a832ed6ceda9a6b0f
	// result, err := userRepositoryDB.GetStockHistory(
	// 	"65c8993c48096b5150cee5d6",
	// 	"65c39a12c4e3672bcbf15b0f",
	// 	0,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET BALANCE
	// result, err := userRepositoryDB.GetBalance("65c4fa33835f044a5c8ed063")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET FAVORITE
	// _, err = userRepositoryDB.GetFavorite("65c4fa33835f044a5c8ed063")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// GET STOCK AMOUNT
	// result, err := userRepositoryDB.GetStockAmount(
	// 	"65c35bf2e383e16e2dd56362",
	// 	"65c35c9a832ed6ceda9a6b0f",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// DELETE FAVORITE STOCK
	// result, err := userRepositoryDB.DeleteFavorite(
	// 	"65c4fa33835f044a5c8ed063", 
	// 	"65c39a03dfb8060d99995934",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// DELETE ACCOUNT
	// result, err := userRepositoryDB.DeleteAccount("65c382330d619735e53b96fc")
	// if err != nil {
	// 	log.Fatal(result)
	// }

	// fmt.Println(result)

	// STOCK

	// CREATE STOCK COLLECTION
	// stockCollection := StockCollection{
	// 	StockImage: "test-image",
	// 	Name:       "TANAGUY",
	// 	Sign:       "TNG",
	// 	Price:      20,
	// 	History:    []StockHistory{},
	// }

	// result, err := stockRepositoryDB.CreateStock(stockCollection)
	// if err != nil {
	// 	log.Fatal(result)
	// }

	// fmt.Println(result)

	// GET ALL STOCK
	// result, err := stockRepositoryDB.GetAllStocks()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// CREATE STOCK ORDER
	// stockHistory := StockHistory{
	// 	ID: "65c39b189f5c807c54a53030",
	// 	Timestamp: uint(time.Now().Unix()),
	// 	Amount: 5,
	// 	Price: 13,
	// }
	// result, err := stockRepositoryDB.CreateStockOrder(
	// 	"65c39a12c4e3672bcbf15b0f",
	// 	stockHistory,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// EDIT NAME
	// result, err := stockRepositoryDB.EditName(
	// 	"65c39a03dfb8060d99995934",
	// 	"KONGPHOP",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// EDIT SIGN
	// result, err := stockRepositoryDB.EditSign(
	// 	"65c39a03dfb8060d99995934",
	// 	"KP",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// SET PRICE
	// result, err := stockRepositoryDB.SetPrice(
	// 	"65c39a03dfb8060d99995934",
	// 	11.11,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET FAVORITE STOCK
	// var stocks []string = []string{"65c39a03dfb8060d99995934", "65c39a12c4e3672bcbf15b0f"}
	// result, err := stockRepositoryDB.GetFavoriteStock(stocks)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET HISTORY STOCK
	result, err := stockRepositoryDB.GetStockHistory("65c39a03dfb8060d99995934")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	// GET STOCK
	// result, err := stockRepositoryDB.GetStock("65c39a03dfb8060d99995934")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET TOP STOCK
	// result, err := stockRepositoryDB.GetTopStocks()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

}

func InitMongoDB() *mongo.Client {
	// ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	uri := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("a")
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("b")
		panic(err)
	}

	return client
}

// objectId, err := primitive.ObjectIDFromHex("65c30de7b654c0e7bf938081")
// if err != nil {
// 	log.Fatal(err)
// }

// userHistory := UserHistory{
// 	Timestamp:   uint(time.Now().Unix()),
// 	StockId:     "65bf707e040d36a26f4bf523",
// 	Price:       101,
// 	Amount:      55,
// 	Status:      "pending",
// 	OrderType:   "auto",
// 	OrderMethod: "buy",
// }

// filter := bson.M{
// 	"_id": bson.M{
// 		"$eq": objectId,
// 	},
// 	"userStock.stockId": bson.M{
// 		"$eq": "65bf707e040d36a26f4bf523",
// 	},
// }
// update := bson.M{
// 	"$inc": bson.M{
// 		"userStock.$.amount": 55,
// 	},
// 	"$push": bson.M{
// 		"userHistory": userHistory,
// 	},
// } // Increment by 55

// result, err := userCol.UpdateOne(context.Background(), filter, update)
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Println(result.MatchedCount)

// func OrderMethodOperation(db *mongo.Collection, userId primitive.ObjectID, order UserHistory) (error) {
// 	ctx := context.Background()
// 	// amount := order.Amount
// 	// if order.OrderMethod == "sale" {
// 	// 	amount = -order.Amount
// 	// }

// 	filter := bson.M{
// 		"_id": userId,
// 	}
// 	var userAccount UserAccount
// 	err := db.FindOne(ctx, filter).Decode(&userAccount)
// 	if err != nil {
// 		return err
// 	}

// 	var userStock model.UserStock
// 	for _, stock := range userAccount.Stock {
// 		if stock.StockId == order.StockId {
// 			userStock = stock
// 			break
// 		}
// 	}

// 	if userStock.Amount == order.Amount {
// 		update := bson.M{
// 			"$pull": bson.M{
// 				"userStock": bson.M{"stockId": order.StockId},
// 			},
// 		}

// 		_, err := db.UpdateOne(ctx, filter, update)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		filter := bson.M{
// 			"_id": bson.M{
// 				"$eq": userId,
// 			},
// 			"userStock.stockId": bson.M{
// 				"$eq": order.StockId,
// 			},
// 		}
// 		update := bson.M{
// 			"$push": bson.M{
// 				"userHistory": order,
// 			},
// 			"$inc": bson.M{
// 				"userStock.$.amount": amount,
// 			},
// 		}

// 		ctx := context.Background()
// 		_, err := db.UpdateOne(ctx, filter, update)
// 		if err != nil {
// 			return err
// 		}

// 	}

// 	return nil
// }
