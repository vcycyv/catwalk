package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vcycyv/blog/domain"
)

type dataSourceHandler struct {
	dataSourceService domain.DataSourceInterface
}

func NewDataSourceHandler(dataSourceService domain.DataSourceInterface) *dataSourceHandler {
	return &dataSourceHandler{
		dataSourceService: dataSourceService,
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
