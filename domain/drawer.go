package domain

import (
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type DrawerRepository interface {
	Add(drawer entity.Drawer) (*entity.Drawer, error)
	Get(id string) (*entity.Drawer, error)
	GetAll() ([]*entity.Drawer, error)
	Update(drawer entity.Drawer) (*entity.Drawer, error)
	Delete(id string) error
}

type DrawerInterface interface {
	Add(drawer rep.Drawer) (*rep.Drawer, error)
	Get(id string) (*rep.Drawer, error)
	GetAll() ([]*rep.Drawer, error)
	Update(drawer rep.Drawer) (*rep.Drawer, error)
	Delete(id string) error
}
