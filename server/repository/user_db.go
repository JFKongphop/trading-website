package repository

import (
	"context"
	"errors"
	"server/model"
	"server/util"

	// "strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()

type userRepositoryDB struct {
	db *mongo.Collection
}
type UserStock = model.UserStock

func NewUserRepositoryDB(db *mongo.Collection) UserRepository {
	return userRepositoryDB{db}
}

func (r userRepositoryDB) Create(data CreateAccount) (string, error) {
	user := UserAccount{
		Name:         data.Name,
		ProfileImage: data.ProfileImage,
		Email:        data.Email,
		RegisterDate: uint(time.Now().Unix()),
		History:      []UserHistory{},
		Stock:        []UserStock{},
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

	userHistory := UserHistory{
		Timestamp:   uint(time.Now().Unix()),
		StockId:     stockId,
		Price:       orderRequest.Price,
		Amount:      amount,
		Status:      "pending",
		OrderType:   orderRequest.OrderType,
		OrderMethod: orderRequest.OrderMethod,
	}

	validStock, _, err := util.CheckValidStock(r.db, objectId, stockId)
	if err != nil {
		return "", err
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

	userHistory := UserHistory{
		Timestamp:   uint(time.Now().Unix()),
		StockId:     stockId,
		Price:       orderRequest.Price,
		Amount:      amount,
		Status:      "pending",
		OrderType:   orderRequest.OrderType,
		OrderMethod: orderRequest.OrderMethod,
	}

	validStock, userStock,  err := util.CheckValidStock(r.db, objectId, stockId)
	if err != nil {
		return "", err
	}

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
			}
	
			_, err := r.db.UpdateOne(ctx, filter, update)
			if err != nil {
				return "Successfully sold stock", err
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
				},
			}
		
			ctx := context.Background()
			_, err := r.db.UpdateOne(ctx, filter, update)
			if err != nil {
				return "Successfully sold stock", err
			}
		} else {
			return "", errors.New("not enough stock")
		}
	} else {
		return "", errors.New("invalid stock")
	}

	return "Successfully sold stock", nil
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

func (r userRepositoryDB) GetAllHistories(id string) ([]UserHistory, error) {
	// objectId, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return []UserHistory{}, err
	// }

	return []UserHistory{}, nil
}

func (r userRepositoryDB) GetStockHistory(string) ([]UserHistory, error) {
	return []UserHistory{}, nil
}

func (r userRepositoryDB) DeleteAccount(string, string) (string, error) {
	return "", nil
}
