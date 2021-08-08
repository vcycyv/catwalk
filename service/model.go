package service

import (
	"github.com/vcycyv/catwalk/assembler"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type modelService struct {
	modelRepo domain.ModelRepository
}

func NewModelService(modelRepo domain.ModelRepository) domain.ModelInterface {
	return &modelService{
		modelRepo,
	}
}

func (s *modelService) Add(model rep.Model) (*rep.Model, error) {
	data, err := s.modelRepo.Add(*assembler.ModelAss.ToData(model))
	if err != nil {
		return &rep.Model{}, err
	}
	return assembler.ModelAss.ToRepresentation(*data), nil
}

func (s *modelService) Get(id string) (*rep.Model, error) {
	data, err := s.modelRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return assembler.ModelAss.ToRepresentation(*data), nil
}

func (s *modelService) GetAll() ([]*rep.Model, error) {
	models, err := s.modelRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rtnVal := []*rep.Model{}
	for _, model := range models {
		rtnVal = append(rtnVal, assembler.ModelAss.ToRepresentation(*model))
	}
	return rtnVal, nil
}

func (s *modelService) Update(model rep.Model) (*rep.Model, error) {
	data, err := s.modelRepo.Update(*assembler.ModelAss.ToData(model))
	if err != nil {
		return nil, err
	}

	return assembler.ModelAss.ToRepresentation(*data), nil
}

func (s *modelService) Delete(id string) error {
	err := s.modelRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}