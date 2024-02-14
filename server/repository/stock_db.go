package repository

import (
	"server/errs"
	"server/model"
	"sort"
	"time"

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

type StockGroup = model.StockGroup

var (
	ErrPrice = errs.ErrPrice
	ErrSign  = errs.ErrSign
)

func NewStockRepositoryDB(db *mongo.Collection) StockRepository {
	return stockRepositoryDB{db}
}

func (r stockRepositoryDB) CreateStock(stockCollection StockCollection) (string, error) {
	if len(stockCollection.StockImage) == 0 ||
		len(stockCollection.Sign) == 0 ||
		len(stockCollection.Name) == 0 ||
		stockCollection.Price < 1 {
		return "", ErrData
	}

	_, err := r.db.InsertOne(ctx, stockCollection)
	if err != nil {
		return "", err
	}

	return "Successfullt created stock collection", nil
}

func (r stockRepositoryDB) CreateStockOrder(stockId string, stockOrder StockHistory) (string, error) {
	if len(stockId) == 0 {
		return "", ErrInvalidStock
	}

	if len(stockOrder.ID) == 0 ||
		stockOrder.Amount == 0 ||
		stockOrder.Price == 0 {
		return "", ErrData
	}

	stockOrder.Timestamp = time.Now().Unix()

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

	return "Successfully created stock order", nil
}

func (r stockRepositoryDB) GetAllStocks() ([]StockCollectionResponse, error) {
	filter := bson.M{}
	projection := bson.M{
		"_id":        1,
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
			ID:         result["_id"].(string),
			StockImage: result["stockImage"].(string),
			Name:       result["name"].(string),
			Sign:       result["sign"].(string),
			Price:      result["price"].(float64),
		}

		stockCollections = append(stockCollections, stockCollection)
	}

	return stockCollections, nil
}

func (r stockRepositoryDB) GetTopStocks() ([]StockGroup, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$unwind", Value: "$stockHistory"}},
		bson.D{{Key: "$sort", Value: bson.D{{
			Key: "stockHistory.timestamp", Value: -1,
		}}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":          "$_id",
			"name":         bson.M{"$first": "$name"},
			"price":        bson.M{"$first": "$price"},
			"sign":         bson.M{"$first": "$sign"},
			"stockImage":   bson.M{"$first": "$stockImage"},
			"stockHistory": bson.M{"$push": "$stockHistory"},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":        1,
			"name":       1,
			"price":      1,
			"sign":       1,
			"stockImage": 1,
			"stockHistory": bson.M{
				"$slice": []interface{}{
					"$stockHistory", 10,
				},
			},
		}}},
	}

	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return []StockGroup{}, err
	}
	defer cursor.Close(ctx)

	var stocks []StockGroup
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return []StockGroup{}, err
		}

		histories := result["stockHistory"].(bson.A)

		var volume float64 = 0
		for _, history := range histories {
			historyDoc := history.(bson.M)
			amount := historyDoc["amount"].(float64)
			price := historyDoc["price"].(float64)
			volume += amount * price
		}

		stock := StockGroup{
			ID:         result["_id"].(string),
			Name:       result["name"].(string),
			Price:      result["price"].(float64),
			Sign:       result["sign"].(string),
			StockImage: result["stockImage"].(string),
			Volume:     volume,
		}

		stocks = append(stocks, stock)
		volume = 0
	}

	sort.Slice(stocks[:], func(i, j int) bool {
		return stocks[i].Volume > stocks[j].Volume
	})

	amountOfStock := len(stocks)
	if amountOfStock > 10 {
		amountOfStock = 10
	}

	var topTenStock = stocks[0:amountOfStock]

	return topTenStock, nil
}

