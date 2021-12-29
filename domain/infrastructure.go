package domain

import (
	"io"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
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
	ConvertTableToCSV(connection rep.Connection, table string, writer io.Writer) ([]string, error)
}

type FileService interface {
	Save(fileName string, reader io.Reader) (string, error)
	DirectContentToWriter(fileID string, writer io.Writer) error //TODO check if it makes sense for all consumers
	Delete(fileID string) error
	GetContent(fileID string) (io.Reader, error)
}

type FolderService interface {
	Create(f entity.Folder) (*entity.Folder, error)
	GetAll() ([]entity.Folder, error)
	GetByID(id string) (*entity.Folder, error)
	GetByPath(path string) (*entity.Folder, error)
	GetChildren(parent string) ([]entity.Folder, error)
	GetDescendants(parentPath string) ([]entity.Folder, error)
	Rename(id string, name string) (*entity.Folder, error)
	Delete(path string) error
}

type ComputeService interface {
	IsAlive(server rep.Server) bool
	BuildModel(server rep.Server, buildModelRequest BuildModelRequest, token string) (*rep.Model, error)
	Score(server rep.Server, scoreRequest ScoreRequest, token string) (*rep.DataSource, error)
}
