package service

import (
	"io"

	"github.com/vcycyv/blog/assembler"
	"github.com/vcycyv/blog/domain"
	"github.com/vcycyv/blog/entity"
	rep "github.com/vcycyv/blog/representation"
)

type dataSourceService struct {
	dataSourceRepo    domain.DataSourceRepository
	tableService      domain.TableServiceInterface
	connectionService domain.ConnectionInterface
	fileService       domain.FileService
}

func NewDataSourceService(dataSourceRepo domain.DataSourceRepository,
	tableService domain.TableServiceInterface,
	connectionService domain.ConnectionInterface,
	fileService domain.FileService) domain.DataSourceInterface {
	return &dataSourceService{
		dataSourceRepo:    dataSourceRepo,
		tableService:      tableService,
		connectionService: connectionService,
		fileService:       fileService,
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

func (s *dataSourceService) ConvertTableToCSV(connectionID string, drawerID string, table string, user string) (*rep.DataSource, error) {
	conn, _ := s.connectionService.Get(connectionID)

	reader, writer := io.Pipe()
	defer reader.Close()

	errChan := make(chan error, 1)
	go func() {
		defer writer.Close()
		err := s.tableService.ConvertTableToCSV(*conn, table, writer)
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	id, err := s.fileService.Save(reader, table)
	if err != nil {
		return nil, err
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

	data := entity.DataSource{
		Base:     entity.Base{Name: table},
		DrawerID: drawerID,
		User:     user,
		FileID:   id,
	}

	newDataSource, err := s.dataSourceRepo.Add(data)
	if err != nil {
		return nil, err
	}

	return assembler.DataSourceAss.ToRepresentation(*newDataSource), nil
}

func (s *dataSourceService) GetContent(id string, writer io.Writer) error {
	dataSourceData, err := s.dataSourceRepo.Get(id)
	if err != nil {
		return err
	}
	return s.fileService.GetContent(dataSourceData.FileID, writer)
}
