package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcycyv/blog/domain"
)

type dataSourceHandler struct {
	dataSourceService domain.DataSourceInterface
	authService       domain.AuthInterface
}

func NewDataSourceHandler(dataSourceService domain.DataSourceInterface, authService domain.AuthInterface) *dataSourceHandler {
	return &dataSourceHandler{
		dataSourceService: dataSourceService,
		authService:       authService,
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

func (s *dataSourceHandler) GetDataSourceContent(c *gin.Context) {
	dataSourceID := c.Param("id")
	err := s.dataSourceService.GetContent(dataSourceID, c.Writer)
	if err != nil {
		_ = c.Error(err)
		return
	}
}
