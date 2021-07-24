package service

import (
	"github.com/vcycyv/blog/domain"
	"github.com/vcycyv/blog/infrastructure/mock"
	"github.com/vcycyv/blog/infrastructure/repository"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	drawerSvc domain.DrawerInterface
)

func init() {
	db = mock.CreateDB()
	repository.InitDB(db)
	drawerRepo := repository.NewDrawerRepo(db)
	drawerSvc = NewDrawerService(drawerRepo)
}
