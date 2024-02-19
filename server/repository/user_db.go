package repository

import (
	"context"
	"server/errs"
	"server/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
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

type UserBalanceHistory struct {
	BalanceHistory []BalanceHistory `bson:"balanceHistory"`
}

type Stock struct {
	UserStock []UserStock `bson:"userStock"`
}

type UserBalance struct {
	Balance float64 `bson:"balance"`
}

type UserFavorite struct {
	Favorite []string `bson:"favorite"`
}

type ValidStock struct {
	Stock []string `bson:"favorite"`
}

var (
	ErrSignin          = errs.ErrSignin
	ErrUser           = errs.ErrUser
	ErrData           = errs.ErrData
	ErrMoney          = errs.ErrMoney
	ErrBalance        = errs.ErrBalance
	ErrOrderType      = errs.ErrOrderType
	ErrOrderMethod    = errs.ErrOrderMethod
	ErrInvalidStock   = errs.ErrInvalidStock
	ErrFavoriteStock  = errs.ErrFavoriteStock
	ErrNotEnoughStock = errs.ErrNotEnoughStock
)

func NewUserRepositoryDB(db *mongo.Collection) UserRepository {
	return userRepositoryDB{db}
}

func (r userRepositoryDB) Create(data CreateAccount) (string, error) {
	name := data.Name
	profileImage := data.ProfileImage
	email := data.Email
	uid := data.UID
	if len(name) == 0 || len(profileImage) == 0 || len(email) == 0 {
		return "", ErrData
	}

	filter := bson.M{
		"uid": uid,
	}
	var result UserAccount
	if err := r.db.FindOne(ctx, filter).Decode(&result); err != nil {
		return "", ErrUser
	}

	if len(result.Email) > 0 {
		return "", ErrSignin
	}

	user := UserAccount{
		UID:            uid,
		Name:           name,
		ProfileImage:   profileImage,
		Email:          email,
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
	userId := orderRequest.UserId
	stockId := orderRequest.StockId
	amount := orderRequest.Amount
	price := orderRequest.Price
	orderType := orderRequest.OrderType
	OrderMethod := orderRequest.OrderMethod

	if len(userId) == 0 {
		return "", ErrUser
	}

	if (len(stockId) == 0) ||
		(len(orderType) == 0) ||
		(len(OrderMethod) == 0) ||
		(amount <= 0) ||
		(price <= 0) {
		return "", ErrData
	}

	if orderType != "auto" && orderType != "order" {
		return "", ErrOrderType
	}

	if OrderMethod != "buy" {
		return "", ErrOrderMethod
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return "", err
	// }

	userHistory := UserHistory{
		StockId:     stockId,
		Price:       price,
		Amount:      amount,
		Status:      "pending",
		Timestamp:   int64(time.Now().Unix()),
		OrderType:   orderRequest.OrderType,
		OrderMethod: orderRequest.OrderMethod,
	}

	validStock, _, balance, err := util.CheckValidStock(r.db, userId, stockId)
	if err != nil {
		return "", err
	}

	stockValue := price * amount
	if stockValue > balance {
		return "", ErrBalance
	}

	if validStock {
		filter := bson.M{
			"uid": bson.M{
				"$eq": userId,
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
			"uid": userId,
		}

		var data bson.M
		err = r.db.FindOneAndUpdate(ctx, filter, update).Decode(&data)
		if err != nil {

			return "", err
		}
	}

	return "Successfully bought stock", nil
}

func (r userRepositoryDB) Sale(orderRequest OrderRequest) (string, error) {
	userId := orderRequest.UserId
	stockId := orderRequest.StockId
	amount := orderRequest.Amount
	price := orderRequest.Price
	orderType := orderRequest.OrderType
	OrderMethod := orderRequest.OrderMethod

	if len(userId) == 0 {
		return "", ErrUser
	}

	if (len(stockId) == 0) ||
		(len(orderType) == 0) ||
		(len(OrderMethod) == 0) ||
		(amount <= 0) ||
		(price <= 0) {
		return "", ErrData
	}

	if orderType != "auto" && orderType != "order" {
		return "", ErrOrderType
	}

	if OrderMethod != "sale" {
		return "", ErrOrderMethod
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return "", err
	// }

	userHistory := UserHistory{
		StockId:     stockId,
		Price:       price,
		Amount:      amount,
		Status:      "pending",
		Timestamp:   int64(time.Now().Unix()),
		OrderType:   orderRequest.OrderType,
		OrderMethod: orderRequest.OrderMethod,
	}

	validStock, userStock, _, err := util.CheckValidStock(r.db, userId, stockId)
	if err != nil {
		return "", err
	}

	if amount > userStock.Amount {
		return "", errs.ErrNotEnoughStock
	}

	stockValue := price * amount
	if validStock {
		if userStock.Amount == amount {
			filter := bson.M{
				"uid": userId,
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
				"uid": bson.M{
					"$eq": userId,
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
			return "", ErrNotEnoughStock
		}
	} else {
		return "", ErrInvalidStock
	}

	return "Successfully sold stock", nil
}

func (r userRepositoryDB) SetFavorite(userId string, stockId string) (string, error) {
	if len(userId) == 0 {
		return "", ErrUser
	}

	if len(stockId) == 0 {
		return "", ErrInvalidStock
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return "", err
	// }

	filter := bson.M{
		"uid": userId,
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$unwind", Value: "$favorite"}},
		bson.D{{Key: "$match", Value: bson.M{
			"favorite": stockId,
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"favorite": 1,
		}}},
	}

	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	validStock := true
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return "", err
		}

		if len(result["favorite"].(string)) == 0 {
			validStock = false
		}
	}

	if validStock {
		return "", ErrFavoriteStock
	}

	update := bson.M{
		"$push": bson.M{
			"favorite": stockId,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully set favorite stock", nil
}

func (r userRepositoryDB) GetBalanceHistory(userId string, method string, skip uint) ([]BalanceHistory, error) {
	if len(userId) == 0 {
		return []BalanceHistory{}, ErrUser
	}

	if len(method) == 0 {
		return []BalanceHistory{}, ErrOrderMethod
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return []BalanceHistory{}, err
	// }

	filter := bson.M{
		"uid": userId,
	}
	var pipeline mongo.Pipeline
	if method == "ALL" {
		pipeline = mongo.Pipeline{
			bson.D{{Key: "$match", Value: filter}},
			bson.D{{Key: "$unwind", Value: "$balanceHistory"}},
			bson.D{{Key: "$sort", Value: bson.D{
				{Key: "balanceHistory.timestamp", Value: -1},
			}}},
			bson.D{{Key: "$skip", Value: skip}},
			bson.D{{Key: "$limit", Value: 10}},
			bson.D{{Key: "$project", Value: bson.M{
				"balanceHistory": 1,
			}}},
		}
	} else if method == "DEPOSIT" || method == "WITHDRAW" {
		pipeline = mongo.Pipeline{
			bson.D{{Key: "$match", Value: filter}},
			bson.D{{Key: "$unwind", Value: "$balanceHistory"}},
			bson.D{{Key: "$match", Value: bson.M{
				"balanceHistory.method": method,
			}}},
			bson.D{{Key: "$sort", Value: bson.D{
				{Key: "balanceHistory.timestamp", Value: -1},
			}}},
			bson.D{{Key: "$skip", Value: skip}},
			bson.D{{Key: "$limit", Value: 10}},
			bson.D{{Key: "$project", Value: bson.M{
				"balanceHistory": 1,
			}}},
		}
	} else {
		return []BalanceHistory{}, ErrOrderMethod
	}

	var balanceHistories []BalanceHistory
	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return []BalanceHistory{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return []BalanceHistory{}, err
		}

		// fmt.Println(result)

		balanceHistoryMap := result["balanceHistory"].(bson.M)
		balanceHistory := BalanceHistory{
			Timestamp: balanceHistoryMap["timestamp"].(int64),
			Balance:   balanceHistoryMap["balance"].(float64),
			Method:    balanceHistoryMap["method"].(string),
		}

		balanceHistories = append(balanceHistories, balanceHistory)
	}

	return balanceHistories, nil
}

func (r userRepositoryDB) GetBalance(userId string) (float64, error) {
	if len(userId) == 0 {
		return 0, ErrUser
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return 0, err
	// }

	filter := bson.M{
		"uid": userId,
	}
	projection := bson.M{
		"balance": 1,
	}

	var userBalance UserBalance
	opts := options.FindOne().SetProjection(projection)
	err := r.db.FindOne(ctx, filter, opts).Decode(&userBalance)
	if err != nil {
		return 0, err
	}

	return userBalance.Balance, nil
}

func (r userRepositoryDB) GetFavorite(userId string) ([]string, error) {
	if len(userId) == 0 {
		return []string{}, ErrUser
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return []string{}, err
	// }

	filter := bson.M{
		"uid": userId,
	}
	projection := bson.M{
		"favorite": 1,
	}

	var userFavorite UserFavorite
	opts := options.FindOne().SetProjection(projection)
	err := r.db.FindOne(ctx, filter, opts).Decode(&userFavorite)
	if err != nil {
		return []string{}, err
	}

	return userFavorite.Favorite, nil
}

func (r userRepositoryDB) Deposit(userId string, depositMoney float64) (string, error) {
	if depositMoney <= 0 {
		return "", ErrMoney
	}

	if len(userId) == 0 {
		return "", ErrUser
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return "", err
	// }

	balanceHistory := BalanceHistory{
		Timestamp: int64(time.Now().Unix()),
		Balance:   depositMoney,
		Method:    "DEPOSIT",
	}

	filter := bson.M{
		"uid": userId,
	}
	update := bson.M{
		"$inc": bson.M{
			"balance": depositMoney,
		},
		"$push": bson.M{
			"balanceHistory": balanceHistory,
		},
	}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully deposited money", nil
}

func (r userRepositoryDB) Withdraw(userId string, withdrawMoney float64) (string, error) {
	if withdrawMoney <= 0 {
		return "", ErrMoney
	}

	if len(userId) == 0 {
		return "", ErrUser
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	fmt.Println("test")
	// 	return "", err
	// }

	balanceHistory := BalanceHistory{
		Timestamp: int64(time.Now().Unix()),
		Balance:   withdrawMoney,
		Method:    "WITHDRAW",
	}

	filter := bson.M{
		"uid": userId,
	}
	var userAccount UserAccount
	err := r.db.FindOne(ctx, filter).Decode(&userAccount)
	if err != nil {
		return "", err
	}

	if userAccount.Balance < withdrawMoney {
		return "", ErrBalance
	}

	update := bson.M{
		"$inc": bson.M{
			"balance": -withdrawMoney,
		},
		"$push": bson.M{
			"balanceHistory": balanceHistory,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully withdrawed money", nil
}

func (r userRepositoryDB) GetAccount(userId string) (userAccount UserAccount, err error) {
	if len(userId) == 0 {
		return UserAccount{}, ErrUser
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return UserAccount{}, err
	// }

	filter := bson.M{
		"uid": userId,
	}
	err = r.db.FindOne(ctx, filter).Decode(&userAccount)
	if err != nil {
		return UserAccount{}, err
	}

	return userAccount, nil
}

func (r userRepositoryDB) GetAllHistories(userId string, start uint) ([]UserHistory, error) {
	if len(userId) == 0 {
		return []UserHistory{}, ErrUser
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return []UserHistory{}, err
	// }

	stop := start + 10
	filter := bson.M{
		"uid": userId,
	}
	projection := bson.M{
		"userHistory": bson.M{
			"$slice": []int{int(start), int(stop)},
		},
	}

	var result History
	opts := options.FindOne().SetProjection(projection)
	err := r.db.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return []UserHistory{}, err
	}

	return result.UserHistory, nil
}

func (r userRepositoryDB) GetUserStockHistory(userId string, stockId string, skip uint) ([]UserHistory, error) {
	if len(userId) == 0 {
		return []UserHistory{}, ErrUser
	}

	if len(stockId) == 0 {
		return []UserHistory{}, ErrInvalidStock
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return []UserHistory{}, err
	// }

	filter := bson.M{
		"uid": userId,
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$unwind", Value: "$userHistory"}},
		bson.D{{Key: "$match", Value: bson.M{"userHistory.stockId": stockId}}},
		bson.D{{Key: "$sort", Value: bson.D{{
			Key: "userHistory.timestamp", Value: -1,
		}}}},
		// pagination
		bson.D{{Key: "$skip", Value: skip}},
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
			Timestamp:   userHistoryMap["timestamp"].(int64),
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
	if len(userId) == 0 {
		return UserStock{}, ErrUser
	}

	if len(stockId) == 0 {
		return UserStock{}, ErrInvalidStock
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return UserStock{}, err
	// }

	filter := bson.M{
		"uid": userId,
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$unwind", Value: "$userStock"}},
		bson.D{{Key: "$match", Value: bson.M{"userStock.stockId": stockId}}},
		bson.D{{Key: "$project", Value: bson.M{"userStock": 1}}},
	}

	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return UserStock{}, err
	}
	defer cursor.Close(ctx)

	var userStock UserStock
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return UserStock{}, err
		}

		userStockMap := result["userStock"].(bson.M)
		userStock.Amount = userStockMap["amount"].(float64)
		userStock.StockId = userStockMap["stockId"].(string)
	}

	if err := cursor.Err(); err != nil {
		return UserStock{}, err
	}

	return userStock, nil
}

func (r userRepositoryDB) DeleteFavorite(userId string, stockId string) (string, error) {
	if len(userId) == 0 {
		return "", ErrUser
	}

	if len(stockId) == 0 {
		return "", ErrInvalidStock
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return "", err
	// }

	filter := bson.M{
		"uid": userId,
	}
	delete := bson.M{
		"$pull": bson.M{
			"favorite": stockId,
		},
	}

	_, err := r.db.UpdateOne(ctx, filter, delete)
	if err != nil {
		return "", err
	}

	return "Successfully deleted favorite stock", nil
}

func (r userRepositoryDB) DeleteAccount(userId string) (string, error) {
	if len(userId) == 0 {
		return "", ErrUser
	}

	// objectUserId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return "", err
	// }

	filter := bson.M{
		"uid": userId,
	}

	_, err := r.db.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}

	return "Successfully deleted account", nil
}
