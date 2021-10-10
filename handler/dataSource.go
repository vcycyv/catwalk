package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type dataSourceHandler struct {
	dataSourceService domain.DataSourceInterface
	authService       domain.AuthInterface
}

func NewDataSourceHandler(dataSourceService domain.DataSourceInterface, authService domain.AuthInterface) *dataSourceHandler {
	return &dataSourceHandler{
		dataSourceService,
		authService,
	}
}

func (s *dataSourceHandler) GetTables(c *gin.Context) {
	connectionID := c.Param("id")
	tables, err := s.dataSourceService.GetTables(connectionID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tables)
}

func (s *dataSourceHandler) GetTableData(c *gin.Context) {
	connectionID := c.Param("id")
	tableName := c.Param("tableName")
	data, err := s.dataSourceService.GetTableData(connectionID, tableName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *dataSourceHandler) ConvertTableToCSV(c *gin.Context) {
	connectionID := c.Param("id")
	tableName := c.Param("tableName")
	drawerID := c.Query("drawerId")

	token := s.authService.ExtractToken(c)
	user, _ := s.authService.GetUserFromToken(token)

	dataSource, err := s.dataSourceService.ConvertTableToCSV(connectionID, drawerID, tableName, user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dataSource)
}

func (s *dataSourceHandler) GetContent(c *gin.Context) {
	dataSourceID := c.Param("id")
	err := s.dataSourceService.GetContent(dataSourceID, c.Writer)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

func (s *dataSourceHandler) Add(c *gin.Context) {
	file, _ := c.FormFile("file")
	drawerID := c.Request.FormValue("drawerId")
	if drawerID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "drawer ID must not be empty.",
		})
		return
	}
	fileName := file.Filename

	token := s.authService.ExtractToken(c)
	user, _ := s.authService.GetUserFromToken(token)

	openedFile, err := file.Open()
	defer func() {
		openedFile.Close()
	}()

	if err != nil {
		_ = c.Error(err)
		return
	}

	rtnVal, err := s.dataSourceService.Add(drawerID, fileName, user, openedFile)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, rtnVal)
}

func (s *dataSourceHandler) Get(c *gin.Context) {
	id := c.Param("id")
	dataSource, err := s.dataSourceService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dataSource)
}

func (s *dataSourceHandler) GetAll(c *gin.Context) {
	dataSourceArray, err := s.dataSourceService.GetAll()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dataSourceArray)
}

func (s *dataSourceHandler) Update(c *gin.Context) {
	id := c.Param("id")
	_, err := s.dataSourceService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	dataSource := &rep.DataSource{}
	if err := c.ShouldBind(dataSource); err != nil {
		_ = c.Error(&rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		})
		return
	}
	logger.Debugf("Received request to add a dataSource %s", dataSource.Name)

	dataSource, err = s.dataSourceService.Update(*dataSource)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The dataSource %s is updated successfully.", dataSource.Name)
	c.JSON(http.StatusOK, dataSource)
}

func (s *dataSourceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	logger.Debugf("Received request to delete a dataSource %s.", id)
	err := s.dataSourceService.Delete(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The dataSource %s is deleted successfully.", id)
	c.JSON(http.StatusNoContent, nil)
}

func (s *dataSourceHandler) GetColumns(c *gin.Context) {
	id := c.Param("id")
	rtnVal, err := s.dataSourceService.GetColumns(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, rtnVal)
}
