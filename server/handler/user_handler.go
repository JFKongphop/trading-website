package handler

import (
	"server/errs"
	"server/model"
	"server/service"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h userHandler) SignUp(c *gin.Context) {
	body := CreateAccount{}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	message, err := h.userService.CreateUserAccount(body)
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

func (h userHandler) DepositBalance(c *gin.Context) {
	body := UserBalanceRequest{}


	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	uid := c.MustGet("uid").(string)
	deposit := body.Balance
	
	message, err := h.userService.DepositBalance(uid, deposit)
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

func (h userHandler) WithdrawBalance(c *gin.Context) {
	body := UserBalanceRequest{}
	
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	uid := c.MustGet("uid").(string)
	withdraw := body.Balance
	message, err := h.userService.WithdrawBalance(uid, withdraw)
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

func (h userHandler) BuyStock(c *gin.Context) {
	body := OrderRequest{}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	body.UserId = c.MustGet("uid").(string)

	message, err := h.userService.BuyStock(body)
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

func (h userHandler) SaleStock(c *gin.Context) {
	body := OrderRequest{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	body.UserId = c.MustGet("uid").(string)

	message, err := h.userService.SaleStock(body)
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

func (h userHandler) SetFavoriteStock(c *gin.Context) {
	body := UserSetFavoriteRequest{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	uid := c.MustGet("uid").(string)
	stockId := body.StockId

	message, err := h.userService.SetFavoriteStock(uid, stockId)
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

func (h userHandler) GetUserBalance(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	balance, err := h.userService.GetUserBalance(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return 
	}
	
	c.JSON(200, gin.H{
		"message": "Successfully fetched user balance",
		"balance": balance,
	})
}

func (h userHandler) GetUserBalanceHistory(c *gin.Context) {
	startPage, err := strconv.Atoi(c.Query("startPage"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}
	uid := c.MustGet("uid").(string)
	method := c.Query("method")

	transaction, err := h.userService.GetUserBalanceHistory(uid, method, uint(startPage))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "successfully fetched transaction balance",
		"transactions": transaction,
	})
}

func (h userHandler) GetUserFavoriteStock(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	stockId, err := h.userService.GetUserFavoriteStock(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	favoriteStock, err := h.stockService.GetFavoriteStock(stockId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}


	c.JSON(200, gin.H{
		"message": "successfully fetched favorite stock",
		"favorites": favoriteStock,
	})
}

func (h userHandler) SignIn(c *gin.Context) {
	body := UserIdRequest{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	user, err := h.userService.GetUserAccount(body.UID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "successfully fetched user profile",
		"user": user,
	})
}

func (h userHandler) GetUserTradingHistories(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	startPage, err := strconv.Atoi(c.Query("startPage"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}

	transactions, err := h.userService.GetUserTradingHistories(uid, uint(startPage))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "successfully fetched all transaction history",
		"transactions": transactions,
	})
}

func (h userHandler) GetUserStockHistory(c *gin.Context) {
	startPage, err := strconv.Atoi(c.Query("startPage"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrData.Error(),
		})

		return
	}
	uid := c.MustGet("uid").(string)
	stockId := c.Query("stockId")
	
	transactions, err := h.userService.GetUserStockHistory(uid, stockId, uint(startPage))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "successfully fetched stock transaction",
		"transactions": transactions,
	})
}

func (h userHandler) GetUserStockAmount(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	stockId := c.Query("stockId")

	stockRatio, err := h.userService.GetUserStockAmount(uid, stockId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "successfully fetched stock ratio",
		"stockRatio": stockRatio,
	})
}

func (h userHandler) DeleteFavoriteStock(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	stockId := c.Query("stockId")

	message, err := h.userService.DeleteFavoriteStock(uid, stockId)
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

func (h userHandler) DeleteUserAccount(c *gin.Context) {
	// uid := c.Locals("uid").(string)

	message, err := h.userService.DeleteUserAccount("test1")
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
