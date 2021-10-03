package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/handler"
	infra "github.com/vcycyv/catwalk/infrastructure"
	"github.com/vcycyv/catwalk/infrastructure/repository"
	"github.com/vcycyv/catwalk/middleware"
	"github.com/vcycyv/catwalk/service"
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
		Addr:         fmt.Sprintf(":%d", infra.AppSetting.HTTPPort),
		Handler:      router,
		ReadTimeout:  time.Duration(60) * time.Second,
		WriteTimeout: time.Duration(60) * time.Second,
		//MaxHeaderBytes: 32 << 22,
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
	var bucket = initiateMongoBucket()
	fileService := infra.NewFileService(*bucket)

	drawerRepo := repository.NewDrawerRepo(db)
	drawerService := service.NewDrawerService(drawerRepo)
	drawerHandler := handler.NewDrawerHandler(drawerService, authService)

	connectionRepo := repository.NewConnectionRepo(db)
	connectionService := service.NewConnectionService(connectionRepo)
	connectionHandler := handler.NewConnectionHandler(connectionService)

	authHandler := handler.NewAuthHandler(authService)

	dataSourceRepo := repository.NewDataSourceRepo(db, fileService)
	tableService := infra.NewTableService()
	dataSourceService := service.NewDataSourceService(dataSourceRepo, tableService, drawerService, connectionService, fileService)
	dataSourceHandler := handler.NewDataSourceHandler(dataSourceService, authService)

	computeService := infra.NewComputeService()
	serverRepo := repository.NewServerRepo(db)
	serverService := service.NewServerService(serverRepo, computeService)
	serverHandler := handler.NewServerHandler(serverService)

	go serverService.RefreshHealth()

	modelFileRepo := repository.NewModelFileRepo(db, fileService)
	modelFileService := service.NewModelFileService(modelFileRepo)

	modelRepo := repository.NewModelRepo(db)
	modelService := service.NewModelService(modelRepo, serverService, computeService)
	modelHandler := handler.NewModelHandler(modelService, modelFileService, authService)

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

		api.POST("/dataSources", dataSourceHandler.Add)
		api.GET("/dataSources/:id", dataSourceHandler.Get)
		api.GET("/dataSources", dataSourceHandler.GetAll)
		api.PUT("/dataSources/:id", dataSourceHandler.Update)
		api.GET("/dataSources/:id/content", dataSourceHandler.GetContent)
		api.DELETE("/dataSources/:id", dataSourceHandler.Delete)

		api.GET("/servers/:id", serverHandler.Get)
		api.GET("/servers", serverHandler.GetAll)
		api.POST("/servers", serverHandler.Add)
		api.DELETE("/servers/:id", serverHandler.Delete)

		api.GET("/models/:id", modelHandler.Get)
		api.GET("/models", modelHandler.GetAll)
		api.POST("/models", modelHandler.AddModel)
		api.POST("/models/:id/files", modelHandler.AddFile)
		api.PUT("/models/:id", modelHandler.Update)
		api.DELETE("/models/:id", modelHandler.Delete)
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
