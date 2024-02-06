package util

import (
	"context"
	"server/model"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStock = model.UserStock

func CheckValidStock(db *mongo.Collection, userId primitive.ObjectID, stockId string) (bool, UserStock, error) {
	filter := bson.M{
		"_id": userId,
		"userStock": bson.M{
			"$elemMatch": bson.M{
				"stockId": stockId,
			},
		},
	}
	
	validStock := true
	var userAccount UserAccount
	ctx := context.Background()
	err := db.FindOne(ctx, filter).Decode(&userAccount)
	if err != nil {
		messageError := err.Error()
		if strings.Contains(messageError, "no documents in result") {
			validStock = false
		} else {
			return validStock, UserStock{}, nil
		}
	}

	var userStock UserStock
	for _, stock := range userAccount.Stock {
		if stock.StockId == stockId {
			userStock = stock
			break
		}
	}

	return validStock, userStock, nil
}