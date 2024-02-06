package util

import (
	"context"
	"errors"
	"server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserAccount = model.UserAccount
type UserHistory = model.UserHistory

func OrderMethodOperation(db *mongo.Collection, userId primitive.ObjectID, order UserHistory) (error) {
	ctx := context.Background()
	amount := order.Amount
	if order.OrderMethod == "sale" {
		amount = -order.Amount

		filter := bson.M{
			"_id": userId,
		}
		var userAccount UserAccount
		err := db.FindOne(ctx, filter).Decode(&userAccount)
		if err != nil {
			return err
		}
	
		var userStock model.UserStock
		for _, stock := range userAccount.Stock {
			if stock.StockId == order.StockId {
				userStock = stock
				break
			}
		}

		if userStock.Amount == order.Amount {
			update := bson.M{
				"$push": bson.M{
					"userHistory": order,
				},
				"$pull": bson.M{
					"userStock": bson.M{"stockId": order.StockId},
				},
			}
	
			_, err := db.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		} else {
			filter := bson.M{
				"_id": bson.M{
					"$eq": userId,
				}, 
				"userStock.stockId": bson.M{
					"$eq": order.StockId,
				},
			}
			update := bson.M{
				"$push": bson.M{
					"userHistory": order,
				},
				"$inc": bson.M{
					"userStock.$.amount": amount,
				},
			}
		
			ctx := context.Background()
			_, err := db.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}

			return nil
		}
	} else if order.OrderMethod == "buy" {
		filter := bson.M{
			"_id": bson.M{
				"$eq": userId,
			}, 
			"userStock.stockId": bson.M{
				"$eq": order.StockId,
			},
		}
		update := bson.M{
			"$push": bson.M{
				"userHistory": order,
			},
			"$inc": bson.M{
				"userStock.$.amount": amount,
			},
		}
	
		_, err := db.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

	}

	return errors.New("invalid order method")
}