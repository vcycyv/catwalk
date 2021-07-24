package domain

type DataSourceInterface interface {
	GetTables(connectionID string) ([]string, error)
	GetTableData(connectionID string, table string) ([][]string, error)
}