func (r stockRepositoryDB) GetStock(stockId string) (StockCollectionResponse, error) {
	if len(stockId) == 0 {
		return StockCollectionResponse{}, ErrInvalidStock
	}

	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return StockCollectionResponse{}, err
	}

	filter := bson.M{
		"_id": objectStockId,
	}
	projection := bson.M{
		"_id":        1,
		"stockImage": 1,
		"name":       1,
		"sign":       1,
		"price":      1,
	}

	var stockCollection StockCollectionResponse
	opts := options.FindOne().SetProjection(projection)
	err = r.db.FindOne(ctx, filter, opts).Decode(&stockCollection)
	if err != nil {
		return StockCollectionResponse{}, nil
	}

	return stockCollection, nil
}

func (r stockRepositoryDB) GetFavoriteStock(favoriteStockIds []string) ([]StockCollectionResponse, error) {
	var objectFavoriteStocks []primitive.ObjectID
	for _, stock := range favoriteStockIds {
		if len(stock) == 0 {
			return []StockCollectionResponse{}, ErrInvalidStock
		}

		objectStockId, err := primitive.ObjectIDFromHex(stock)
		if err != nil {
			return []StockCollectionResponse{}, err
		}
		objectFavoriteStocks = append(objectFavoriteStocks, objectStockId)
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": objectFavoriteStocks,
		},
	}
	projection := bson.M{
		"name":       1,
		"sign":       1,
		"price":      1,
		"stockImage": 1,
	}

	opts := options.Find().SetProjection(projection)
	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return []StockCollectionResponse{}, err
	}
	defer cursor.Close(ctx)

	var favoriteStocks []StockCollectionResponse
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return []StockCollectionResponse{}, err
		}

		favoriteStock := StockCollectionResponse{
			Name:       result["name"].(string),
			Sign:       result["sign"].(string),
			StockImage: result["stockImage"].(string),
			Price:      result["price"].(float64),
		}

		favoriteStocks = append(favoriteStocks, favoriteStock)
	}

	return favoriteStocks, nil
}

func (r stockRepositoryDB) GetStockHistory(stockId string) ([]StockHistoryResponse, error) {
	if len(stockId) == 0 {
		return []StockHistoryResponse{}, ErrInvalidStock
	}

	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return []StockHistoryResponse{}, err
	}

	filter := bson.M{
		"_id": objectStockId,
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$unwind", Value: "$stockHistory"}},
		bson.D{{Key: "$sort", Value: bson.D{{
			Key: "stockHistory.timestamp", Value: -1,
		}}}},
		bson.D{{Key: "$limit", Value: 2}},
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
			Price:  stockHistoryMap["price"].(float64),
			Amount: stockHistoryMap["amount"].(float64),
		}

		stockHistories = append(stockHistories, history)
	}

	return stockHistories, nil
}

func (r stockRepositoryDB) SetPrice(stockId string, price float64) (string, error) {
	if len(stockId) == 0 {
		return "", ErrInvalidStock
	}

	if price < 1 {
		return "", ErrPrice
	}

	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
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
	if len(stockId) == 0 {
		return "", ErrInvalidStock
	}

	if len(name) == 0 {
		return "", errs.ErrName
	}

	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
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

	return "Successfully updated name", nil
}

func (r stockRepositoryDB) EditSign(stockId string, sign string) (string, error) {
	if len(stockId) == 0 {
		return "", ErrInvalidStock
	}

	if len(sign) == 0 {
		return "", ErrSign
	}

	objectStockId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return "", err
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

	return "Successfully updated sign", nil
}

func (r stockRepositoryDB) DeleteStock(stockId string) (string, error) {
	if len(stockId) == 0 {
		return "", ErrInvalidStock
	}

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

// idValues := []string{"id1", "id2", "id3"} // Replace with your actual _id values

// // Construct the filter to find documents by _id values
// filter := bson.M{"_id": bson.M{"$in": idValues}}

// // Perform the find operation
// cursor, err := collection.Find(context.Background(), filter)
// if err != nil {
// 	log.Fatal(err)
// }
// defer cursor.Close(context.Background())
