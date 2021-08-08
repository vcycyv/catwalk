package domain

import (
	"io"

	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type ModelFileRepository interface {
	Add(modelFile entity.ModelFile, reader io.Reader) (*entity.ModelFile, error)
	Get(id string) (*entity.ModelFile, error)
	GetContent(id string, writer io.Writer) error
	GetAll() ([]*entity.ModelFile, error)
	Update(modelFile entity.ModelFile) (*entity.ModelFile, error)
	Delete(id string) error
}

type ModelFileInterface interface {
	Add(modelFile rep.ModelFile, reader io.Reader) (*rep.ModelFile, error)
	Get(id string) (*rep.ModelFile, error)
	GetAll() ([]*rep.ModelFile, error)
	Update(modelFile rep.ModelFile) (*rep.ModelFile, error)
	Delete(id string) error
}
