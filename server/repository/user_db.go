package repository

import (
	"context"
	"errors"
	"fmt"
	"server/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type userRepositoryDB struct {
	db *mongo.Collection
}

type History struct {
	UserHistory []UserHistory `bson:"userHistory"`
}

type Stock struct {
	UserStock []UserStock `bson:"userStock"`
}

func NewUserRepositoryDB(db *mongo.Collection) UserRepository {
	return userRepositoryDB{db}
}

func (r userRepositoryDB) Create(data CreateAccount) (string, error) {
	user := UserAccount{
		Name:           data.Name,
		ProfileImage:   data.ProfileImage,
		Email:          data.Email,
		Balance:        0,
		BalanceHistory: []BalanceHistory{},
		Favorite:       []string{},
		History:        []UserHistory{},
		Stock:          []UserStock{},
	}

	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return "Successfully created account", nil
}

func (r userRepositoryDB) Buy(orderRequest OrderRequest) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(orderRequest.UserId)
	if err != nil {
		return "", err
	}

	stockId := orderRequest.StockId
	amount := orderRequest.Amount
	price := orderRequest.Price

	userHistory := UserHistory{
		StockId:     stockId,
		Price:       price,
		Amount:      amount,
		Status:      "pending",
		Timestamp:   uint(time.Now().Unix()),
		OrderType:   orderRequest.OrderType,
		OrderMethod: orderRequest.OrderMethod,
	}

	validStock, _, balance, err := util.CheckValidStock(r.db, objectId, stockId)
	if err != nil {
		return "", err
	}

	stockValue := price * amount
	if stockValue > balance {
		return "", errors.New("balance is not enough")
	}

	if validStock {
		filter := bson.M{
			"_id": bson.M{
				"$eq": objectId,
			},
			"userStock.stockId": bson.M{
				"$eq": stockId,
			},
		}
		update := bson.M{
			"$push": bson.M{
				"userHistory": userHistory,
			},
			"$inc": bson.M{
				"userStock.$.amount": amount,
				"balance":            -stockValue,
			},
		}

		_, err := r.db.UpdateOne(ctx, filter, update)
		if err != nil {
			return "", err
		}
	} else {
		userStock := UserStock{
			StockId: stockId,
			Amount:  amount,
		}

		update := bson.M{
			"$push": bson.M{
				"userHistory": userHistory,
				"userStock":   userStock,
			},
			"$inc": bson.M{
				"balance": -stockValue,
			},
		}

		filter := bson.M{
			"_id": objectId,
		}

		var test bson.M
		err = r.db.FindOneAndUpdate(ctx, filter, update).Decode(&test)
		if err != nil {

			return "", err
		}
	}

	return "Successfully bought stock", nil
}

func (r userRepositoryDB) Sale(orderRequest OrderRequest) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(orderRequest.UserId)
	if err != nil {
		return "", err
	}

	stockId := orderRequest.StockId
	amount := orderRequest.Amount
	price := orderRequest.Price

	userHistory := UserHistory{
		StockId:     stockId,
		Price:       price,
		Amount:      amount,
		Status:      "pending",
		Timestamp:   uint(time.Now().Unix()),
		OrderType:   orderRequest.OrderType,
		OrderMethod: orderRequest.OrderMethod,
	}

	validStock, userStock, _, err := util.CheckValidStock(r.db, objectId, stockId)
	if err != nil {
		return "", err
	}

	if amount > userStock.Amount {
		return "", errors.New("stock not enough")
	}

	stockValue := price * amount
	if validStock {
		if userStock.Amount == amount {
			filter := bson.M{
				"_id": objectId,
			}
			update := bson.M{
				"$push": bson.M{
					"userHistory": userHistory,
				},
				"$pull": bson.M{
					"userStock": bson.M{"stockId": stockId},
				},
				"$inc": bson.M{
					"balance": stockValue,
				},
			}

			_, err := r.db.UpdateOne(ctx, filter, update)
			if err != nil {
				return "", err
			}
		} else if userStock.Amount > amount {
			filter := bson.M{
				"_id": bson.M{
					"$eq": objectId,
				},
				"userStock.stockId": bson.M{
					"$eq": stockId,
				},
			}
			update := bson.M{
				"$push": bson.M{
					"userHistory": userHistory,
				},
				"$inc": bson.M{
					"userStock.$.amount": -amount,
					"balance":            stockValue,
				},
			}

			_, err := r.db.UpdateOne(ctx, filter, update)
			if err != nil {
				return "", err
			}
		} else {
			return "", errors.New("not enough stock")
		}
	} else {
		return "", errors.New("invalid stock")
	}

	return "Successfully sold stock", nil
}

