package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	rep "github.com/vcycyv/blog/representation"
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

type AuthInterface interface {
	Auth(user string, password string) error
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (*Claims, error)
	GetUserFromToken(token string) (string, error)
	ExtractToken(c *gin.Context) string
}

type TableServiceInterface interface {
	GetTables(connection rep.Connection) ([]string, error)
	GetTableData(connection rep.Connection, table string) ([][]string, error)
}

type FileService interface {
}
