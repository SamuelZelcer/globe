package userHandler

import userService "globe/internal/service/user"


type Handler struct {
	userService userService.Service
}

func InitHandler(userService userService.Service) *Handler {
	return &Handler{userService: userService}
}