package infrastructure

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

func (s *computeService) BuildModel(server rep.Server, buildModelRequest domain.BuildModelRequest, token string) (*rep.Model, error) {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(buildModelRequest)
	req, _ := http.NewRequest("POST", "http://"+server.Host+":"+strconv.Itoa(server.Port)+"/models", buffer)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &rep.AppError{
			Code:    500,
			Message: "failed to send request during building model",
		}
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return nil, &rep.AppError{
			Code:    500,
			Message: "failed to build model: " + string(body),
		}
	}

	model := &rep.Model{}
	err = json.Unmarshal(body, model)
	if err != nil {
		return nil, &rep.AppError{
			Code:    500,
			Message: "failed to unmarshal response during building model: " + string(body),
		}
	}

	return model, nil
}
