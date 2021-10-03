package representation

type DataSource struct {
	Base

	DrawerID    string `json:"drawerId"`
	DrawerName  string `json:"drawerName"`
	Description string `json:"description"`
	User        string `json:"user"`
	FileID      string `json:"fileId"`

	Columns []string `json:"columns"`
}
