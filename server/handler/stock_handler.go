package handler

import (
	"server/errs"
	"server/model"
	"server/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
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
	Price  float64 `json:"price"`
}

type CreateOrderRequest struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type EditNameRequest struct {
	Name  string  `json:"name"`
}

type EditSignRequest struct {
	Sign  string  `json:"sign"`
}


var ErrData = errs.ErrData

func NewStockHandler(stockService service.StockService) stockHandler {
	return stockHandler{stockService}
}

func (h stockHandler) CreateStockCollection(c *fiber.Ctx) error {
	file, err := c.FormFile("stock_image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error file",
		})
	}

	priceStr := c.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid price format",
		})
	}

	blobFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error bolb",
		})
	}

	stock := StockCollectionRequest{
		StockImage: blobFile,
		Name:       c.FormValue("name"),
		Sign:       c.FormValue("sign"),
		Price:      price,
	}

	message, err := h.stockService.CreateStockCollection(stock)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h stockHandler) CreateStockOrder(c *fiber.Ctx) error {
	stockId := c.Params("stockId")
	body := CreateOrderRequest{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrData.Error(),
		})
	}

	order := StockHistory{
		ID:        c.Locals("uid").(string),
		Timestamp: int64(time.Now().Unix()),
		Price:     body.Price,
		Amount:    body.Amount,
	}

	message, err := h.stockService.CreateStockOrder(stockId, order)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h stockHandler) GetAllStockCollections(c *fiber.Ctx) error {
	stocks, err := h.stockService.GetAllStockCollections()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched all stocks",
		"stocks": stocks,
	})
}

func (h stockHandler) GetTop10Stocks(c *fiber.Ctx) error {
	topStocks, err := h.stockService.GetTop10Stocks()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched top volume stock",
		"topStocks": topStocks,
	})
}

func (h stockHandler) GetStockCollection(c *fiber.Ctx) error {
	stockId := c.Params("stockId")

	stock, err := h.stockService.GetStockCollection(stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched stock",
		"stock": stock,
	})
}

func (h stockHandler) GetStockHistory(c *fiber.Ctx) error {
	stockId := c.Params("stockId")

	transactions, err := h.stockService.GetStockHistory(stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched transactions",
		"transactions": transactions,
	})
}

// func (h stockHandler) SetStockPrice(c *fiber.Ctx) error {
	// stockId := c.Params("stockId")
	// body := SetPriceRequest{}

// 	if err := c.BodyParser(&body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": ErrData.Error(),
// 		})
// 	}

// 	message, err := h.stockService.SetStockPrice(stockId, body.Price)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": message,
// 	})
// }

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

func (h stockHandler) EditStockName(c *fiber.Ctx) error {
	stockId := c.Params("stockId")
	body := EditNameRequest{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrData.Error(),
		})
	}

	message, err := h.stockService.EditStockName(stockId, body.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h stockHandler) EditStockSign(c *fiber.Ctx) error {
	stockId := c.Params("stockId")
	body := EditSignRequest{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrData.Error(),
		})
	}

	message, err := h.stockService.EditStockSign(stockId, body.Sign)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h stockHandler) DeleteStockCollection(c *fiber.Ctx) error {
	stockId := c.Params("stockId")

	message, err := h.stockService.DeleteStockCollection(stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}
