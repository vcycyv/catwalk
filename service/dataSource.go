package service

import (
	"github.com/vcycyv/blog/domain"
)

type dataSourceService struct {
	tableService      domain.TableServiceInterface
	connectionService domain.ConnectionInterface
}

func NewDataSourceService(tableService domain.TableServiceInterface,
	connectionService domain.ConnectionInterface) domain.DataSourceInterface {
	return &dataSourceService{
		tableService:      tableService,
		connectionService: connectionService,
	}
}

func (s *dataSourceService) GetTables(connectionID string) ([]string, error) {
	conn, _ := s.connectionService.Get(connectionID)
	return s.tableService.GetTables(*conn)
}

func (s *dataSourceService) GetTableData(connectionID string, table string) ([][]string, error) {
	conn, _ := s.connectionService.Get(connectionID)
	return s.tableService.GetTableData(*conn, table)
}
