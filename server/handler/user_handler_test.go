package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/errs"
	"server/handler"
	"server/model"
	"server/service"
	"testing"

	"github.com/gin-gonic/gin"
)

type UserAccount = model.UserAccount
type CreateAccount = model.CreateAccount
type UserStock = model.UserStock
type UserBalanceRequest = model.UserBalanceRequest
type UserSetFavoriteRequest = model.UserSetFavoriteRequest
type OrderRequest = model.OrderRequest
type BalanceHistory = model.BalanceHistory
type StockCollectionResponse = model.StockCollectionResponse
type UserResponse = model.UserResponse
type UserIdRequest = model.UserIdRequest
type UserHistory = model.UserHistory

var (
	userId          = "test12345"
	ErrData         = errs.ErrData
	ErrUser         = errs.ErrUser
	ErrMoney        = errs.ErrMoney
	ErrOrderMethod  = errs.ErrOrderMethod
	ErrInvalidStock = errs.ErrInvalidStock
)

func userPath(route string) string {
	return fmt.Sprintf("/api/v1/user/%s", route)
}

func TestSignUp(t *testing.T) {
	expectedMessage := "Successfully created account"
	url := userPath("signup")

	t.Run("Successfully signup", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := CreateAccount{
			UID:          "123",
			Name:         "Test User",
			ProfileImage: "test.jpg",
			Email:        "test@example.com",
		}

		userService.
			On("CreateUserAccount", testBody).
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.POST(url, userHandler.SignUp)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

		t.Run("Error on service create user account", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := CreateAccount{
			UID:          "123",
			ProfileImage: "test.jpg",
			Email:        "test@example.com",
		}

		userService.
			On("CreateUserAccount", testBody).
			Return(expectedMessage, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.POST(url, userHandler.SignUp)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service create user account", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := CreateAccount{
			UID:          "123",
			ProfileImage: "test.jpg",
			Email:        "test@example.com",
		}

		userService.
			On("CreateUserAccount", testBody).
			Return(expectedMessage, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.POST(url, userHandler.SignUp)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestDepositBalance(t *testing.T) {
	expectedMessage := "Successfully deposited money"
	url := userPath("deposit")

	t.Run("Successfully deposit", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserBalanceRequest{
			Balance: 1000,
		}

		userService.
			On("DepositBalance", userId, testBody.Balance).
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.DepositBalance)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service deposit balance", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserBalanceRequest{
			Balance: -1,
		}

		userService.
			On("DepositBalance", userId, testBody.Balance).
			Return(expectedMessage, ErrMoney)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.DepositBalance)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrMoney.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestWithdrawBalance(t *testing.T) {
	expectedMessage := "Successfully withdrawed money"
	url := userPath("withdraw")

	t.Run("Successfully withdraw", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserBalanceRequest{
			Balance: 1000,
		}

		userService.
			On("WithdrawBalance", userId, testBody.Balance).
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.WithdrawBalance)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service withdraw balance", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserBalanceRequest{
			Balance: -1,
		}

		userService.
			On("WithdrawBalance", userId, testBody.Balance).
			Return(expectedMessage, ErrMoney)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.WithdrawBalance)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrMoney.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestBuyStock(t *testing.T) {
	expectedMessage := "Successfully bought stock"
	url := userPath("buy")

	t.Run("Successfully buy stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := OrderRequest{
			StockId:     "65c39a03dfb8060d99995934",
			UserId:      "test12345",
			Price:       60,
			Amount:      8,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		userService.
			On("BuyStock", testBody).
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.BuyStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service buy stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := OrderRequest{
			StockId:     "65c39a03dfb8060d99995934",
			UserId:      "test12345",
			Price:       60,
			Amount:      8,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		userService.
			On("BuyStock", testBody).
			Return(expectedMessage, ErrOrderMethod)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.BuyStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(`{"message":"%s"}`, ErrOrderMethod.Error())

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestSaleStock(t *testing.T) {
	expectedMessage := "Successfully sold stock"
	url := userPath("sale")

	t.Run("Successfully sale stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := OrderRequest{
			StockId:     "65c39a03dfb8060d99995934",
			UserId:      "test12345",
			Price:       60,
			Amount:      8,
			OrderType:   "auto",
			OrderMethod: "sale",
		}

		userService.
			On("SaleStock", testBody).
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.SaleStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service sale stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := OrderRequest{
			StockId:     "65c39a03dfb8060d99995934",
			UserId:      "test12345",
			Price:       60,
			Amount:      8,
			OrderType:   "auto",
			OrderMethod: "buy",
		}

		userService.
			On("SaleStock", testBody).
			Return(expectedMessage, ErrOrderMethod)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.SaleStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrOrderMethod.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestSetFavoriteStock(t *testing.T) {
	expectedMessage := "Successfully set favorite stock"
	url := userPath("set-favorite")

	t.Run("Successfully set favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserSetFavoriteRequest{
			StockId: "65c39a03dfb8060d99995934",
		}

		userService.
			On("SetFavoriteStock", userId, testBody.StockId).
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.SetFavoriteStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service set favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserSetFavoriteRequest{
			StockId: "",
		}

		userService.
			On("SetFavoriteStock", userId, testBody.StockId).
			Return(expectedMessage, ErrInvalidStock)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.SetFavoriteStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrInvalidStock.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestGetUserBalance(t *testing.T) {
	expectedMessage := "Successfully fetched user balance"
	expectedBalance := 1001
	url := userPath("balance")

	t.Run("Successfully buy stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserBalance", userId).
			Return(expectedBalance, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.GET(url, userHandler.GetUserBalance)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"balance":%d,"message":"%s"}`,
			expectedBalance,
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service buy stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserBalance", "").
			Return(expectedBalance, ErrUser)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "")
		})

		router.GET(url, userHandler.GetUserBalance)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrUser.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestGetUserBalanceHistory(t *testing.T) {
	expectedMessage := "Successfully fetched transaction balance"
	expectedTransactions := []BalanceHistory{
		{
			Timestamp: 1708855073,
			Balance:   1000,
			Method:    "WITHDRAW",
		},
		{
			Timestamp: 1708763789,
			Balance:   499,
			Method:    "DEPOSIT",
		},
	}
	path := userPath("balance-transaction")

	t.Run("Successfully buy stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserBalanceHistory", userId, "ALL", uint(0)).
			Return(expectedTransactions, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("startPage", "0")
		queryParams.Set("method", "ALL")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserBalanceHistory)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonTransactions, _ := json.Marshal(expectedTransactions)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","transactions":%v}`,
			expectedMessage,
			string(expectedJsonTransactions),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler query parameter", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserBalanceHistory", userId, "ALL", uint(0)).
			Return(expectedTransactions, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserBalanceHistory)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service get user balancehistory", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserBalanceHistory", userId, "SOME", uint(0)).
			Return(expectedTransactions, ErrOrderMethod)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("startPage", "0")
		queryParams.Set("method", "SOME")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserBalanceHistory)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrOrderMethod.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestGetUserFavoriteStock(t *testing.T) {
	returnUserStock := []string{"12345"}
	returnFavoriteStock := []StockCollectionResponse{
		{
			ID:         "12345",
			StockImage: "test",
			Name:       "test",
			Sign:       "test",
			Price:      1.1,
		},
	}
	expectedMessage := "Successfully fetched favorite stock"
	url := userPath("get-favorite")

	t.Run("Successfully get user favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserFavoriteStock", userId).
			Return(returnUserStock, nil)

		stockService.
			On("GetFavoriteStock", returnUserStock).
			Return(returnFavoriteStock, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.GET(url, userHandler.GetUserFavoriteStock)
		router.ServeHTTP(recorder, req)

		expectedJsonFavoriteStock, _ := json.Marshal(returnFavoriteStock)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"favorites":%v,"message":"%s"}`,
			string(expectedJsonFavoriteStock),
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service at user get user favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserFavoriteStock", "").
			Return(returnUserStock, ErrUser)

		stockService.
			On("GetFavoriteStock", returnUserStock).
			Return(returnFavoriteStock, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "")
		})

		router.GET(url, userHandler.GetUserFavoriteStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrUser.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service at stock get favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserFavoriteStock", userId).
			Return(returnUserStock, nil)

		stockService.
			On("GetFavoriteStock", returnUserStock).
			Return(returnFavoriteStock, ErrInvalidStock)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.GET(url, userHandler.GetUserFavoriteStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrInvalidStock.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestSignIn(t *testing.T) {
	expectedMessage := "Successfully fetched user profile"
	expectedProfile := UserResponse{
		Name:         "test",
		ProfileImage: "test",
		Email:        "test",
	}
	url := userPath("signin")

	t.Run("Successfully signin", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserIdRequest{
			UID: userId,
		}

		userService.
			On("GetUserAccount", testBody.UID).
			Return(expectedProfile, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", userId)
		})

		router.POST(url, userHandler.SignIn)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonUser, _ := json.Marshal(expectedProfile)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","user":%s}`,
			expectedMessage,
			expectedJsonUser,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler body", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		testBody := UserIdRequest{}

		userService.
			On("GetUserAccount", testBody.UID).
			Return(expectedProfile, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		jsonBody, _ := json.Marshal(testBody.UID)

		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer(jsonBody),
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.POST(url, userHandler.SignIn)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestGetUserTradingHistories(t *testing.T) {
	expectedMessage := "Successfully fetched all transactions history"
	expectedTransactions := []UserHistory{
		{
			Timestamp:   1708855336,
			StockId:     "65cc5fd45aa71b64fbb551a9",
			Price:       10,
			Amount:      1,
			Status:      "pending",
			OrderType:   "auto",
			OrderMethod: "sale",
		},
	}
	path := userPath("trade-transaction")

	t.Run("Successfully get user trading histories", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserTradingHistories", userId, uint(0)).
			Return(expectedTransactions, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("startPage", "0")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserTradingHistories)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonTransactions, _ := json.Marshal(expectedTransactions)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","transactions":%v}`,
			expectedMessage,
			string(expectedJsonTransactions),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler query parameter", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserTradingHistories", userId, uint(0)).
			Return(expectedTransactions, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserTradingHistories)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestGetUserStockHistory(t *testing.T) {
	expectedMessage := "Successfully fetched stock transaction"
	expectedTransactions := []UserHistory{
		{
			Timestamp:   1708855336,
			StockId:     "65cc5fd45aa71b64fbb551a9",
			Price:       10,
			Amount:      1,
			Status:      "pending",
			OrderType:   "auto",
			OrderMethod: "sale",
		},
	}
	path := userPath("stock-transaction")

	t.Run("Successfully test get user stock history", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserStockHistory", userId, "65cc5fd45aa71b64fbb551a9", uint(0)).
			Return(expectedTransactions, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("startPage", "0")
		queryParams.Set("stockId", "65cc5fd45aa71b64fbb551a9")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserStockHistory)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonTransactions, _ := json.Marshal(expectedTransactions)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","transactions":%v}`,
			expectedMessage,
			string(expectedJsonTransactions),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler query parameter", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserStockHistory", userId, "65cc5fd45aa71b64fbb551a9", uint(0)).
			Return(expectedTransactions, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserStockHistory)
		router.ServeHTTP(recorder, req)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on service get user stock history", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserStockHistory", userId, "", uint(0)).
			Return(expectedTransactions, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("startPage", "0")
		queryParams.Set("stockId", "")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserStockHistory)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestGetUserStockAmount(t *testing.T) {
	expectedMessage := "Successfully fetched stock ratio"
	expectedStockRatio := UserStock{
		StockId: "test",
		Amount: 1,
	}
	path := userPath("stock-ratio")

	t.Run("Successfully get user stock amount", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserStockAmount", userId, "test").
			Return(expectedStockRatio, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("stockId", "test")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserStockAmount)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonStockRatio, _ := json.Marshal(expectedStockRatio)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","stockRatio":%v}`,
			expectedMessage,
			string(expectedJsonStockRatio),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler query parameter", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("GetUserStockAmount", userId, "").
			Return(expectedStockRatio, ErrData)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"GET",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("stockId", "")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.GET(path, userHandler.GetUserStockAmount)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody :=  fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrData.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestDeleteFavoriteStock(t *testing.T) {
	expectedMessage := "Successfully deleted favorite stock"
	path := userPath("delete-favorite")

	t.Run("Successfully delete favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("DeleteFavoriteStock", userId, "test").
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"DELETE",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("stockId", "test")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.DELETE(path, userHandler.DeleteFavoriteStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler query parameter", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("DeleteFavoriteStock", userId, "").
			Return(expectedMessage, ErrInvalidStock)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"DELETE",
			path,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		queryParams := url.Values{}
		queryParams.Set("stockId", "")
		req.URL.RawQuery = queryParams.Encode()
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.DELETE(path, userHandler.DeleteFavoriteStock)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrInvalidStock.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}

func TestDeleteUserAccount(t *testing.T) {
	expectedMessage := "Successfully deleted account"
	url := userPath("delete-account")

	t.Run("Successfully delete favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("DeleteUserAccount", userId).
			Return(expectedMessage, nil)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"DELETE",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "test12345")
		})

		router.DELETE(url, userHandler.DeleteUserAccount)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			expectedMessage,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Successfully delete favorite stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()

		userService.
			On("DeleteUserAccount", "").
			Return(expectedMessage, ErrUser)

		userHandler := handler.NewUserHandler(userService, stockService)

		req, err := http.NewRequest(
			"DELETE",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Set("uid", "")
		})

		router.DELETE(url, userHandler.DeleteUserAccount)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s"}`, 
			ErrUser.Error(),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})
}
