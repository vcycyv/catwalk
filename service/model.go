package service

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/vcycyv/catwalk/assembler"
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/infrastructure"
	"github.com/vcycyv/catwalk/infrastructure/util"
	rep "github.com/vcycyv/catwalk/representation"
)

type modelService struct {
	modelRepo         domain.ModelRepository
	serverService     domain.ServerInterface
	computeService    domain.ComputeService
	fileService       domain.FileService
	dataSourceService domain.DataSourceInterface
}

func NewModelService(modelRepo domain.ModelRepository,
	serverService domain.ServerInterface,
	computeService domain.ComputeService,
	fileService domain.FileService,
	dataSourceService domain.DataSourceInterface) domain.ModelInterface {
	return &modelService{
		modelRepo,
		serverService,
		computeService,
		fileService,
		dataSourceService,
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

func (s *modelService) BuildModel(request domain.BuildModelRequest, token string) (*rep.Model, error) {
	server, err := s.serverService.GetAvailableServer()
	if err != nil {
		return nil, err
	}

	return s.computeService.BuildModel(*server, request, token)
}

func (s *modelService) Score(request domain.ScoreRequest, token string) (*rep.DataSource, error) {
	server, err := s.serverService.GetAvailableServer()
	if err != nil {
		return nil, err
	}

	model, err := s.Get(request.ModelID)
	if err != nil {
		return nil, err
	}

	for _, file := range model.Files {
		if s.isScoreFile(file) {
			tmpFile, _ := ioutil.TempFile("", "model_")
			err := s.fileService.DirectContentToWriter(file.FileID, tmpFile)
			if err != nil {
				return nil, err
			}
			request.ScoreFile = tmpFile
			scoreInputTable, err := s.dataSourceService.Get(request.ScoreInputTableID)
			if err != nil {
				return nil, err
			}
			request.ScoreInputTableURL = "http://" + util.GetOutboundIP() + ":" + strconv.Itoa(infrastructure.AppSetting.HTTPPort) + "/dataSources/" + scoreInputTable.ID

			rtnVal, err := s.computeService.Score(*server, request, token)
			if err != nil {
				return nil, err
			}

			return rtnVal, nil
		}
	}
	return nil, fmt.Errorf("No scoring file")
}

func (s *modelService) isScoreFile(modelFile rep.ModelFile) bool {
	if strings.EqualFold(modelFile.Role, "score") {
		return true
	}

	if strings.HasSuffix(modelFile.Name, ".pickle") {
		return true
	}

	return false
}