func (r userRepositoryDB) SetFavorite(string) (string, error)
func (r userRepositoryDB)  GetBalanceHistory(string) ([]BalanceHistory, error)
func (r userRepositoryDB) GetBalance(string) (float64, error)
func (r userRepositoryDB) GetFavorite(string) ([]StockCollectionResponse, error)

func (r userRepositoryDB) Deposit(userId string, depositMoney float64) (string, error) {
	if depositMoney <= 0 {
		return "", errors.New("invalid money")
	}

	if len(userId) == 0 {
		return "", errors.New("invalid user")
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": objectId,
	}
	update := bson.M{
		"$inc": bson.M{
			"balance": depositMoney,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully deposited money", nil
}

func (r userRepositoryDB) Withdraw(userId string, withdrawMoney float64) (string, error) {
	if withdrawMoney <= 0 {
		return "", errors.New("invalid money")
	}

	if len(userId) == 0 {
		return "", errors.New("invalid user")
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println("test")
		return "", err
	}

	filter := bson.M{
		"_id": objectId,
	}
	var userAccount UserAccount
	err = r.db.FindOne(ctx, filter).Decode(&userAccount)
	if err != nil {
		return "", err
	}

	if userAccount.Balance < withdrawMoney {
		return "", errors.New("balance not enough")
	}

	update := bson.M{
		"$inc": bson.M{
			"balance": -withdrawMoney,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully withdrawed money", nil
}

func (r userRepositoryDB) GetAccount(id string) (userAccount UserAccount, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return UserAccount{}, err
	}

	filter := bson.M{"_id": objectId}
	err = r.db.FindOne(ctx, filter).Decode(&userAccount)
	if err != nil {
		return UserAccount{}, err
	}

	return userAccount, nil
}

func (r userRepositoryDB) GetAllHistories(userId string) ([]UserHistory, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return []UserHistory{}, err
	}

	filter := bson.M{
		"_id": objectId,
	}
	projection := bson.M{
		"userHistory": bson.M{
			"$slice": []int{0, 1},
		},
	}

	var result History
	opts := options.FindOne().SetProjection(projection)
	err = r.db.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return []UserHistory{}, err
	}

	return result.UserHistory, nil
}

func (r userRepositoryDB) GetStockHistory(userId string, stockId string) ([]UserHistory, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return []UserHistory{}, err
	}

	filter := bson.M{
		"_id": objectId,
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$unwind", Value: "$userHistory"}},
		bson.D{{Key: "$match", Value: bson.M{"userHistory.stockId": stockId}}},
		// pagination
		// bson.D{{Key: "$skip", Value: 1}},
		bson.D{{Key: "$limit", Value: 10}},
		bson.D{{Key: "$project", Value: bson.M{"userHistory": 1}}},
	}

	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return []UserHistory{}, err
	}
	defer cursor.Close(ctx)

	var userHistories []UserHistory
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return []UserHistory{}, err
		}

		userHistoryMap := result["userHistory"].(bson.M)
		history := UserHistory{
			Price:       userHistoryMap["price"].(float64),
			Amount:      userHistoryMap["amount"].(float64),
			Status:      userHistoryMap["status"].(string),
			Timestamp:   userHistoryMap["timestamp"].(uint),
			OrderType:   userHistoryMap["orderType"].(string),
			OrderMethod: userHistoryMap["orderMethod"].(string),
		}

		userHistories = append(userHistories, history)
	}

	if err := cursor.Err(); err != nil {
		return []UserHistory{}, err
	}

	return userHistories, nil
}

func (r userRepositoryDB) GetStockAmount(userId string, stockId string) (UserStock, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return UserStock{}, err
	}

	filter := bson.M{
		"_id": objectId,
	}
	projection := bson.M{
		"userStock": bson.M{
			"$elemMatch": bson.M{
				"stockId": stockId,
			},
		},
	}

	var result Stock
	opts := options.FindOne().SetProjection(projection)
	err = r.db.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return UserStock{}, err
	}

	return result.UserStock[0], nil
}

func (r userRepositoryDB) DeleteFavorite(string) (string, error)

func (r userRepositoryDB) DeleteAccount(userId string) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": objectId,
	}

	_, err = r.db.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}

	return "Successfully deleted account", nil
}
