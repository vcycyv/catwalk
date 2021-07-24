package service

import (
	"github.com/vcycyv/blog/assembler"
	"github.com/vcycyv/blog/domain"
	rep "github.com/vcycyv/blog/representation"
)

type connectionService struct {
	connectionRepo domain.ConnectionRepository
}

func NewConnectionService(connectionRepo domain.ConnectionRepository) domain.ConnectionInterface {
	return &connectionService{
		connectionRepo,
	}
}

func (s *connectionService) Add(connection rep.Connection) (*rep.Connection, error) {
	data, err := s.connectionRepo.Add(*assembler.ConnectionAss.ToData(connection))
	if err != nil {
		return &rep.Connection{}, err
	}
	return assembler.ConnectionAss.ToRepresentation(*data), nil
}

func (s *connectionService) Get(id string) (*rep.Connection, error) {
	data, err := s.connectionRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return assembler.ConnectionAss.ToRepresentation(*data), nil
}

func (s *connectionService) GetAll() ([]*rep.Connection, error) {
	connections, err := s.connectionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rtnVal := []*rep.Connection{}
	for _, connection := range connections {
		rtnVal = append(rtnVal, assembler.ConnectionAss.ToRepresentation(*connection))
	}
	return rtnVal, nil
}

func (s *connectionService) Update(connection rep.Connection) (*rep.Connection, error) {
	data, err := s.connectionRepo.Update(*assembler.ConnectionAss.ToData(connection))
	if err != nil {
		return nil, err
	}

	return assembler.ConnectionAss.ToRepresentation(*data), nil
}

func (s *connectionService) Delete(id string) error {
	err := s.connectionRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
