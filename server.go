package main

import (
	"BitmaxGinGorilla/config"
	"BitmaxGinGorilla/controller"
	"BitmaxGinGorilla/middleware"
	"BitmaxGinGorilla/migrations"
	"BitmaxGinGorilla/repository"
	"BitmaxGinGorilla/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

var (
	db             *gorm.DB                       = config.SetupDatabaseConnection()
	userRepository repository.UserRepository      = repository.NewUserRepository(db)
	subRepository  repository.SubscribeRepository = repository.NewSubRepository(db)
	jwtService     service.JWTService             = service.NewJWTService()
	userService    service.UserService            = service.NewUserService(userRepository)
	subService     service.SubscribeService       = service.NewSubService(subRepository)
	authService    service.AuthService            = service.NewAuthService(userRepository)
	authController controller.AuthController      = controller.NewAuthController(authService, jwtService)
	userController controller.UserController      = controller.NewUserController(userService, jwtService)
	subController  controller.SubscribeController = controller.NewSubController(subService, jwtService)
	Migrations                                    = migrations.DbMigrate
)

func main() {
	defer config.CloseDatabaseConnection(db)
	//id := db.Raw("SELECT id FROM subscribes WHERE id= ?", 1)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		userRoutes.POST("/subscribe", subController.Insert)
		userRoutes.DELETE("/unsubscribe/:id", subController.Delete)
	}

	go Migrations()

	r.Run(":8080")
}
