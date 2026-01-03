package main

import (
	userHandler "globe/internal/handler/user"
	"globe/internal/repository"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
	userService "globe/internal/service/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    DB := repository.InitDB()

    transactions := transactions.Init(DB)
    
    userRepository := user.InitRepository(DB)
    unverifiedUserRepository := unverifiedUser.InitRepository(DB)

    email := email.InitEmail()
    jwtManager := JWT.Init()

    userService := userService.Init(
        userRepository,
        unverifiedUserRepository,
        email,
        jwtManager,
        transactions,
    )

    userHandler := userHandler.Init(userService)
    
    e := echo.New()
    e.Use(middleware.RequestLogger())
    e.Use(middleware.CORS())
    
    e.POST("/user/sign-up", userHandler.SignUp)
    e.POST("/user/verification", userHandler.Verification)
    e.POST("/user/verification/get-new-code", userHandler.GetNewCode)

    e.Start("127.0.0.1:8080")
}
