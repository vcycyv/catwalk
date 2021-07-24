package service

import (
	"github.com/vcycyv/blog/assembler"
	"github.com/vcycyv/blog/domain"
	rep "github.com/vcycyv/blog/representation"
)

type drawerService struct {
	drawerRepo domain.DrawerRepository
}

func NewDrawerService(drawerRepo domain.DrawerRepository) domain.DrawerInterface {
	return &drawerService{
		drawerRepo,
	}
}

func (s *drawerService) Add(drawer rep.Drawer) (*rep.Drawer, error) {
	data, err := s.drawerRepo.Add(*assembler.DrawerAss.ToData(drawer))
	if err != nil {
		return &rep.Drawer{}, err
	}
	return assembler.DrawerAss.ToRepresentation(*data), nil
}

func (s *drawerService) Get(id string) (*rep.Drawer, error) {
	data, err := s.drawerRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return assembler.DrawerAss.ToRepresentation(*data), nil
}

func (s *drawerService) GetAll() ([]*rep.Drawer, error) {
	drawers, err := s.drawerRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rtnVal := []*rep.Drawer{}
	for _, drawer := range drawers {
		rtnVal = append(rtnVal, assembler.DrawerAss.ToRepresentation(*drawer))
	}
	return rtnVal, nil
}

func (s *drawerService) Update(drawer rep.Drawer) (*rep.Drawer, error) {
	data, err := s.drawerRepo.Update(*assembler.DrawerAss.ToData(drawer))
	if err != nil {
		return nil, err
	}

	return assembler.DrawerAss.ToRepresentation(*data), nil
}

func (s *drawerService) Delete(id string) error {
	err := s.drawerRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
