package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/handler"
	"server/model"
	"server/service"
	"testing"

	"github.com/gin-gonic/gin"
)

type UserAccount = model.UserAccount
type CreateAccount = model.CreateAccount

func TestSignUp(t *testing.T) {
	// expected := "Successfully created account"
	t.Run("Create user account", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		request := CreateAccount{
			Name:         "test",
			ProfileImage: "test",
			Email:        "test@gmail.com",
		}
		jsonStr, _ := json.Marshal(request)

		res := httptest.NewRecorder()
		c, r := gin.CreateTestContext(res)
		c.Request = httptest.NewRequest(
			http.MethodPost, 
			"/api/v1/stock/signup", 
			bytes.NewBuffer(jsonStr),
		)

		userService := service.NewUserServiceMock()
		stockService := service.NewStockServiceMock()
		userHandler := handler.NewUserHandler(userService, stockService)
		r.POST("/api/v1/stock/signup", userHandler.SignUp)
		r.ServeHTTP(res, c.Request)

		fmt.Println(res.Code)





		// userAccount := CreateAccount{
		// 	Name:         "test",
		// 	ProfileImage: "test",
		// 	Email:        "test@gmail.com",
		// }
		// userService.On("CreateUserAccount", userAccount).Return(expected, nil)

		// 

		// app := gin.Default()
		// app.POST("/api/v1/stock/signup", userHandler.SignUp)

		// req := httptest.NewRequest("POST", "/api/v1/stock/signup", nil)
		// res, _ := app.Test(req)
		// defer res.Body.Close()
	})
}