package service

import (
	"net/http"
	"time"

	"github.com/vcycyv/catwalk/assembler"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type serverService struct {
	serverRepo     domain.ServerRepository
	computeService domain.ComputeService
}

func NewServerService(serverRepo domain.ServerRepository, computeService domain.ComputeService) domain.ServerInterface {
	return &serverService{
		serverRepo,
		computeService,
	}
}

func (s *serverService) Add(server rep.Server) (*rep.Server, error) {
	maps := make(map[string]interface{})
	maps["host"] = server.Host
	maps["port"] = server.Port
	servers, err := s.GetAll(maps)
	if err != nil {
		return nil, err
	}
	if len(servers) > 0 {
		return nil, &rep.AppError{
			Code:    http.StatusConflict,
			Message: "The server is already registered.",
		}
	}

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

func (s *serverService) GetAll(maps interface{}) ([]*rep.Server, error) {
	servers, err := s.serverRepo.GetAll(maps)
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

func (s *serverService) GetAvailableServer() (*rep.Server, error) {
	count := 0
	for i := 1; i < 13; i++ {
		maps := make(map[string]interface{})
		maps["status"] = "available"
		maps["health"] = true
		servers, _ := s.serverRepo.GetAll(maps)
		for _, server := range servers {
			rtnVal := assembler.ServerAss.ToRepresentation(*server)
			if s.computeService.IsAlive(*rtnVal) {
				return rtnVal, nil
			}
		}

		time.Sleep(time.Second * 5)
		count += i
	}
	return nil, &rep.AppError{
		Code:    500,
		Message: "No compute server available.",
	}
}

func (s *serverService) RefreshHealth() {
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		maps := make(map[string]interface{})
		servers, _ := s.GetAll(maps)
		for _, server := range servers {
			alive := s.computeService.IsAlive(*server)
			if !server.Health == alive {
				server.Health = alive
				_, _ = s.serverRepo.Update(*assembler.ServerAss.ToData(*server))
			}
		}
	}
}
