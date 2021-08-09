package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type computeService struct{}

func NewComputeService() domain.ComputeService {
	return &computeService{}
}

func (s *computeService) IsAlive(server rep.Server) bool {
	resp, err := http.Get("http://" + server.Host + ":" + strconv.Itoa(server.Port) + "/status")
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}
