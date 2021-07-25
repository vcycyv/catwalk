package entity

type Connection struct {
	Base

	Type     string
	Host     string
	User     string
	Password string
	DbName   string
}
