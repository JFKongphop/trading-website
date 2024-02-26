package handler

import (
	"server/errs"
	"server/model"
	"server/service"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type stockHandler struct {
	stockService service.StockService
}

type StockCollectionRequest = model.StockCollectionRequest
type StockHistory = model.StockHistory

type CreateStockRequest struct {
	Name  string  `json:"name"`
	Sign  string  `json:"sign"`
	Price float64 `json:"price"`
}

type SetPriceRequest struct {
	Price float64 `json:"price"`
}

type CreateOrderRequest struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type EditNameRequest struct {
	Name string `json:"name"`
}

type EditSignRequest struct {
	Sign string `json:"sign"`
}

var ErrData = errs.ErrData

func NewStockHandler(stockService service.StockService) stockHandler {
	return stockHandler{stockService}
}

func (h stockHandler) CreateStockCollection(c *gin.Context) {
	file, err := c.FormFile("stock_image")

	if err != nil {
		c.JSON(200, gin.H{
			"message": "invalid file",
		})

		return
	}

	priceStr := c.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "invalid price",
		})

		return
	}

	blobFile, err := file.Open()
	if err != nil {
		c.JSON(200, gin.H{
			"message": "invalid file",
		})

		return
	}

	stock := StockCollectionRequest{
		StockImage: blobFile,
		Name:       c.PostForm("name"),
		Sign:       c.PostForm("sign"),
		Price:      price,
	}

	message, err := h.stockService.CreateStockCollection(stock)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}

// func (h stockHandler) CreateStockOrder(c *fiber.Ctx) error {
// 	stockId := c.Params("stockId")
// 	body := CreateOrderRequest{}

// 	if err := c.BodyParser(&body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": ErrData.Error(),
// 		})
// 	}

// 	order := StockHistory{
// 		ID:        c.Locals("uid").(string),
// 		Timestamp: int64(time.Now().Unix()),
// 		Price:     body.Price,
// 		Amount:    body.Amount,
// 	}

// 	message, err := h.stockService.CreateStockOrder(stockId, order)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": message,
// 	})
// }

func (h stockHandler) CreateStockOrder(c *gin.Context) {
	stockId := c.Param("stockId")
	body := CreateOrderRequest{}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	order := StockHistory{
		ID:        c.MustGet("uid").(string),
		Timestamp: int64(time.Now().Unix()),
		Price:     body.Price,
		Amount:    body.Amount,
	}

	message, err := h.stockService.CreateStockOrder(stockId, order)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}

func (h stockHandler) GetAllStockCollections(c *gin.Context) {
	stocks, err := h.stockService.GetAllStockCollections()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "successfully fetched all stocks",
		"stocks":  stocks,
	})
}

func (h stockHandler) GetTop10Stocks(c *gin.Context) {
	topStocks, err := h.stockService.GetTop10Stocks()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message":   "successfully fetched top volume stock",
		"topStocks": topStocks,
	})
}

func (h stockHandler) GetStockCollection(c *gin.Context) {
	stockId := c.Param("stockId")

	stock, err := h.stockService.GetStockCollection(stockId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "successfully fetched stock",
		"stock":   stock,
	})
}

func (h stockHandler) GetStockHistory(c *gin.Context) {
	stockId := c.Param("stockId")

	transactions, err := h.stockService.GetStockHistory(stockId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message":      "successfully fetched transactions",
		"transactions": transactions,
	})
}

func (h stockHandler) GetStockPrice(c *gin.Context) {
	stockId := c.Param("stockId")

	price, err := h.stockService.GetStockPrice(stockId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "Sucessfully fetched stock price",
		"price":   price,
	})
}

func (h stockHandler) SetStockPrice(c *gin.Context) {
	stockId := c.Param("stockId")
	body := SetPriceRequest{}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	message, err := h.stockService.SetStockPrice(stockId, body.Price)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}

func (h stockHandler) EditStockName(c *gin.Context) {
	stockId := c.Param("stockId")
	body := EditNameRequest{}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	message, err := h.stockService.EditStockName(stockId, body.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}

func (h stockHandler) EditStockSign(c *gin.Context) {
	stockId := c.Param("stockId")
	body := EditSignRequest{}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	message, err := h.stockService.EditStockSign(stockId, body.Sign)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}

func (h stockHandler) DeleteStockCollection(c *gin.Context) {
	stockId := c.Param("stockId")

	message, err := h.stockService.DeleteStockCollection(stockId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": message,
	})
}
