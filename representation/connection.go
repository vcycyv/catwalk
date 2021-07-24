package representation

type Connection struct {
	Base

	Type     string `json:"type"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
}
