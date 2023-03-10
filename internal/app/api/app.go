package api

import (
	"fmt"
	"log"
	"os"
	"time"

	"signupin-api/internal/pkg/user"

	userrepo "signupin-api/internal/pkg/user/persistence"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"

	kkva "github.com/kkodecaffeine/go-common/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App interface {
	Init()
	RegisterRoute(driver *gin.Engine)
	Clean() error
}

type apiApp struct {
	client *mongo.Client
}

func (app *apiApp) Init() {
	err := godotenv.Load("../config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_ = mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: 12 * time.Second}, "kkodecaffeine", options.Client().ApplyURI(os.Getenv("MONGO_URL")))
}

func (app *apiApp) RegisterRoute(driver *gin.Engine) {
	v := validator.New()

	user_uc := user.NewUsecase(userrepo.New(app.client))
	NewController(driver, v, user_uc)
}

func (app *apiApp) Clean() error {
	return nil
}

// CreateAPIApp returns new core.App implementation
func CreateAPIApp() {
	router := gin.Default()
	router.RouterGroup = *router.Group("/api")

	app := &apiApp{}
	app.Init()

	router.Use(gin.Recovery())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("customEmail", kkva.RegexEmail())
		v.RegisterValidation("customPhone", kkva.RegexPhone())
	}

	frontserver := os.Getenv("FRONT_SERVER_HOST")

	// Enable CORS policy
	router.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{frontserver},
			AllowMethods:     []string{"GET, POST"},
			AllowHeaders:     []string{"Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))

	app.RegisterRoute(router)

	router.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))

	app.Clean()
}
