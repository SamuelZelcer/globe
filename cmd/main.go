package main

import (
	userHandler "globe/internal/handler/user"
	"globe/internal/repository"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/transaction"
	"globe/internal/service/jwt"
	userService "globe/internal/service/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	DB := repository.InitDB()

	userRepository := user.InitRepository(DB)
	unverifiedUserRepository := unverifiedUser.InitRepository(DB)

	transactions := transaction.InitTransactions(DB)

	jwtManager := jwt.InitJWTManager()

	userService := userService.InitService(
		userRepository,
		unverifiedUserRepository,
		transactions,
		jwtManager,
	)

	userHandler := userHandler.InitHandler(userService)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORS())
	
	e.POST("/user/create", userHandler.Create)
	e.POST("/user/verify", userHandler.Verify)

	e.Start("127.0.0.1:8080")
}