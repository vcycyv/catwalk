package domain

import (
	"io"

	"github.com/vcycyv/blog/entity"
	rep "github.com/vcycyv/blog/representation"
)

type DataSourceRepository interface {
	Add(drawer entity.DataSource) (*entity.DataSource, error)
	Get(id string) (*entity.DataSource, error)
	GetAll() ([]*entity.DataSource, error)
	Update(drawer entity.DataSource) (*entity.DataSource, error)
	Delete(id string) error
}

type DataSourceInterface interface {
	GetTables(connectionID string) ([]string, error)
	GetTableData(connectionID string, table string) ([][]string, error)
	ConvertTableToCSV(connectionID string, drawerID string, table string, user string) (*rep.DataSource, error)
	GetContent(id string, writer io.Writer) error
}
