package service

import (
	"context"
	"encoding/json"
	"fmt"
	"server/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRepository = repository.UserRepository

type userService struct {
	userRepo    UserRepository
	redisClient *redis.Client
}

var ctx = context.Background()

func NewUserService(userRepo UserRepository, redisClient *redis.Client) UserService {
	return userService{userRepo, redisClient}
}

func (s userService) CreateUserAccount(userAccount CreateAccount) (string, error) {
	result, err := s.userRepo.Create(userAccount)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s userService) DepositBalance(userId string, depositMoney float64) (message string, err error) {
	message, err = s.userRepo.Deposit(userId, depositMoney)
	if err != nil {
		return "", err
	}

	balanceKey := fmt.Sprintf("balance:%s", userId)
	s.redisClient.Del(ctx, balanceKey)

	return message, nil
}

func (s userService) WithdrawBalance(userId string, withdrawMoney float64) (message string, err error) {
	message, err = s.userRepo.Withdraw(userId, withdrawMoney)
	if err != nil {
		return "", err
	}

	balanceKey := fmt.Sprintf("balance:%s", userId)
	s.redisClient.Del(ctx, balanceKey)

	return message, nil
}

func (s userService) BuyStock(orderRequest OrderRequest) (message string, err error) {
	message, err = s.userRepo.Buy(orderRequest)
	if err != nil {
		return "", err
	}

	userId := orderRequest.UserId 
	stockId := orderRequest.StockId
	stockAmountKey := fmt.Sprintf("stockAmount:%s:%s", userId, stockId)
	userHistoryKey := fmt.Sprintf("userHistory:%s", userId)
	s.redisClient.Del(ctx, stockAmountKey, userHistoryKey)

	return message, nil
}

func (s userService) SaleStock(orderRequest OrderRequest) (message string, err error) {
	message, err = s.userRepo.Sale(orderRequest)
	if err != nil {
		return "", err
	}

	userId := orderRequest.UserId 
	stockId := orderRequest.StockId
	stockAmountKey := fmt.Sprintf("stockAmount:%s:%s", userId, stockId)
	userHistoryKey := fmt.Sprintf("userHistory:%s", userId)
	s.redisClient.Del(ctx, stockAmountKey, userHistoryKey)

	return message, nil
}

func (s userService) SetFavoriteStock(userId string, stockId string) (message string, err error) {
	message, err = s.userRepo.SetFavorite(userId, stockId)
	if err != nil {
		return "", err
	}

	favoriteKey := fmt.Sprintf("favorite:%s", userId)
	s.redisClient.Del(ctx, favoriteKey)

	return message, nil
}

func (s userService) GetUserBalanceHistory(userId string, method string, perpage uint) (balanceHistories []BalanceHistory, err error) {
	balanceHistories, err = s.userRepo.GetBalanceHistory(userId, method, perpage)
	if err != nil {
		return []BalanceHistory{}, err
	}

	return balanceHistories, nil
}

func (s userService) GetUserBalance(userId string) (balance float64, err error) {
	balanceKey := fmt.Sprintf("balance:%s", userId)
	result, err := s.userRepo.GetBalance(userId)
	if err != nil {
		return 0, err
	}

	if balanceJson, err := s.redisClient.Get(context.Background(), balanceKey).Result(); err == nil {
		if json.Unmarshal([]byte(balanceJson), &balance) == nil {
			fmt.Println("redis")
			return balance, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(context.Background(), balanceKey, string(data), time.Second*3600)
	}

	fmt.Println("db")
	return result, nil
}

func (s userService) GetFavoriteStock(userId string) (favoriteStocks []string, err error) {
	favoriteKey := fmt.Sprintf("favorite:%s", userId)
	result, err := s.userRepo.GetFavorite(userId)
	if err != nil {
		return []string{}, err
	}

	if balanceJson, err := s.redisClient.Get(context.Background(), favoriteKey).Result(); err == nil {
		if json.Unmarshal([]byte(balanceJson), &favoriteStocks) == nil {
			fmt.Println("redis")
			return favoriteStocks, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(context.Background(), favoriteKey, string(data), time.Second*3600)
	}

	fmt.Println("db")
	return result, nil
}

func (s userService) GetUserAccount(userId string) (userResponse UserResponse, err error) {
	userKey := fmt.Sprintf("user:%s", userId)
	result, err := s.userRepo.GetAccount(userId)
	if err != nil {
		return UserResponse{}, err
	}

	if userResponseJson, err := s.redisClient.Get(context.Background(), userKey).Result(); err == nil {
		if json.Unmarshal([]byte(userResponseJson), &userResponse) == nil {
			fmt.Println("redis")
			return userResponse, nil
		}
	}

	userResponse = UserResponse{
		Name:         result.Name,
		ProfileImage: result.ProfileImage,
		Email:        result.Email,
	}

	if data, err := json.Marshal(userResponse); err == nil {
		s.redisClient.Set(context.Background(), userKey, string(data), time.Second*3600)
	}

	fmt.Println("db")
	return userResponse, nil
}

// will change
func (s userService) GetUserTradingHistories(userId string, startPage uint) (userHistories []ResponseUserHistory, err error) {
	userHistoryKey := fmt.Sprintf("userHistory:%s", userId)
	result, err := s.userRepo.GetAllHistories(userId, startPage)
	if err != nil {
		return []ResponseUserHistory{}, err
	}

	if historyJson, err := s.redisClient.Get(context.Background(), userHistoryKey).Result(); err == nil {
		if json.Unmarshal([]byte(historyJson), &userHistories) == nil {
			fmt.Println("redis")
			return userHistories, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(context.Background(), userHistoryKey, string(data), time.Second*3600)
	}

	return result, nil
}

// will change
func (s userService) GetUserStockHistory(userId string, stockId string, start uint) (userStockHistories []ResponseUserHistory, err error) {
	userStockHistoryKey := fmt.Sprintf("userStockHistory:%s", userId)
	result, err := s.userRepo.GetUserStockHistory(userId, stockId, start)
	if err != nil {
		return []ResponseUserHistory{}, err
	}

	if stockHistoryJson, err := s.redisClient.Get(context.Background(), userStockHistoryKey).Result(); err == nil {
		if json.Unmarshal([]byte(stockHistoryJson), &userStockHistories) == nil {
			fmt.Println("redis")
			return userStockHistories, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(context.Background(), userStockHistoryKey, string(data), time.Second*3600)
	}

	return result, nil
}

func (s userService) GetUserStockAmount(userId string, stockId string) (userStock UserStock, err error) {
	stockAmountKey := fmt.Sprintf("stockAmount:%s:%s", userId, stockId)
	result, err := s.userRepo.GetStockAmount(userId, stockId)
	if err != nil {
		return UserStock{}, err
	}

	if userStockJson, err := s.redisClient.Get(context.Background(), stockAmountKey).Result(); err == nil {
		if json.Unmarshal([]byte(userStockJson), &userStock) == nil {
			fmt.Println("redis")
			return userStock, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(context.Background(), stockAmountKey, string(data), time.Second*3600)
	}

	fmt.Println("db")
	return result, nil
}

func (s userService) DeleteFavoriteStock(userId string, stockId string) (message string, err error) {
	message, err = s.userRepo.DeleteFavorite(userId, stockId)
	if err != nil {
		return "", err
	}

	favoriteKey := fmt.Sprintf("favorite:%s", userId)
	s.redisClient.Del(ctx, favoriteKey)

	return message, nil
}

func (s userService) DeleteUserAccount(userId string) (message string, err error) {
	message, err = s.userRepo.DeleteAccount(userId)
	if err != nil {
		return "", err
	}

	userKey := fmt.Sprintf("user:%s", userId)
	s.redisClient.Del(ctx, userKey)

	return message, nil
}
