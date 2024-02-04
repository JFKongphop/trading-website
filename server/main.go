package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"server/repository"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
	City string             `bson:"city"`
}

type UserHistory struct {
	Timestamp uint   `bson:"timestamp"`
	Name      string `bson:"name"`
	Price     uint   `bson:"price"`
	Amount    uint   `bson:"amount"`
	Status    string `bson:"status"`
}

type UserStock struct {
	Name   string `bson:"name"`
	Amount uint   `bson:"amount"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	ProfileImage string             `bson:"profileImage"`
	Email        string             `bson:"email"`
	RegisterDate uint               `bson:"registerDate"`
	History      []UserHistory      `bson:"userHistory"`
	Stock        []UserStock        `bson:"userStock"`
}

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := "trading-system"
	userColName := "user"
	stockColName := "stock"
	// collName := "person"

	ctx := context.Background()

	client := InitMongoDB()

	// Get the database and collection
	db := client.Database(dbName)
	userCol := db.Collection(userColName)
	stockCol := db.Collection(stockColName)

	_ = repository.NewStockepositoryDB(userCol)
	_ = repository.NewStockepositoryDB(stockCol)


	user := User{
		Name:         "Kongphop",
		ProfileImage: "helloworld",
		Email:        "kongphopleo@gmail.com",
		RegisterDate: uint(time.Now().Unix()),
		History:      []UserHistory{},
		Stock:        []UserStock{},
	}

	userResult, err := userCol.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID", userResult.InsertedID)

	stock := Stock{
		StockImage: "helloworld",
		Name:     "KONGPHOP",
		Sign:     "KP",
		Price:    100,
		History:  []StockHistory{},
	}

	stockResult, err := stockCol.InsertOne(ctx, stock)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID", stockResult.InsertedID)

	// Create a person document
	// person := Person{Name: "John Doe", Age: 30, City: "New York"}

	// // Insert the document
	// result, err := col.InsertOne(ctx, person)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Inserted document with ID:", result.InsertedID)

	// // Find all documents
	// cursor, err := col.Find(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cursor.Close(ctx)

	// for cursor.Next(ctx) {
	// 	var p Person
	// 	err := cursor.Decode(&p)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("Found person:", p)
	// }

	// if err := cursor.Err(); err != nil {
	// 	log.Fatal(err)
	// }
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
