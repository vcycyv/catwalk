package service

import (
	"github.com/vcycyv/catwalk/assembler"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type serverService struct {
	serverRepo domain.ServerRepository
}

func NewServerService(serverRepo domain.ServerRepository) domain.ServerInterface {
	return &serverService{
		serverRepo,
	}
}

func (s *serverService) Add(server rep.Server) (*rep.Server, error) {
	data, err := s.serverRepo.Add(*assembler.ServerAss.ToData(server))
	if err != nil {
		return &rep.Server{}, err
	}
	return assembler.ServerAss.ToRepresentation(*data), nil
}

func (s *serverService) Get(id string) (*rep.Server, error) {
	data, err := s.serverRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return assembler.ServerAss.ToRepresentation(*data), nil
}

func (s *serverService) GetAll() ([]*rep.Server, error) {
	servers, err := s.serverRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rtnVal := []*rep.Server{}
	for _, server := range servers {
		rtnVal = append(rtnVal, assembler.ServerAss.ToRepresentation(*server))
	}
	return rtnVal, nil
}

func (s *serverService) Delete(id string) error {
	err := s.serverRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
