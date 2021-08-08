package handler

import (
	"github.com/vcycyv/catwalk/infrastructure/mock"
	"github.com/vcycyv/catwalk/infrastructure/repository"
	"github.com/vcycyv/catwalk/service"
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
