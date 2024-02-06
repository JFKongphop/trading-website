package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"server/model"
	"server/repository"
	// "time"

	// "time"

	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()


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

type StockHistory struct {
	ID        primitive.ObjectID `bson:"userId"`
	Timestamp uint               `bson:"timestamp"`
	Price     uint               `bson:"price"`
}

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := "trading-system"
	userColName := "user"
	stockColName := "stock"

	client := InitMongoDB()

	// Get the database and collection
	db := client.Database(dbName)
	userCol := db.Collection(userColName)
	stockCol := db.Collection(stockColName)

	userRepositoryDB := repository.NewUserRepositoryDB(userCol)
	_ = repository.NewStockRepositoryDB(stockCol)

	// objectId, err := primitive.ObjectIDFromHex("65bf707e040d36a26f4bf522")
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
	// 	Name:         "dsfeiugher",
	// 	ProfileImage: "goiregjegijor",
	// 	Email:        "awfwefwfw@gmail.com",
	// }

	// result, err := userRepositoryDB.Create(account)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)

	// GETACCOUNT
	// result, err := userRepositoryDB.GetAccount("65bf707e040d36a26f4bf522")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)

	// BUY
	// orderRequest := model.OrderRequest{
	// 	StockId: "65bf707e040d36a26f4bf523",
	// 	UserId: "65bf707e040d36a26f4bf522",
	// 	Price: 100,
	// 	Amount: 50,
	// 	OrderType: "auto", 
	// 	OrderMethod: "buy",
	// }

	// result, err := userRepositoryDB.Buy(orderRequest)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)	

	// SALE
	orderRequest := model.OrderRequest{
		StockId: "65bf707e040d36a26f4bf523",
		UserId: "65bf707e040d36a26f4bf522",
		Price: 100,
		Amount: 50,
		OrderType: "auto", 
		OrderMethod: "sale",
	}

	result, err := userRepositoryDB.Sale(orderRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)	

	// update := bson.M{
	// 	"$pull": bson.M{
	// 		"userStock": bson.M{"stockId": "65bf707e040d36a26f4bf523"},
	// 	},
	// }

	// result, err := userCol.UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)


}

func InitMongoDB() *mongo.Client {
	ctx := context.Background()
	uri := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("a")
		panic(err)
	}
	// defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("b")
		panic(err)
	}

	return client
}


// objectId, err := primitive.ObjectIDFromHex("65bf707e040d36a26f4bf522")
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