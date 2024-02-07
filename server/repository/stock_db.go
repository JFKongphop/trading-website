package repository

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type stockRepositoryDB struct {
	db *mongo.Collection
}

type StockPrice struct {
	Price float64 `bson:"price"`
}

func NewStockRepositoryDB(db *mongo.Collection) StockRepository {
	return stockRepositoryDB{db}
}

func (r stockRepositoryDB) CreateStock(stockCollection StockCollection) (string, error) {
	_, err := r.db.InsertOne(ctx, stockCollection)
	if err != nil {
		return "", err
	}

	return "Successfullt created stock collection", nil
}

func (r stockRepositoryDB) CreateStockOrder(stockId string, stockOrder StockHistory) (string, error) {
	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": objectStockId,
	}
	projection := bson.M{
		"price": 1,
	}

	var stockPrice StockPrice
	opts := options.FindOne().SetProjection(projection)
	err = r.db.FindOne(ctx, filter, opts).Decode(&stockPrice)
	if err != nil {
		return "", err
	}

	stockOrder.Price = stockPrice.Price

	update := bson.M{
		"$push": bson.M{
			"stockHistory": stockOrder,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully creted stock order", nil
}

func (r stockRepositoryDB) GetAllStocks() ([]StockCollectionResponse, error) {
	filter := bson.M{}
	projection := bson.M{
		"stockImage": 1,
		"name":       1,
		"sign":       1,
		"price":      1,
	}

	opts := options.Find().SetProjection(projection)
	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return []StockCollectionResponse{}, err
	}

	var stockCollections = []StockCollectionResponse{}
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return []StockCollectionResponse{}, err
		}

		stockCollection := StockCollectionResponse{
			StockImage: result["stockImage"].(string),
			Name:       result["name"].(string),
			Sign:       result["sign"].(string),
			Price:      result["price"].(float64),
		}

		stockCollections = append(stockCollections, stockCollection)
	}

	return stockCollections, nil
}

func (r stockRepositoryDB) GetTopStocks() ([]StockCollection, error) {

	return []StockCollection{}, nil
}

func (r stockRepositoryDB) GetStock(string) (StockCollection, error) {

	return StockCollection{}, nil
}

func (r stockRepositoryDB) GetStockHistory(stockId string) ([]StockHistoryResponse, error) {
	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return []StockHistoryResponse{}, nil
	}

	filter := bson.M{
		"_id": objectStockId,
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$unwind", Value: "$stockHistory"}},
		bson.D{{Key: "$limit", Value: 2}},
		bson.D{{Key: "$sort", Value: bson.D{{
			Key: "stockHistory.timestamp", Value: -1,
		}}}},
		bson.D{{Key: "$project", Value: bson.M{
			"stockHistory.price":  1,
			"stockHistory.amount": 1,
		}}},
	}

	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return []StockHistoryResponse{}, nil
	}
	defer cursor.Close(ctx)

	var stockHistories []StockHistoryResponse
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return []StockHistoryResponse{}, err
		}

		stockHistoryMap := result["stockHistory"].(bson.M)
		history := StockHistoryResponse{
			Price: stockHistoryMap["price"].(float64),
			Amount: stockHistoryMap["amount"].(float64), 
		}

		stockHistories = append(stockHistories, history)
	}

	return stockHistories, nil
}

func (r stockRepositoryDB) SetPrice(stockId string, price float64) (string, error) {
	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
	}

	if price <= 0 {
		return "", errors.New("invalid price")
	}

	filter := bson.M{
		"_id": objectStockId,
	}
	update := bson.M{
		"$set": bson.M{
			"price": price,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully set price", nil
}

func (r stockRepositoryDB) EditName(stockId string, name string) (string, error) {
	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
	}

	if len(name) == 0 {
		return "", errors.New("invalid name")
	}

	filter := bson.M{
		"_id": objectStockId,
	}
	update := bson.M{
		"$set": bson.M{
			"name": name,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully update name", nil
}

func (r stockRepositoryDB) EditSign(stockId string, sign string) (string, error) {
	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
	}

	if len(sign) == 0 {
		return "", errors.New("invalid sign")
	}

	filter := bson.M{
		"_id": objectStockId,
	}
	update := bson.M{
		"$set": bson.M{
			"sign": sign,
		},
	}

	_, err = r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return "Successfully update sign", nil
}

func (r stockRepositoryDB) DeleteStock(stockId string) (string, error) {
	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
	}

	filter := bson.M{
		"_id": objectStockId,
	}

	_, err = r.db.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}

	return "Successfully deleted stock", nil
}
