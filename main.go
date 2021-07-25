package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vcycyv/blog/handler"
	infra "github.com/vcycyv/blog/infrastructure"
	"github.com/vcycyv/blog/infrastructure/repository"
	"github.com/vcycyv/blog/middleware"
	"github.com/vcycyv/blog/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	router := initRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", infra.AppSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    time.Duration(60) * time.Second,
		WriteTimeout:   time.Duration(60) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func initRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.JSONAppErrorReporter())

	logrus.SetLevel(logrus.DebugLevel)
	r.Use(middleware.LoggerToFile(logrus.StandardLogger()))
	r.Use(middleware.CORS())
	gin.SetMode(infra.AppSetting.RunMode)

	db := createDB()
	repository.InitDB(db)

	authService := infra.NewAuthService()
	drawerRepo := repository.NewDrawerRepo(db)
	drawerService := service.NewDrawerService(drawerRepo)
	drawerHandler := handler.NewDrawerHandler(drawerService, authService)
	connectionRepo := repository.NewConnectionRepo(db)
	connectionService := service.NewConnectionService(connectionRepo)
	connectionHandler := handler.NewConnectionHandler(connectionService)
	authHandler := handler.NewAuthHandler(authService)
	dataSourceRepo := repository.NewDataSourceRepo(db)
	tableService := infra.NewTableService()
	var bucket = initiateMongoBucket()
	fileService := infra.NewFileService(*bucket)
	dataSourceService := service.NewDataSourceService(dataSourceRepo, tableService, connectionService, fileService)
	dataSourceHandler := handler.NewDataSourceHandler(dataSourceService, authService)

	r.POST("/auth", authHandler.GetAuth)

	api := r.Group("/")

	api.Use(middleware.NewJWTMiddleware(authService).JWT())
	{
		api.GET("/drawers/:id", drawerHandler.Get)
		api.GET("/drawers", drawerHandler.GetAll)
		api.POST("/drawers", drawerHandler.Add)
		api.PUT("/drawers/:id", drawerHandler.Update)
		api.DELETE("/drawers/:id", drawerHandler.Delete)

		api.GET("/connections/:id", connectionHandler.Get)
		api.GET("/connections", connectionHandler.GetAll)
		api.POST("/connections", connectionHandler.Add)
		api.PUT("/connections/:id", connectionHandler.Update)
		api.DELETE("/connections/:id", connectionHandler.Delete)
		api.GET("/connections/:id/tables", dataSourceHandler.GetTables)
		api.GET("/connections/:id/tables/:tableName", dataSourceHandler.GetTableData)
		api.POST("/connections/:id/tables/:tableName/csv", dataSourceHandler.ConvertTableToCSV)

		api.GET("/dataSource/:id/content", dataSourceHandler.GetDataSourceContent)
	}
	return r
}

func createDB() *gorm.DB {
	var (
		dbName, user, password, host string
		port                         int
	)

	dbName = infra.DatabaseSetting.DBName
	user = infra.DatabaseSetting.DBUser
	password = infra.DatabaseSetting.DBPassword
	host = infra.DatabaseSetting.DBHost
	port = infra.DatabaseSetting.DBPort

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, port),
		PreferSimpleProtocol: true,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	}})

	if err != nil {
		log.Fatal("failed to open database")
	}

	return db
}

func initiateMongoBucket() *gridfs.Bucket {
	uri := infra.MongodbSetting.Uri
	opts := options.Client()
	credential := options.Credential{
		Username: infra.MongodbSetting.User,
		Password: infra.MongodbSetting.Password,
	}
	opts.ApplyURI(uri).SetAuth(credential)
	opts.SetMaxPoolSize(5)

	var err error
	var client *mongo.Client
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		log.Fatal("failed to open mongodb")
	}

	bucket, err := gridfs.NewBucket(
		client.Database("myfiles"),
	)
	if err != nil {
		log.Fatal("failed to get bucket")
	}
	return bucket
}
