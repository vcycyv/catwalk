package domain

import (
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type ConnectionRepository interface {
	Add(connection entity.Connection) (*entity.Connection, error)
	Get(id string) (*entity.Connection, error)
	GetAll() ([]*entity.Connection, error)
	Update(connection entity.Connection) (*entity.Connection, error)
	Delete(id string) error
}

type ConnectionInterface interface {
	Add(connection rep.Connection) (*rep.Connection, error)
	Get(id string) (*rep.Connection, error)
	GetAll() ([]*rep.Connection, error)
	Update(connection rep.Connection) (*rep.Connection, error)
	Delete(id string) error
}
