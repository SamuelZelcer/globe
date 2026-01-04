package main

import (
	userHandler "globe/internal/handler/user"
	"globe/internal/repository"
	"globe/internal/repository/entities/refreshToken"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
	userService "globe/internal/service/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    DB := repository.InitDB()
    REDIS := redis.InitRedis()

    // repository
    redisRepository := redis.InitRepository(REDIS)

    transactions := transactions.Init(DB)
    
    userRepository := user.InitRepository(DB)
    unverifiedUserRepository := unverifiedUser.InitRepository(DB)
    refreshTokenRepository := refreshToken.InitRepository(DB)

    // service
    email := email.InitEmail()
    jwtManager := JWT.Init()

    userService := userService.Init(
        userRepository,
        unverifiedUserRepository,
        email,
        jwtManager,
        transactions,
        redisRepository,
        refreshTokenRepository,
    )

    // handler
    userHandler := userHandler.Init(userService)
    
    e := echo.New()
    e.Use(middleware.RequestLogger())
    e.Use(middleware.CORS())
    
    // API
    e.POST("/user/sign-up", userHandler.SignUp)
    e.POST("/user/verification", userHandler.Verification)
    e.POST("/user/verification/get-new-code", userHandler.GetNewCode)
    e.POST("/user/verification/send-code-again", userHandler.SendCodeAgain)
    e.POST("/user/sign-in", userHandler.SignIn)

    e.Start("127.0.0.1:8080")
}
