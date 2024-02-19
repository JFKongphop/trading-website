package handler

import (
	"server/service"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) userHandler {
	return userHandler{userService}
}

// func (h userHandler) Signin(c *fiber.Ctx) {
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{

// 	})
// }

