package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	// "math/rand"
	"os"
	"time"

	"server/config"
	"server/handler"
	"server/model"
	"server/redis"
	"server/repository"
	"server/service"
	"server/ws-handler"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
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
	UID      string   `bson:"uid"`
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
type StockCollection = model.StockCollection
type StockHistory = model.StockHistory
type StockOrder = model.StockHistory

var ctx = context.Background()
var uploader *model.ClientUploader

func main() {
	app := gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	app.Use(cors.Default())
	// app.Use(gin.Logger())
	// app.Use(func(c *gin.Context) {
	// 	c.Set("uid", "test12345")
	// 	c.Next()
	// })

	hub := wshandler.H
	go hub.Run()

	firebase, err := config.InitializeFirebase()
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

	// ClearStocKHistory()
	// for i := 0; i < 200; i++ {
	// 	a := time.Duration(i * 12 * int(time.Minute))
	// 	rand.Seed(time.Now().UnixNano())

	// 	stockOrder := StockHistory{
	// 		ID: "",
	// 		Amount: float64(rand.Intn(10) + 1),
	// 		Price: float64(rand.Intn(61) + 40),
	// 		Timestamp: time.Now().Add(a).Unix(),
	// 	}

	// 	update := bson.M{
	// 		"$push": bson.M{
	// 			"stockHistory": stockOrder,
	// 		},
	// 	}

	// 	objectStockId, _ := primitive.ObjectIDFromHex("65d60a2dc25b2ff14700a3c2")
	// 	stockCollection.UpdateOne(ctx, bson.M{
	// 		"_id": objectStockId,
	// 	}, update)
	// }

	// fmt.Println(graph)

	userHandler := handler.NewUserHandler(userService, stockService)
	stockHandler := handler.NewStockHandler(stockService)
	stockWebsocket := wshandler.NewStockWebsocket(stockService)

	apiV1 := app.Group("/api/v1")

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	userGroup := apiV1.Group("/user")
	stockGroup := apiV1.Group("/stock")

	websocketGroup := app.Group("/ws/v1")

	websocketGroup.GET("/price", func(c *gin.Context) {
		stockWebsocket.ServePriceWs(hub, c.Writer, c.Request)
	})
	websocketGroup.GET("/transaction", func(c *gin.Context) {
		stockWebsocket.ServeTransactionWs(hub, c.Writer, c.Request)
	})
	websocketGroup.GET("/graph", func(c *gin.Context) {
		stockWebsocket.ServeGraphWs(hub, c.Writer, c.Request)
	})

	// app.DELETE("/stock-history/:stockId", ClearStocKHistory)

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
	stockGroup.GET("/collections", stockHandler.GetAllStockCollections)
	stockGroup.GET("/top-stocks", stockHandler.GetTop10Stocks)
	stockGroup.GET("/collection/:stockId", stockHandler.GetStockCollection)
	stockGroup.GET("/transaction/:stockId", stockHandler.GetStockHistory)
	stockGroup.GET("/price/:stockId", stockHandler.GetStockPrice)
	stockGroup.GET("/graph/:stockId", stockHandler.GetStockGraph)
	stockGroup.POST("/set-price/:stockId", stockHandler.SetStockPrice)
	stockGroup.POST("/edit-name/:stockId", stockHandler.EditStockName)
	stockGroup.POST("/edit-sign/:stockId", stockHandler.EditStockSign)
	stockGroup.DELETE("/delete/:stockId", stockHandler.DeleteStockCollection)

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

func ClearStocKHistory( /*c *gin.Context*/ ) {
	mongoDB := initMongoDB()
	db := mongoDB.Database(os.Getenv("MONGO_DATABASE"))
	stockCollectionName := os.Getenv("MONGO_COLLECTION_STOCK")
	stockCollection := db.Collection(stockCollectionName)

	objectStockId, _ := primitive.ObjectIDFromHex("65d60a2dc25b2ff14700a3c2")

	filter := bson.M{
		"_id": objectStockId,
	}
	update := bson.M{
		"$set": bson.M{
			"stockHistory": []StockHistory{},
		},
	}

	stockCollection.UpdateOne(ctx, filter, update)

	fmt.Println("done")
}
