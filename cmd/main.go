package main

import (
	productHandler "globe/internal/handler/product"
	userHandler "globe/internal/handler/user"
	"globe/internal/repository"
	serviceDTOs "globe/internal/repository/dtos/service"
	"globe/internal/repository/entities/product"
	"globe/internal/repository/entities/refreshToken"
	"globe/internal/repository/entities/unverifiedUser"
	"globe/internal/repository/entities/user"
	"globe/internal/repository/redis"
	"globe/internal/repository/transactions"
	"globe/internal/service/email"
	JWT "globe/internal/service/jwt"
	productService "globe/internal/service/product"
	refreshTokenService "globe/internal/service/refreshToken"
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
    productRepository := product.InitRepository(DB)

    // service
    email := email.InitEmail()
    jwtManager := JWT.Init()

    refreshTokenService := refreshTokenService.Init(
        refreshTokenRepository,
        userRepository,
        jwtManager,
        redisRepository,
    )
    userService := userService.Init(
        &serviceDTOs.UserDependencies{
            UserRepository: userRepository,
            UnverifiedUserRepository: unverifiedUserRepository,
            Email: email,
            JWTManager: jwtManager,
            Transactions: transactions,
            Redis: redisRepository,
            RefreshTokenService: refreshTokenService,
        },
    )
    productService := productService.Init(
        &serviceDTOs.ProductDependencies{
            ProductRepository: productRepository,
            UserRepository: userRepository,
            Email: email,
            Redis: redisRepository,
            JWTManager: jwtManager,
            RefreshTokenService: refreshTokenService,
        },
    )

    // handler
    userHandler := userHandler.Init(userService)
    productHandler := productHandler.Init(productService)

    // initialize HTTP server
    e := echo.New()
    e.Use(middleware.RequestLogger())
    e.Use(middleware.CORS())
    
    // authorization & authentication
    e.POST("/user/sign-up", userHandler.SignUp)
    e.POST("/user/sign-in", userHandler.SignIn)
    
    // verification
    e.POST("/user/verification", userHandler.Verification)
    e.POST("/user/verification/get-new-code", userHandler.GetNewCode)
    e.POST("/user/verification/send-code-again", userHandler.SendCodeAgain)

    // update
    e.POST("/user/update/username", userHandler.UpdateUsername)
    e.POST("/user/update/email", userHandler.UpdateEmail)
    e.POST("/user/update/email/verification", userHandler.NewEmailVerification)
    e.POST("/user/update/password", userHandler.UpdatePassword)
    e.POST("/user/update/password/verification", userHandler.NewPasswordVerification)

    // product
    e.POST("/product/create", productHandler.Create)
    e.POST("/product/update", productHandler.Update)
    e.POST("/product/delete", productHandler.Delete)
    e.POST("/product/search", productHandler.Search)

    e.Start("127.0.0.1:8080")
}
