package entity

type DataSource struct {
	Base

	DrawerID    string
	Description string
	User        string
	FileID      string

	Columns []Column
}

type Column struct {
	ID           string `gorm:"column:id;type:uuid;primary_key;"`
	Name         string
	DataSourceID string
}

func (DataSource) TableName() string {
	return "DataSource"
}

func (Column) TableName() string {
	return "Column"
}
