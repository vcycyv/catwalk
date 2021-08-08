package service

import (
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/infrastructure/mock"
	"github.com/vcycyv/catwalk/infrastructure/repository"
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
