package representation

type DataSource struct {
	Base

	DrawerID    string `json:"drawerId"`
	Description string `json:"description"`
	User        string `json:"user"`
	FileID      string `json:"fileId"`

	Columns []string `json:"columns"`
}
