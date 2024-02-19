package handler

import "server/service"

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) userHandler {
	return userHandler{userService}
}

func (h userHandler)

