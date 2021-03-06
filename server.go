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
	"net/http"
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
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
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
		authRoutes.GET("/login", func(c *gin.Context) {
			c.HTML(
				http.StatusOK,
				"login.html",
				gin.H{
					"title": "Login",
				},
			)
		})
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		userRoutes.POST("/subscribe", subController.Insert)
		userRoutes.DELETE("/unsubscribe/:id", subController.Delete)
		userRoutes.GET("/main", func(c *gin.Context) {
			c.HTML(
				http.StatusOK,
				"index.html",
				gin.H{
					"title": "Home Page",
				},
			)
		})
	}

	r.GET("/main", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "Home Page",
			},
		)

	})

	go Migrations()

	r.Run(":8080")
}

func bitmexInstrument() {
	//verb := "GET"
	//path := "/api/v1/instrument"
	//expires := fmt.Sprint(time.Now().Local().Add(time.Minute * time.Duration(10)).Unix())
	//var secret = "mvK7p-zYF5He2eistXxXUvASoJWRGvp6eOO5TF2gn4BHI2iB"

	//signature := hmac.New(sha256.New, []byte(secret))
	//data := verb + path + expires
	//signature.Write([]byte(data))
	//sha := hex.EncodeToString(signature.Sum(nil))
	//fmt.Println(sha)
}

//func WsHandler() {
//	r := gin.Default()
//	r.LoadHTMLGlob("templates/*")
//}
