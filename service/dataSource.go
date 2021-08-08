package service

import (
	"io"

	"github.com/vcycyv/catwalk/assembler"
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
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
		dataSourceRepo,
		tableService,
		connectionService,
		fileService,
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
	var columns []string
	go func() {
		defer writer.Close()
		var err error
		columns, err = s.tableService.ConvertTableToCSV(*conn, table, writer)
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	id, err := s.fileService.Save(table+".csv", reader)
	if err != nil {
		return nil, err
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

	var columnArrayData []entity.Column
	for _, column := range columns {
		columnArrayData = append(columnArrayData, entity.Column{Name: column})
	}
	data := entity.DataSource{
		Base:     entity.Base{Name: table},
		DrawerID: drawerID,
		User:     user,
		FileID:   id,
		Columns:  columnArrayData,
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

func (s *dataSourceService) Add(drawerID string, tableName string, user string, reader io.ReadCloser) (*rep.DataSource, error) {
	fileId, err := s.fileService.Save(tableName, reader)
	if err != nil {
		return nil, err
	}

	data := &entity.DataSource{
		Base:     entity.Base{Name: tableName},
		DrawerID: drawerID,
		User:     user,
		FileID:   fileId,
	}

	data, err = s.dataSourceRepo.Add(*data)
	if err != nil {
		return nil, err
	}
	return assembler.DataSourceAss.ToRepresentation(*data), nil
}

func (s *dataSourceService) Get(id string) (*rep.DataSource, error) {
	data, err := s.dataSourceRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return assembler.DataSourceAss.ToRepresentation(*data), nil
}

func (s *dataSourceService) GetAll() ([]*rep.DataSource, error) {
	dataSources, err := s.dataSourceRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rtnVal := []*rep.DataSource{}
	for _, dataSource := range dataSources {
		rtnVal = append(rtnVal, assembler.DataSourceAss.ToRepresentation(*dataSource))
	}
	return rtnVal, nil
}

func (s *dataSourceService) Update(dataSource rep.DataSource) (*rep.DataSource, error) {
	data, err := s.dataSourceRepo.Update(*assembler.DataSourceAss.ToData(dataSource))
	if err != nil {
		return nil, err
	}

	return assembler.DataSourceAss.ToRepresentation(*data), nil
}

func (s *dataSourceService) Delete(id string) error {
	err := s.dataSourceRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
