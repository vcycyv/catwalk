package domain

import (
	"github.com/vcycyv/blog/entity"
	rep "github.com/vcycyv/blog/representation"
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
