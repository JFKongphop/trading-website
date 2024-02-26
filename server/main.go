package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"server/config"
	"server/handler"
	"server/model"
	"server/redis"
	"server/repository"
	"server/service"
	wshandler "server/ws-handler"

	// "server/specific"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"

	// "github.com/gofiber/fiber/v2"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	"go.mongodb.org/mongo-driver/bson"
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

type FavoriteStock struct {
	Favorite []string `bson:"favorite"`
	UID string `bson:"uid"`
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

var ctx = context.Background()
var uploader *model.ClientUploader

func main() {
	// app := fiber.New(fiber.Config{
	// 	Prefork: true,
	// })
	// app.Use(cors.New())
	// app.Use(logger.New())

	app := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	app.Use(cors.Default())
	app.Use(gin.Logger())
	app.Use(func(c *gin.Context) {		
		c.Set("uid", "test12345")
		c.Next()
	})

	hub := wshandler.H
	go hub.Run()


	firebase, err := config.InitializeFirebase();
	if err != nil {
		log.Fatal(err)
	}

	_, err = firebase.Auth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	

	// app.Use(func (c *fiber.Ctx) error {
	// 	// path := c.Path()

	// 	// if strings.Contains(path, "signin") || strings.Contains(path, "signup") {
	// 	// 	return c.Next()
	// 	// }
		
	// 	// authorization := c.GET("Authorization")
	// 	// if len(authorization) == 0 {
	// 	// 	log.Fatalf("error token: %v\n", err)
	// 	// }
	// 	// authToken := strings.Split(authorization, " ")[1]
	// 	// token, err := client.VerifyIDToken(ctx, authToken)
	// 	// if err != nil {
	// 	// 	log.Fatalf("error verifying ID token: %v\n", err)
	// 	// }

	// 	// c.Locals("uid", token.UID)

	// 	// mockup set 
	// 	c.Locals("uid", "MuwWsOQmD3PPRuMOlXh6SUbEVtn2")

	// 	return c.Next()
	// })

	mongoDB := initMongoDB()
	redisClient := redis.InitRedis()
	initTimeZone()



	db := mongoDB.Database(os.Getenv("MONGO_DATABASE"))
	userCollectionName := os.Getenv("MONGO_COLLECTION_USER")
	stockCollectionName := os.Getenv("MONGO_COLLECTION_STOCK")
	userCollection := db.Collection(userCollectionName)
	stockCollection := db.Collection(stockCollectionName)

	userRepositoryDB := repository.NewUserRepositoryDB(userCollection)
	stockRepositoryDB := repository.NewStockRepositoryDB(stockCollection)

	userService := service.NewUserService(userRepositoryDB, redisClient)
	stockService := service.NewStockService(stockRepositoryDB, redisClient, uploader)

	// var result model.UserAccount
	// userCollection.FindOne(ctx, bson.M{"uid": "MuwWsOQmD3PPRuMOlXh6SUbEVtn2"}).Decode(&result)
	// fmt.Println(result)

	// result, err := userRepositoryDB.GetAccount("MuwWsOQmD3PPRuMOlXh6SUbEVtn2")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// result, err := userRepositoryDB.DeleteAccount("MuwWsOQmD3PPRuMOlXh6SUbEVtn2")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	userHandler := handler.NewUserHandler(userService, stockService)
	stockHandler := handler.NewStockHandler(stockService)

	stockWebsocket := wshandler.NewStockWebsocket(stockService)

	// apiV1 := app.Group("/api/v1", func(c *fiber.Ctx) error {
	// 	c.Set("version1", "v1")
	// 	return c.Next()
	// })

	apiV1 := app.Group("/api/v1")

	userGroup := apiV1.Group("/user")
	stockGroup := apiV1.Group("/stock")

	userGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "api is running in user",
		})
	})

	stockGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "api is running in stock",
		})
	})

	websocketGroup := app.Group("/ws/v1")

	websocketGroup.GET("/price", func(c *gin.Context) {
		stockWebsocket.ServePriceWs(hub, c.Writer, c.Request)
	})

	websocketGroup.GET("/transaction", func(c *gin.Context) {
		stockWebsocket.ServeTransactionWs(hub, c.Writer, c.Request)
	})
	

	app.DELETE("/stock-history/:stockId", clearStocKHistory)



	// userGroup := apiV1.Group("/user", func(c *fiber.Ctx) error {
	// 	c.Set("user", "user")
	// 	return c.Next()
	// })

	// stockGroup := apiV1.Group("/stock", func(c *fiber.Ctx) error {
	// 	c.Set("stock", "stock")
	// 	return c.Next()
	// })


	userGroup.POST("/signup", userHandler.SignUp)
	userGroup.POST("/deposit", userHandler.DepositBalance)
	userGroup.POST("/withdraw", userHandler.WithdrawBalance)
	userGroup.POST("/buy", userHandler.BuyStock)
	userGroup.POST("/sale", userHandler.SaleStock)
	userGroup.POST("/set-favorite", userHandler.SetFavoriteStock)
	userGroup.GET("/balance-transaction", userHandler.GetUserBalanceHistory)
	userGroup.GET("/balance", userHandler.GetUserBalance)
	userGroup.GET("/get-favorite", userHandler.GetUserFavoriteStock)
	userGroup.POST("/signin", userHandler.SignIn)
	userGroup.GET("/trade-transaction", userHandler.GetUserTradingHistories)
	userGroup.GET("/stock-transaction", userHandler.GetUserStockHistory)
	userGroup.GET("/stock-ratio", userHandler.GetUserStockAmount)
	userGroup.DELETE("/delete-favorite", userHandler.DeleteFavoriteStock)
	userGroup.DELETE("/delete-account", userHandler.DeleteUserAccount)

	stockGroup.POST("/create-stock", stockHandler.CreateStockCollection)
	stockGroup.POST("/create-order/:stockId", stockHandler.CreateStockOrder)
	// stockGroup.GET("/collections", stockHandler.GetAllStockCollections)
	// stockGroup.GET("/top-stocks", stockHandler.GetTop10Stocks)
	// stockGroup.GET("/collection/:stockId", stockHandler.GetStockCollection)
	stockGroup.GET("/transaction/:stockId", stockHandler.GetStockHistory)
	stockGroup.GET("/price/:stockId", stockHandler.GetStockPrice)
	stockGroup.POST("/set-price/:stockId", stockHandler.SetStockPrice)
	// stockGroup.POST("/edit-name/:stockId", stockHandler.EditStockName)
	// stockGroup.POST("/edit-sign/:stockId", stockHandler.EditStockSign)
	// stockGroup.DELETE("/delete/:stockId", stockHandler.DeleteStockCollection)

	// stockGroup.GET("/test", func(c *fiber.Ctx) error {
	// 	// 
	// 	// fmt.Println(c.GET("Authorization"), c.Locals("uid"))
	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 		"message": "token test",
	// 	})
	// })

	// stockGroup.POST("/test", func(c *fiber.Ctx) error {
	// 	authorization := c.GET("Authorization")
	// 	authToken := strings.Split(authorization, " ")[0]
		

	// 	fmt.Println(authToken)

	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 		"message": "success",
	// 	})
	// })

	// stockGroup.GET("/", func(c *fiber.Ctx) error {
	// 	fmt.Println("test")
	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 		"message": "stock is running",
	// 	})
	// })
	// stockGroup.POST("/create-stock", stockHandler.CreateStockCollection)

	// excludeStockIds := []string{"65cc5fd45aa71b64fbb551a9", "65cc5fff0ca63a9e1e8b4db6"}
	// specific.DeleteExceptId(excludeStockIds, db.Collection("stock"))

	// excludeUserIds := []string{"65c8993c48096b5150cee5d6"}
	// specific.DeleteExceptId(excludeUserIds, db.Collection("user"))

	// var objectStockIdTest = []primitive.ObjectID{}
	// for _, id := range excludeIDs {
	// 	objectStockId, err := primitive.ObjectIDFromHex(id)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	objectStockIdTest = append(objectStockIdTest, objectStockId)
	// }
	// filter := bson.M{"_id": bson.M{"$nin": objectStockIdTest}}
	// db.Collection("stock").DeleteMany(context.Background(), filter)

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
	// 	UID:           "MuwWsOQmD3PPRuMOlXh6SUbEVtn2",
	// 	Name:         "JFKongphop",
	// 	ProfileImage: "test",
	// 	Email:        "kongphopleo@gmail.com",
	// }

	// result, err := userRepositoryDB.Create(account)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// result, err := userService.CreateUserAccount(account)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)

	// // DEPOSIT
	// result, err := userService.DepositBalance(
	// 	"65c8993c48096b5150cee5d6",
	// 	1,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// // WITHDRAW
	// result, err := userService.WithdrawBalance(
	// 	"65c8993c48096b5150cee5d6",
	// 	1,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// SET FAVORITE
	// result, err := userRepositoryDB.SetFavorite(
	// 	"MuwWsOQmD3PPRuMOlXh6SUbEVtn2",
	// 	"65c39a12c4e3672bcbf15b0f",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// projection := bson.M{
	// 	"favorite": 1,
	// 	"uid": 1,
	// }
	// var favoriteStock FavoriteStock
	// opts := options.FindOne().SetProjection(projection)
	// userCollection.FindOne(ctx, bson.M{"uid": "MuwWsOQmD3PPRuMOlXh6SUbEVtn2"}, opts).Decode(&favoriteStock)
	// fmt.Println(favoriteStock)

	// update := bson.M{
	// 	"$push": bson.M{
	// 		"favorite": "tettet",
	// 	},
	// }
	// a, _ := userCollection.UpdateOne(ctx, bson.M{"uid": "MuwWsOQmD3PPRuMOlXh6SUbEVtn2"}, update)
	// fmt.Println(a)




	// GET BALANCE HISTORY
	// result, err := userRepositoryDB.GetBalanceHistory(
	// 	"MuwWsOQmD3PPRuMOlXh6SUbEVtn2",
	// 	"DEPOSIT",
	// 	0,
	// )
	// if err != nil {
	// 	log.Fatal()
	// }

	// fmt.Println(result)

	// BUY editd id
	// orderRequest := model.OrderRequest{
	// 	StockId:     "65c39a03dfb8060d99995934",
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
	// result, err := userService.GetUserAccount("65c8993c48096b5150cee5d6")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(result)

	// GET HISTORIES
	// result, err := userService.GetUserTradingHistories("65c8993c48096b5150cee5d6", 0)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET STOCK HISTORIES
	// stockId 65c35c9a832ed6ceda9a6b0f
	// result, err := userService.GetUserStockHistory(
	// 	"65c8993c48096b5150cee5d6",
	// 	"65c39a12c4e3672bcbf15b0f",
	// 	1,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET BALANCE
	// result, err := userService.GetBalance("65c4fa33835f044a5c8ed063")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// result, err := userService.GetUserBalance("65c8993c48096b5150cee5d6")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET FAVORITE
	// result, err := userService.GetFavoriteStock("65c8993c48096b5150cee5d6")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET STOCK AMOUNT
	// result, err := userService.GetUserStockAmount(
	// 	"65c8993c48096b5150cee5d6",
	// 	"65c39a12c4e3672bcbf15b0f",
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
	// 	Name:       "test",
	// 	Sign:       "test",
	// 	Price:      20,
	// 	History:    []StockHistory{},
	// }

	// result, err := stockService.CreateStockCollection(stockCollection)
	// if err != nil {
	// 	log.Fatal(result)
	// }

	// fmt.Println(result)

	// GET ALL STOCK
	// result, err := stockService.GetAllStockCollections()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result[2])

	// CREATE STOCK ORDER
	// stockHistory := StockHistory{
	// 	ID: "65c8993c48096b5150cee5d6",
	// 	Timestamp: int64(time.Now().Unix()),
	// 	Amount: 1,
	// 	Price: 12,
	// }
	// result, err := stockService.CreateStockOrder(
	// 	"65ccda6623a24436ee69d21f",
	// 	stockHistory,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// EDIT NAME
	// result, err := stockService.EditStockName(
	// 	"65ccda6623a24436ee69d21f",
	// 	"TEST",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// EDIT SIGN
	// result, err := stockService.EditStockSign(
	// 	"65ccda6623a24436ee69d21f",
	// 	"T",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// SET PRICE
	// result, err := stockService.SetStockPrice(
	// 	"65ccda6623a24436ee69d21f",
	// 	11.11,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET FAVORITE STOCK
	// var stocks []string = []string{"65ccda6623a24436ee69d21f"}
	// result, err := stockService.GetFavoriteStock(stocks)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET HISTORY STOCK
	// result, err := stockService.GetStockHistory("65ccda6623a24436ee69d21f")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET STOCK
	// result, err := stockService.GetStockCollection("65ccda6623a24436ee69d21f")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// GET TOP STOCK
	// result, err := stockService.GetTop10Stocks()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	// result, err := stockService.DeleteStockCollection("65ccddc168d59a209e58cc82")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(result)

	app.Run(":4000")
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleCloudProjectId := os.Getenv("GOOGLE_STORAGE_PROJECT_ID")
	googleCloudBucketName := os.Getenv("GOOGLE_STORAGE_BUCKET_NAME")
	googleCloudPrivateKeyId := os.Getenv("GOOGLE_STORAGE_PRIVATE_KEY_ID")
	googleCloudPrivateKey := os.Getenv("GOOGLE_STORAGE_PRIVATE_KEY")
	googleCloudClientEmail := os.Getenv("GOOGLE_STORAGE_CLIENT_EMAIL")
	googleCloudClientId := os.Getenv("GOOGLE_STORAGE_CLIENT_ID")
	googleCloudCertUri := os.Getenv("GOOGLE_STORAGE_CERT_URI")
	googleCloudFolder := os.Getenv("GOOGLE_STORAGE_FOLDER")

	keyFile := map[string]interface{}{
		"type":                        "service_account",
		"project_id":                  googleCloudProjectId,
		"private_key_id":              googleCloudPrivateKeyId,
		"private_key":                 googleCloudPrivateKey,
		"client_email":                googleCloudClientEmail,
		"client_id":                   googleCloudClientId,
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        googleCloudCertUri,
		"universe_domain":             "googleapis.com",
	}

	jsonBytes, err := json.Marshal(keyFile)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(jsonBytes)))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uploader = &model.ClientUploader{
		Cl:         client,
		BucketName: googleCloudBucketName,
		ProjectID:  googleCloudProjectId,
		UploadPath: googleCloudFolder,
	}
}

func initMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	uri := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Fail to connect mongodb")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Fail to check ping")
	}

	return client
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func clearStocKHistory(c *gin.Context) {
	mongoDB := initMongoDB()
	db := mongoDB.Database(os.Getenv("MONGO_DATABASE"))
	stockCollectionName := os.Getenv("MONGO_COLLECTION_STOCK")
	stockCollection := db.Collection(stockCollectionName)
	
	objectStockId, _ := primitive.ObjectIDFromHex(c.Param("stockId"))

	filter := bson.M{
		"_id": objectStockId,
	}
	update := bson.M{
		"$set": bson.M{
			"stockHistory": []StockHistory{},
		},
	}

	stockCollection.UpdateOne(ctx, filter, update)
}