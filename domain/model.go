package domain

import (
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type ModelRepository interface {
	Add(model entity.Model) (*entity.Model, error)
	Get(id string) (*entity.Model, error)
	GetAll() ([]*entity.Model, error)
	Update(model entity.Model) (*entity.Model, error)
	Delete(id string) error
}

type ModelInterface interface {
	Add(model rep.Model) (*rep.Model, error)
	Get(id string) (*rep.Model, error)
	GetAll() ([]*rep.Model, error)
	Update(model rep.Model) (*rep.Model, error)
	Delete(id string) error
	BuildModel(request BuildModelRequest, token string) (*rep.Model, error)
	Score(request ScoreRequest, token string) (*rep.DataSource, error)
}
