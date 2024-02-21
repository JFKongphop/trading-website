package handler

import (
	"server/model"
	"server/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type stockHandler struct {
	stockService service.StockService
}

type StockCollectionRequest = model.StockCollectionRequest

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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) GetAllStockCollections(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) GetTop10Stocks(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) GetStockCollection(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) GetStockHistory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) SetStockPrice(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) EditStockName(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) EditStockSign(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h stockHandler) DeleteStockCollection(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}