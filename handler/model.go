package handler

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/infrastructure"
	"github.com/vcycyv/catwalk/infrastructure/util"
	rep "github.com/vcycyv/catwalk/representation"
)

type modelHandler struct {
	modelService     domain.ModelInterface
	modelFileService domain.ModelFileInterface
	authService      domain.AuthInterface
}

func NewModelHandler(modelService domain.ModelInterface,
	modelFileService domain.ModelFileInterface,
	authService domain.AuthInterface) modelHandler {
	return modelHandler{
		modelService,
		modelFileService,
		authService,
	}
}

func (s *modelHandler) Get(c *gin.Context) {
	id := c.Param("id")
	model, err := s.modelService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, model)
}

func (s *modelHandler) GetAll(c *gin.Context) {
	models, err := s.modelService.GetAll()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models)
}

func (s *modelHandler) AddModel(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		s.importMultiForm(c)
	} else {
		s.buildModel(c)
	}
}

func (s *modelHandler) AddFile(c *gin.Context) {
	modelID := c.Param("id")
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		_ = c.Error(err)
		return
	}
	form := c.Request.MultipartForm
	files := form.File["files"]
	var savedFiles []rep.ModelFile
	for _, file := range files {
		savedModelFile, err := s.addMartipartFileToModel(modelID, *file)
		if err != nil {
			_ = c.Error(err)
			return
		}
		savedFiles = append(savedFiles, *savedModelFile)
	}

	c.JSON(http.StatusOK, savedFiles)
}

func (s *modelHandler) Update(c *gin.Context) {
	id := c.Param("id")
	_, err := s.modelService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	model := &rep.Model{}
	if err := c.ShouldBind(model); err != nil {
		_ = c.Error(&rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		})
		return
	}
	logger.Debugf("Received request to add a model %s", model.Name)

	model, err = s.modelService.Update(*model)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The model %s is updated successfully.", model.Name)
	c.JSON(http.StatusOK, model)
}

func (s *modelHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := s.modelService.Delete(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (s *modelHandler) importMultiForm(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		_ = c.Error(err)
		return
	}
	form := c.Request.MultipartForm
	files := form.File["files"]

	model := rep.Model{}
	model.Name = c.Request.PostFormValue("name")
	model.Description = c.Request.PostFormValue("description")
	model.Function = c.Request.PostFormValue("function")

	savedModel, err := s.modelService.Add(model)
	if err != nil {
		_ = c.Error(err)
		return
	}

	for _, file := range files {
		savedModelFile, err := s.addMartipartFileToModel(savedModel.ID, *file)
		if err != nil {
			_ = c.Error(err)
			return
		}
		savedModel.Files = append(savedModel.Files, *savedModelFile)
	}

	c.JSON(http.StatusOK, savedModel)
}

// A payload example is below:
// {
// 	"trainTable": "d0095fa4-4d4b-4fc8-8394-6dcb5d686d08",
// 	"predictors": ["sepal_length","sepal_width","petal_length","petal_width"],
// 	"target": "species",
// 	"modelName":"iris9",
// 	"description":"test",
// 	"function": "classification"
// }
func (s *modelHandler) buildModel(c *gin.Context) {
	modelRequest := domain.BuildModelRequest{}
	if err := c.ShouldBind(&modelRequest); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	modelRequest.TrainTable = "http://" + util.GetOutboundIP() + ":" + strconv.Itoa(infrastructure.AppSetting.HTTPPort) + "/dataSources/" + modelRequest.TrainTable + "/content"
	model, err := s.modelService.BuildModel(modelRequest, s.authService.ExtractToken(c))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, model)
}

func (s *modelHandler) addMartipartFileToModel(modelID string, file multipart.FileHeader) (*rep.ModelFile, error) {
	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = openedFile.Close()
	}()

	modelFile := rep.ModelFile{}
	modelFile.Name = file.Filename
	modelFile.ModelID = modelID

	return s.modelFileService.Add(modelFile, openedFile)
}
