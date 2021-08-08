package entity

type Server struct {
	Base

	Host   string
	Port   int
	Status string
	Health bool
}

func (Server) TableName() string {
	return "Server"
}
