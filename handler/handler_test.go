package handler

import (
	"github.com/vcycyv/blog/infrastructure/mock"
	"github.com/vcycyv/blog/infrastructure/repository"
	"github.com/vcycyv/blog/service"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	drawerHdlr drawerHandler
)

func init() {
	authService := mock.NewMockAuth()

	db = mock.CreateDB()
	repository.InitDB(db)
	drawerRepo := repository.NewDrawerRepo(db)
	drawerService := service.NewDrawerService(drawerRepo)
	drawerHdlr = NewDrawerHandler(drawerService, authService)
}
