package service

import (
	"io"

	"github.com/vcycyv/catwalk/assembler"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type modelFileService struct {
	modelFileRepo domain.ModelFileRepository
}

func NewModelFileService(modelFileRepo domain.ModelFileRepository) domain.ModelFileInterface {
	return &modelFileService{
		modelFileRepo,
	}
}

func (s *modelFileService) Add(modelFile rep.ModelFile, reader io.Reader) (*rep.ModelFile, error) {
	data, err := s.modelFileRepo.Add(*assembler.ModelFileAss.ToData(modelFile), reader)
	if err != nil {
		return &rep.ModelFile{}, err
	}
	return assembler.ModelFileAss.ToRepresentation(*data), nil
}

func (s *modelFileService) Get(id string) (*rep.ModelFile, error) {
	data, err := s.modelFileRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return assembler.ModelFileAss.ToRepresentation(*data), nil
}

func (s *modelFileService) GetAll() ([]*rep.ModelFile, error) {
	modelFiles, err := s.modelFileRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rtnVal := []*rep.ModelFile{}
	for _, modelFile := range modelFiles {
		rtnVal = append(rtnVal, assembler.ModelFileAss.ToRepresentation(*modelFile))
	}
	return rtnVal, nil
}

func (s *modelFileService) Update(modelFile rep.ModelFile) (*rep.ModelFile, error) {
	data, err := s.modelFileRepo.Update(*assembler.ModelFileAss.ToData(modelFile))
	if err != nil {
		return nil, err
	}

	return assembler.ModelFileAss.ToRepresentation(*data), nil
}

func (s *modelFileService) Delete(id string) error {
	err := s.modelFileRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
