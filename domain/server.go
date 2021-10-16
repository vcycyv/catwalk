package domain

import (
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type ServerRepository interface {
	Add(server entity.Server) (*entity.Server, error)
	Get(id string) (*entity.Server, error)
	GetAll(maps interface{}) ([]*entity.Server, error)
	Update(server entity.Server) (*entity.Server, error)
	Delete(id string) error
}

type ServerInterface interface {
	Add(server rep.Server) (*rep.Server, error)
	Get(id string) (*rep.Server, error)
	GetAll(maps interface{}) ([]*rep.Server, error)
	Delete(id string) error
	GetAvailableServer() (*rep.Server, error)
	RefreshHealth()
}

type BuildModelRequest struct {
	TrainTable  string   `json:"trainTable"`
	Predictors  []string `json:"predictors"`
	Target      string   `json:"target"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Function    string   `json:"function"`
	Algorithm   string   `json:"algorithm"`
}
