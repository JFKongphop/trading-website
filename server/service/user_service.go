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
	balanceKey := fmt.Sprintf("balance:%s", userId)
	stockAmountKey := fmt.Sprintf("stockAmount:%s:%s", userId, stockId)
	userHistoryKey := fmt.Sprintf("userHistory:%s", userId)
	userStockHistoryKey := fmt.Sprintf("userStockHistory:%s", userId)
	s.redisClient.Del(ctx, stockAmountKey, userHistoryKey, userStockHistoryKey, balanceKey)

	return message, nil
}

func (s userService) SaleStock(orderRequest OrderRequest) (message string, err error) {
	message, err = s.userRepo.Sale(orderRequest)
	if err != nil {
		return "", err
	}

	userId := orderRequest.UserId 
	stockId := orderRequest.StockId
	balanceKey := fmt.Sprintf("balance:%s", userId)
	stockAmountKey := fmt.Sprintf("stockAmount:%s:%s", userId, stockId)
	userHistoryKey := fmt.Sprintf("userHistory:%s", userId)
	userStockHistoryKey := fmt.Sprintf("userStockHistory:%s", userId)
	s.redisClient.Del(ctx, stockAmountKey, userHistoryKey, userStockHistoryKey, balanceKey)

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

	if balanceJson, err := s.redisClient.Get(ctx, balanceKey).Result(); err == nil {
		if json.Unmarshal([]byte(balanceJson), &balance) == nil {
			return balance, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(ctx, balanceKey, string(data), time.Second*3600)
	}

	return result, nil
}

func (s userService) GetFavoriteStock(userId string) (favoriteStocks []string, err error) {
	favoriteKey := fmt.Sprintf("favorite:%s", userId)
	result, err := s.userRepo.GetFavorite(userId)
	if err != nil {
		return []string{}, err
	}

	if balanceJson, err := s.redisClient.Get(ctx, favoriteKey).Result(); err == nil {
		if json.Unmarshal([]byte(balanceJson), &favoriteStocks) == nil {
			return favoriteStocks, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(ctx, favoriteKey, string(data), time.Second*3600)
	}

	return result, nil
}

func (s userService) GetUserAccount(userId string) (userResponse UserResponse, err error) {
	userKey := fmt.Sprintf("user:%s", userId)
	result, err := s.userRepo.GetAccount(userId)
	if err != nil {
		return UserResponse{}, err
	}

	if userResponseJson, err := s.redisClient.Get(ctx, userKey).Result(); err == nil {
		if json.Unmarshal([]byte(userResponseJson), &userResponse) == nil {
			return userResponse, nil
		}
	}

	userResponse = UserResponse{
		Name:         result.Name,
		ProfileImage: result.ProfileImage,
		Email:        result.Email,
	}

	if data, err := json.Marshal(userResponse); err == nil {
		s.redisClient.Set(ctx, userKey, string(data), time.Second*3600)
	}

	return userResponse, nil
}

func (s userService) GetUserTradingHistories(userId string, startPage uint) (userHistories []ResponseUserHistory, err error) {
	result, err := s.userRepo.GetAllHistories(userId, startPage)
	if err != nil {
		return []ResponseUserHistory{}, err
	}

	if startPage == 0 {
		userHistoryKey := fmt.Sprintf("userHistory:%s", userId)
		if historyJson, err := s.redisClient.Get(ctx, userHistoryKey).Result(); err == nil {
			if json.Unmarshal([]byte(historyJson), &userHistories) == nil {
				return userHistories, nil
			}
		}
	
		if data, err := json.Marshal(result); err == nil {
			s.redisClient.Set(ctx, userHistoryKey, string(data), time.Second*3600)
		}
	}

	fmt.Println("db all")
	return result, nil
}

func (s userService) GetUserStockHistory(userId string, stockId string, startPage uint) (userStockHistories []ResponseUserHistory, err error) {
	userStockHistoryKey := fmt.Sprintf("userStockHistory:%s", userId)
	result, err := s.userRepo.GetUserStockHistory(userId, stockId, startPage)
	if err != nil {
		return []ResponseUserHistory{}, err
	}

	if startPage == 0 {
		if stockHistoryJson, err := s.redisClient.Get(ctx, userStockHistoryKey).Result(); err == nil {
			if json.Unmarshal([]byte(stockHistoryJson), &userStockHistories) == nil {
				return userStockHistories, nil
			}
		}
	
		if data, err := json.Marshal(result); err == nil {
			s.redisClient.Set(ctx, userStockHistoryKey, string(data), time.Second*3600)
		}
	}

	return result, nil
}

func (s userService) GetUserStockAmount(userId string, stockId string) (userStock UserStock, err error) {
	stockAmountKey := fmt.Sprintf("stockAmount:%s:%s", userId, stockId)
	result, err := s.userRepo.GetStockAmount(userId, stockId)
	if err != nil {
		return UserStock{}, err
	}

	if userStockJson, err := s.redisClient.Get(ctx, stockAmountKey).Result(); err == nil {
		if json.Unmarshal([]byte(userStockJson), &userStock) == nil {
 			return userStock, nil
		}
	}

	if data, err := json.Marshal(result); err == nil {
		s.redisClient.Set(ctx, stockAmountKey, string(data), time.Second*3600)
	}

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
