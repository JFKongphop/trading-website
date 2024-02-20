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
type OrderRequest = model.OrderRequest

type UserIdRequest struct {
	UID string `json:"uid"`
}

type UserBalanceRequest struct {
	Balance float64 `json:"balance"`
}

type UserSetFavoriteRequest struct {
	StockId string `json:"stockId"`
}

type FilterBalanceRequest struct {
	Method string `json:"method"`
	StartPage uint `json:"startPage"`
}

var ErrUser = errs.ErrUser
var ErrMoney = errs.ErrMoney

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
	body := OrderRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrUser.Error(),
		})
	}

	body.UserId = c.Locals("uid").(string)
	message, err := h.userService.BuyStock(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h userHandler) SaleStock(c *fiber.Ctx) error {
	body := OrderRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrUser.Error(),
		})
	}

	body.UserId = c.Locals("uid").(string)
	message, err := h.userService.SaleStock(body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h userHandler) SetFavoriteStock(c *fiber.Ctx) error {
	body := UserSetFavoriteRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": ErrUser.Error(),
		})
	}

	stockId := body.StockId
	uid := c.Locals("uid").(string)
	message, err := h.userService.SetFavoriteStock(uid, stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func (h userHandler) GetUserBalance(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	balance, err := h.userService.GetUserBalance(uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched balance",
		"balance": balance,
	})
}

func (h userHandler) GetUserBalanceHistory(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	startPage := c.QueryInt("startPage")
	method := c.Query("method")

	transaction, err := h.userService.GetUserBalanceHistory(uid, method, uint(startPage))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		}) 
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched transaction balance",
		"transactions": transaction,
	})
}

func (h userHandler) GetUserFavoriteStock(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	stockId, err := h.userService.GetUserFavoriteStock(uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	favoriteStock, err := h.stockService.GetFavoriteStock(stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched favorite stock",
		"favorites": favoriteStock,
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
	uid := c.Locals("uid").(string)
	startPage := c.QueryInt("startPage")
	transactions, err := h.userService.GetUserTradingHistories(uid, uint(startPage))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched all transaction history",
		"transactions": transactions,
	})
}

func (h userHandler) GetUserStockHistory(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	startPage := c.QueryInt("startPage")
	stockId := c.Query("stockId")
	transactions, err := h.userService.GetUserStockHistory(uid, stockId, uint(startPage))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched stock transaction",
		"transactions": transactions,
	})
}

func (h userHandler) GetUserStockAmount(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	stockId := c.Query("stockId")

	stockRatio, err := h.userService.GetUserStockAmount(uid, stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully fetched stock ratio",
		"stockRatio": stockRatio,
	})
}

func (h userHandler) DeleteFavoriteStock(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete favorite stock",
	})
}

func (h userHandler) DeleteUserAccount(c *fiber.Ctx) error {
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully delete account",
	})
}
