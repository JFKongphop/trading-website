package handler

import (
	"server/service"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userService service.UserService
	stockService service.StockService
}

func NewUserHandler(userService service.UserService, stockService service.StockService) userHandler {
	return userHandler{userService, stockService}
}

func (h userHandler) SignUp(c *fiber.Ctx) error {
	
	h.userService.CreateUserAccount()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) DepositBalance(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) WithdrawBalance(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) BuyStock(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) SaleStock(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) SetFavoriteStock(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) GetUserBalance(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) GetUserBalanceHistory(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) GetUserFavoriteStock(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

// get user account
func (h userHandler) SignIn(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) GetUserTradingHistories(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) GetUserStockHistory(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) GetUserStockAmount(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) DeleteFavoriteStock(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}

func (h userHandler) DeleteUserAccount(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}
