package service

import (
	"encoding/json"
	"fmt"
	"server/model"
	"server/repository"
	"server/util"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

type StockRepository = repository.StockRepository
type StockGraph = model.StockGraph

type stockService struct {
	stockRepo   StockRepository
	redisClient *redis.Client
	googleCloudUpload *model.ClientUploader
}


func NewStockService(stockRepo StockRepository, redisClient *redis.Client, googleCloudUpload *model.ClientUploader) StockService {
	return stockService{stockRepo, redisClient, googleCloudUpload}
}

func (s stockService) CreateStockCollection(stockCollection StockCollectionRequest) (message string, err error) {
	err = util.UploadFile(
		stockCollection.StockImage, 
		stockCollection.Name, 
		s.googleCloudUpload,
	)
	if err != nil {
		return "", err
	}

	stockCollectionsKey := "stockCollections"
	stock := StockCollection{
		StockImage: stockCollection.Name,
		Name:       stockCollection.Name,
		Sign:       stockCollection.Sign,
		Price:      stockCollection.Price,
		History:    []StockHistory{},
	}
	message, err = s.stockRepo.CreateStock(stock)
	if err != nil {
		return "", err
	}

	s.redisClient.Del(ctx, stockCollectionsKey)

	return message, nil
}

func (s stockService) CreateStockOrder(stockId string, stockOrder StockHistory) (message string, err error) {
	message, err = s.stockRepo.CreateStockOrder(stockId, stockOrder)
	if err != nil {
		return "", err
	}

	return message, nil
}

func (s stockService) GetAllStockCollections() (stockCollections []StockCollectionResponse, err error) {
	stockCollectionsKey := "stockCollections"
	result, err := s.stockRepo.GetAllStocks()
	if err != nil {
		return []StockCollectionResponse{}, err
	}

	if stockCollectionJson, err := s.redisClient.Get(ctx, stockCollectionsKey).Result(); err == nil {
		if json.Unmarshal([]byte(stockCollectionJson), &stockCollections) == nil {
			return stockCollections, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(ctx, stockCollectionsKey, string(data), time.Second*3600)
	}

	return result, nil
}

func (s stockService) GetTop10Stocks() (top10Stock []TopStock, err error) {
	top10StockKey := "top10Stock"
	result, err := s.stockRepo.GetTopStocks()
	if err != nil {
		return []TopStock{}, err
	}

	if stockCollectionJson, err := s.redisClient.Get(ctx, top10StockKey).Result(); err == nil {
		if json.Unmarshal([]byte(stockCollectionJson), &top10Stock) == nil {
			return top10Stock, nil
		}
	}

	for _, stock := range result {
		topStock := TopStock{
			ID:    stock.ID,
			Sign:  stock.Sign,
			Price: stock.Price,
		}

		top10Stock = append(top10Stock, topStock)
	}

	if data, err := json.Marshal(top10Stock); err == nil {
		s.redisClient.Set(ctx, top10StockKey, string(data), time.Second*3600)
	}

	return top10Stock, nil
}

func (s stockService) GetStockCollection(stockId string) (stockCollection StockCollectionResponse, err error) {
	stockCollectionKey := fmt.Sprintf("stockCollection:%s", stockId)
	result, err := s.stockRepo.GetStock(stockId)
	if err != nil {
		return StockCollectionResponse{}, err
	}

	if stockCollectionJson, err := s.redisClient.Get(ctx, stockCollectionKey).Result(); err == nil {
		if json.Unmarshal([]byte(stockCollectionJson), &stockCollection) == nil {
			return stockCollection, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(ctx, stockCollectionKey, string(data), time.Second*3600)
	}

	return result, nil
}

func (s stockService) GetFavoriteStock(favoriteStockIds []string) (favoriteStocks []StockCollectionResponse, err error) {
	favoriteStocks, err = s.stockRepo.GetFavoriteStock(favoriteStockIds)
	if err != nil {
		return []StockCollectionResponse{}, err
	}

	return favoriteStocks, nil
}

func (s stockService) GetStockHistory(stockId string) (stockHistories []StockHistoryResponse, err error) {
	stockHistories, err = s.stockRepo.GetStockHistory(stockId)
	if err != nil {
		return []StockHistoryResponse{}, err
	}

	return stockHistories, nil
}

func (s stockService) GetStockPrice(stockId string) (price float64, err error) {
	price, err = s.stockRepo.GetPrice(stockId)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func (s stockService) GetStockGraph(stockId string) (graph []Graph, err error) {
	stockGraph, err := s.stockRepo.GetGraph(stockId)
	if err != nil {
		return []Graph{}, err
	}

	groupedData := make(map[int64][]StockGraph)

	for _, item := range stockGraph {
		interval := item.Timestamp / 5 * 5
		groupedData[interval] = append(groupedData[interval], item)
	}

	var keys []int64
	for key := range groupedData {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, timestamp := range keys {
		items := groupedData[timestamp]
		fmt.Printf("Timestamp: %d\n", timestamp)
		min := items[0].Price
		max := items[0].Price
		for i, item := range items {
			if items[i].Price < min {
				min = item.Price
			}
			if items[i].Price > max {
				max = item.Price
			}
			// fmt.Println("Price:", item.Price, "timestamp:", item.Timestamp)
		}

		// Y -> [open, max, min, close]
		graph = append(graph, Graph{
			X: timestamp,
			Y: []float64{
				items[len(items) - 1].Price,
				max,
				min,
				items[0].Price,
			},
		})
	}

	return graph, nil
}


func (s stockService) SetStockPrice(stockId string, price float64) (message string, err error) {
	message, err = s.stockRepo.SetPrice(stockId, price)
	if err != nil {
		return "", err
	}
	
	return message, nil
}

func (s stockService) EditStockName(stockId string, name string) (message string, err error) {
	stockCollectionsKey := "stockCollections"
	message, err = s.stockRepo.EditName(stockId, name)
	if err != nil {
		return "", err
	}

	s.redisClient.Del(ctx, stockCollectionsKey)
	
	return message, nil
}

func (s stockService) EditStockSign(stockId string, sign string) (message string, err error) {
	stockCollectionsKey := "stockCollections"
	message, err = s.stockRepo.EditSign(stockId, sign)
	if err != nil {
		return "", err
	}

	s.redisClient.Del(ctx, stockCollectionsKey)
	
	return message, nil
}

func (s stockService) DeleteStockCollection(stockId string) (message string, err error) {
	stockCollectionsKey := "stockCollections"
	message, err = s.stockRepo.DeleteStock(stockId)
	if err != nil {
		return "", err
	}

	s.redisClient.Del(ctx, stockCollectionsKey)
	
	return message, nil
}
