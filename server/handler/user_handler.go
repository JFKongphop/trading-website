package handler

import (
	"server/errs"
	"server/model"
	"server/service"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userService service.UserService
	stockService service.StockService
}

type CreateAccount = model.CreateAccount

type UserIdRequest struct {
	UID string `json:"uid"`
}

type UserBalanceRequest struct {
	Balance float64 `json:"balance"`
}

var ErrUser = errs.ErrUser

func NewUserHandler(userService service.UserService, stockService service.StockService) userHandler {
	return userHandler{userService, stockService}
}

func (h userHandler) SignUp(c *fiber.Ctx) error {
	body := CreateAccount{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrUser.Error(),
		})
	}

	message, err := h.userService.CreateUserAccount(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h userHandler) DepositBalance(c *fiber.Ctx) error {
	body := UserBalanceRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrUser.Error(),
		})
	}

	deposit := body.Balance
	uid := c.Locals("uid").(string)
	message, err := h.userService.DepositBalance(uid, deposit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h userHandler) WithdrawBalance(c *fiber.Ctx) error {
	body := UserBalanceRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrUser.Error(),
		})
	}

	withdraw := body.Balance
	uid := c.Locals("uid").(string)
	message, err := h.userService.WithdrawBalance(uid, withdraw)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
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

func (h userHandler) SignIn(c *fiber.Ctx) error {
	body := UserIdRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrUser.Error(),
		})
	}

	user, err := h.userService.GetUserAccount(body.UID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		}) 
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched user profile",
		"user": user,
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
