package main

import (
	"BitmaxGinGorilla/config"
	"BitmaxGinGorilla/controller"
	"BitmaxGinGorilla/middleware"
	"BitmaxGinGorilla/migrations"
	"BitmaxGinGorilla/repository"
	"BitmaxGinGorilla/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
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
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

func main() {
	defer config.CloseDatabaseConnection(db)
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
	r.LoadHTMLGlob("template/*")
	{
		r.GET("/template", func(c *gin.Context) {
			c.HTML(
				http.StatusOK,
				"index.html",
				gin.H{
					"title": "Home Page",
				},
			)

		})
	}

	go Migrations()

	r.Run(":8080")
}

func bitmexInstrument() {
	//verb := "GET"
	//path := "/api/v1/instrument"
	//expires := fmt.Sprint(time.Now().Local().Add(time.Minute * time.Duration(10)).Unix())
	//var secret = "mvK7p-zYF5He2eistXxXUvASoJWRGvp6eOO5TF2gn4BHI2iB"
	req, err := http.NewRequest(
		http.MethodGet,
		"https://testnet.bitmex.com/api/v1/instrument",
		nil,
	)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	//req.Header.Add("Accept", "application/json")
	//req.Header.Add("apiKey", os.Getenv("API_KEY"))
	//req.Header.Add("apiSecret", os.Getenv("API_SECRET"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error sending HTTP request: %v", err)
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading HTTP response body: %v", err)
	}

	log.Println("We got the response:", string(responseBytes))

	//signature := hmac.New(sha256.New, []byte(secret))
	//data := verb + path + expires
	//signature.Write([]byte(data))
	//sha := hex.EncodeToString(signature.Sum(nil))
	//fmt.Println(sha)

}
