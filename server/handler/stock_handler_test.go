package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"server/errs"
	"server/handler"
	"server/model"
	"server/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type StockCollectionRequest = model.StockCollectionRequest
type CreateOrderRequest = model.CreateOrderRequest
type StockHistory = model.StockHistory
type TopStock = model.TopStock
type StockHistoryResponse = model.StockHistoryResponse
type Graph = model.Graph
type SetPriceRequest = model.SetPriceRequest
type EditNameRequest = model.EditNameRequest
type EditSignRequest = model.EditSignRequest

var (
	ErrPrice = errs.ErrPrice
)

func stockPath(route string) string {
	return fmt.Sprintf("/api/v1/stock/%s", route)
}

func TestCreateStockCollection(t *testing.T) {
	expectedMessage := "Successfully created stock collection"
	url := stockPath("create-stock")

	t.Run("Successfully create stock collection", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		stockService := service.NewStockServiceMock()

		stockService.
			On("CreateStockCollection", mock.Anything).
			Return(expectedMessage, nil)

		stockHandler := handler.NewStockHandler(stockService)

		file, err := os.Open("test.png")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.CreateFormFile("stock_image", filepath.Base(file.Name()))
		writer.WriteField("name", "test")
		writer.WriteField("sign", "test")
		writer.WriteField("price", "1.5")
		writer.Close()

		req, err := http.NewRequest(
			"POST",
			url,
			body,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		recorder := httptest.NewRecorder()

		router.POST(url, stockHandler.CreateStockCollection)
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

	t.Run("Error on handler file", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		stockService := service.NewStockServiceMock()

		stockService.
			On("CreateStockCollection", mock.Anything).
			Return(expectedMessage, nil)

		stockHandler := handler.NewStockHandler(stockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		req, err := http.NewRequest(
			"POST",
			url,
			body,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		recorder := httptest.NewRecorder()

		router.POST(url, stockHandler.CreateStockCollection)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedResponseBody := `{"message":"invalid file"}`

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler post form price", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		stockService := service.NewStockServiceMock()

		stockService.
			On("CreateStockCollection", mock.Anything).
			Return(expectedMessage, ErrPrice)

		stockHandler := handler.NewStockHandler(stockService)

		file, err := os.Open("test.png")
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.CreateFormFile("stock_image", filepath.Base(file.Name()))
		writer.WriteField("name", "test")
		writer.WriteField("sign", "test")
		writer.WriteField("price", "test")
		writer.Close()

		req, err := http.NewRequest(
			"POST",
			url,
			body,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		recorder := httptest.NewRecorder()

		router.POST(url, stockHandler.CreateStockCollection)
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
			ErrPrice.Error(),
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

func TestCreateStockOrder(t *testing.T) {
	expectedMessage := "Successfully created stock order"

	t.Run("Successfully create stock order", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("create-order/12345")

		stockService := service.NewStockServiceMock()

		testBody := CreateOrderRequest{
			Amount: 1,
			Price:  1,
		}

		order := StockHistory{
			ID:        userId,
			Timestamp: int64(time.Now().Unix()),
			Price:     1,
			Amount:    1,
		}

		stockService.
			On("CreateStockOrder", "12345", order).
			Return(expectedMessage, nil)

		stockHandler := handler.NewStockHandler(stockService)

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
			c.Params = gin.Params{
				{
					Key:   "stockId",
					Value: "12345",
				},
			}
		})

		router.POST(url, stockHandler.CreateStockOrder)
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

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("create-order/")

		stockService := service.NewStockServiceMock()

		testBody := CreateOrderRequest{
			Amount: 1,
			Price:  1,
		}

		order := StockHistory{
			ID:        userId,
			Timestamp: int64(time.Now().Unix()),
			Price:     1,
			Amount:    1,
		}

		stockService.
			On("CreateStockOrder", "", order).
			Return(expectedMessage, ErrData)

		stockHandler := handler.NewStockHandler(stockService)

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

		router.POST(url, stockHandler.CreateStockOrder)
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

func TestGetAllStockCollections(t *testing.T) {
	expectedMessage := "Successfully fetched all stocks"
	expectedStockCollections := []StockCollectionResponse{
		{
			ID:         "1",
			StockImage: "test",
			Name:       "test",
			Sign:       "test",
			Price:      1,
		},
	}
	url := stockPath("collections")

	t.Run("Successfully get all stock collections", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetAllStockCollections").
			Return(expectedStockCollections, nil)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		router.GET(url, stockHandler.GetAllStockCollections)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonStockCollections, _ := json.Marshal(expectedStockCollections)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","stocks":%v}`,
			expectedMessage,
			string(expectedJsonStockCollections),
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

func TestGetTop10Stocks(t *testing.T) {
	expectedMessage := "Successfully fetched top volume stock"
	expectedTopStocks := []TopStock{
		{
			ID:    "1",
			Sign:  "test",
			Price: 1,
		},
	}
	url := stockPath("top-stocks")

	t.Run("Successfully get top 10 stocks", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetTop10Stocks").
			Return(expectedTopStocks, nil)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		router.GET(url, stockHandler.GetTop10Stocks)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonStockCollections, _ := json.Marshal(expectedTopStocks)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","topStocks":%v}`,
			expectedMessage,
			string(expectedJsonStockCollections),
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

func TestGetStockCollection(t *testing.T) {
	expectedMessage := "Successfully fetched stock"
	expectedStock := StockCollectionResponse{
		ID:         "1",
		StockImage: "test",
		Name:       "test",
		Sign:       "test",
		Price:      1,
	}
	

	t.Run("Successfully get stock collection", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("collection/12345")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockCollection", "12345").
			Return(expectedStock, nil)

		stockHandler := handler.NewStockHandler(stockService)

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
			c.Params = gin.Params{
				{
					Key:   "stockId",
					Value: "12345",
				},
			}
		})

		router.GET(url, stockHandler.GetStockCollection)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedJsonStockCollections, _ := json.Marshal(expectedStock)

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","stock":%v}`,
			expectedMessage,
			string(expectedJsonStockCollections),
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("create-order/")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockCollection", "").
			Return(expectedStock, ErrData)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		router.GET(url, stockHandler.GetStockCollection)
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

func TestGetStockHistory(t *testing.T) {
	expectedMessage := "Successfully fetched transactions"
	expectedTransactions := []StockHistoryResponse{
		{
			Amount: 1,
			Price: 1,
		},
	}

	t.Run("Successfully get stock history", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("transaction/12345")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockHistory", "12345").
			Return(expectedTransactions, nil)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "12345",
				},
			}
		})

		router.GET(url, stockHandler.GetStockHistory)
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

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("transaction/")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockHistory", "").
			Return(expectedTransactions, ErrInvalidStock)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		router.GET(url, stockHandler.GetStockHistory)
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

func TestGetStockPrice(t *testing.T) {
	expectedMessage := "Sucessfully fetched stock price"
	expectedPrice := 1.1

	t.Run("Successfully get stock price", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("price/12345")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockPrice", "12345").
			Return(expectedPrice, nil)

		stockHandler :=  handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "12345",
				},
			}
		})

		router.GET(url, stockHandler.GetStockPrice)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"message":"%s","price":%.1f}`,
			expectedMessage,
			expectedPrice,
		)

		if recorder.Body.String() != expectedResponseBody {
			t.Errorf(
				"Expected response body %s, got %s",
				expectedResponseBody,
				recorder.Body.String(),
			)
		}
	})

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("price/")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockPrice", "").
			Return(expectedPrice, ErrInvalidStock)

		stockHandler :=  handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		router.GET(url, stockHandler.GetStockPrice)
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

func TestGetStockGraph(t *testing.T) {
	expectedMessage := "Successfully fetched stock graph"
	expectedGraph := []Graph{
		{
			X: 1, 
			Y: []float64{1.1, 1.2, 1.3, 1.4},
		},
	}

	t.Run("Successfully get stock graph", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("graph/12345")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockGraph", "12345").
			Return(expectedGraph, nil)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "12345",
				},
			}
		})

		router.GET(url, stockHandler.GetStockGraph)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf(
				"Expected status code %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}

		expectedJsonGraph, err := json.Marshal(expectedGraph)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		expectedResponseBody := fmt.Sprintf(
			`{"graph":%v,"message":"%s"}`,
			string(expectedJsonGraph),
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

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("graph/")

		stockService := service.NewStockServiceMock()

		stockService.
			On("GetStockGraph", "").
			Return(expectedGraph, ErrInvalidStock)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"GET",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()

		router.GET(url, stockHandler.GetStockGraph)
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

func TestSetStockPrice(t *testing.T) {
	expectedMessage := "Successfully set price"

	t.Run("Successfully set stock price", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("set-price/12345")

		testBody := SetPriceRequest{
			Price: 1.1,
		}

		stockService := service.NewStockServiceMock()

		stockService.
			On("SetStockPrice", "12345", testBody.Price).
			Return(expectedMessage, nil)

		stockHandler := handler.NewStockHandler(stockService)

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
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "12345",
				},
			}
		})

		router.POST(url, stockHandler.SetStockPrice)
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

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("set-price/")

		testBody := SetPriceRequest{
			Price: 1.1,
		}

		stockService := service.NewStockServiceMock()

		stockService.
			On("SetStockPrice", "", testBody.Price).
			Return(expectedMessage, ErrInvalidStock)

		stockHandler := handler.NewStockHandler(stockService)

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

		router.POST(url, stockHandler.SetStockPrice)
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
			ErrInvalidStock,
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

func TestEditStockName(t *testing.T) {
	expectedMessage := "Successfully updated name"

	t.Run("Successfully edit stock name", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("edit-name/12345")

		stockService := service.NewStockServiceMock()

		testBody := EditNameRequest{
			Name: "test",
		}

		stockService.
			On("EditStockName", "12345", testBody.Name).
			Return(expectedMessage, nil)

		stockHandler := handler.NewStockHandler(stockService)

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
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "12345",
				},
			}
		})

		router.POST(url, stockHandler.EditStockName)
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

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("edit-name/")

		stockService := service.NewStockServiceMock()

		testBody := EditNameRequest{
			Name: "test",
		}

		stockService.
			On("EditStockName", "", testBody.Name).
			Return(expectedMessage, ErrInvalidStock)

		stockHandler := handler.NewStockHandler(stockService)

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

		router.POST(url, stockHandler.EditStockName)
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

func TestEditStockSign(t *testing.T) {
	expectedMessage := "Successfully updated sign"

	t.Run("Successfully edit stock sign", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("edit-sign/12345")

		stockService := service.NewStockServiceMock()

		testBody := EditSignRequest{
			Sign: "test",
		}

		stockService.
			On("EditStockSign", "12345", testBody.Sign).
			Return(expectedMessage, nil)

		stockHandler := handler.NewStockHandler(stockService)

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
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "12345",
				},
			}
		})

		router.POST(url, stockHandler.EditStockSign)
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

	t.Run("Error on handler param", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("edit-sign/")

		stockService := service.NewStockServiceMock()

		testBody := EditSignRequest{
			Sign: "test",
		}

		stockService.
			On("EditStockSign", "", testBody.Sign).
			Return(expectedMessage, ErrInvalidStock)

		stockHandler := handler.NewStockHandler(stockService)

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

		router.POST(url, stockHandler.EditStockSign)
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

func TestDeleteStockCollection(t *testing.T) {
	expectedMessage := "Successfully deleted stock"

	t.Run("Successfully deleted stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("delete/12345")

		stockService := service.NewStockServiceMock()

		stockService.
			On("DeleteStockCollection", "12345").
			Return(expectedMessage, nil)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"DELETE",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "12345",
				},
			}
		})

		router.DELETE(url, stockHandler.DeleteStockCollection)
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

	t.Run("Successfully deleted stock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		url := stockPath("delete/")

		stockService := service.NewStockServiceMock()

		stockService.
			On("DeleteStockCollection", "").
			Return(expectedMessage, ErrInvalidStock)

		stockHandler := handler.NewStockHandler(stockService)

		req, err := http.NewRequest(
			"DELETE",
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		router.Use(func(c *gin.Context) {
			c.Params = []gin.Param{
				{
					Key: "stockId",
					Value: "",
				},
			}
		})

		router.DELETE(url, stockHandler.DeleteStockCollection)
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
 
