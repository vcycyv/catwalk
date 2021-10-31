package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/entity"
	infra "github.com/vcycyv/catwalk/infrastructure"
)

type folderRepo struct {
	conn *sqlx.DB
}

func NewFolderRepo() domain.FolderService {
	var (
		dbName, user, password, host string
		port                         int
	)

	dbName = infra.DatabaseSetting.DBName
	user = infra.DatabaseSetting.DBUser
	password = infra.DatabaseSetting.DBPassword
	host = infra.DatabaseSetting.DBHost
	port = infra.DatabaseSetting.DBPort

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return &folderRepo{
		conn: db,
	}
}

//Create create folder record
func (r *folderRepo) Create(f entity.Folder) (*entity.Folder, error) {
	query := "INSERT INTO folder (path) VALUES ($1) RETURNING *"
	var folder entity.Folder
	err := r.conn.QueryRowx(query, f.Path).StructScan(&folder)
	if err != nil {
		return nil, err
	}

	return r.GetByID(folder.ID)
}

func (r *folderRepo) GetAll() ([]entity.Folder, error) {
	query := "Select * From folder Order By path"
	rtnVal, err := getCollection(query, r)
	if err != nil {
		return nil, err
	}

	return rtnVal, nil
}

//GetByID ...
func (r *folderRepo) GetByID(id string) (*entity.Folder, error) {
	query := "Select * From folder where id = $1"
	rtnVal := entity.Folder{}
	row := r.conn.QueryRowx(query, id)
	err := row.StructScan(&rtnVal)
	if err != nil || len(rtnVal.Path) == 0 {
		log.Errorf("failed to get folder by id %s", id)
		return nil, nil
	}

	lastIndex := strings.LastIndex(rtnVal.Path, ".")
	if lastIndex != -1 {
		parent, err := r.GetByPath(rtnVal.Path[0:lastIndex])
		if err != nil {
			return nil, err
		}
		rtnVal.ParentID = parent.ID
		rtnVal.Name = rtnVal.Path[lastIndex+1:]
	} else {
		rtnVal.Name = rtnVal.Path
	}

	return &rtnVal, err
}

//GetByPath ...
func (r *folderRepo) GetByPath(path string) (*entity.Folder, error) {
	query := "Select * From folder where path = $1"
	rtnVal := entity.Folder{}
	row := r.conn.QueryRowx(query, path)
	err := row.StructScan(&rtnVal)

	return &rtnVal, err
}

//GetChildren get children by parent
func (r *folderRepo) GetChildren(parent string) ([]entity.Folder, error) {
	query := "SELECT id, path FROM folder WHERE path ~ '" + parent + ".*{1}' Order By path"
	rtnVal, err := getCollection(query, r)
	if err != nil {
		return nil, err
	}

	return rtnVal, nil
}

//GetDescendants get descendants by parent path
func (r *folderRepo) GetDescendants(parentPath string) ([]entity.Folder, error) {
	query := "SELECT * FROM folder WHERE '" + parentPath + "' @> path Order By path"
	rtnVal, err := getCollection(query, r)
	if err != nil {
		return nil, err
	}

	return rtnVal, nil
}

//Delete delete by path
func (r *folderRepo) Delete(path string) error {
	query := "delete from folder where '" + path + "' @> path"

	return executeStatement(r.conn, query)
}

//Rename rename a folder by id and new name
func (r *folderRepo) Rename(id string, name string) (*entity.Folder, error) {
	folder, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	newName := folder.Path[0:strings.LastIndex(folder.Path, ".")+1] + name

	rootQuery := "update folder set path = '" + newName + "' where path = '" + folder.Path + "'"

	tx, _ := r.conn.Beginx()

	err = executeStatement(tx, rootQuery)
	if err != nil {
		return nil, err
	}

	descendantsQuery := "update folder set path = '" + newName + "' || subpath(path,nlevel(text2ltree('" + folder.Path + "'))) where path ~ '" + folder.Path + ".*' and path != '" + folder.Path + "'"
	err = executeStatement(tx, descendantsQuery)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to rename folder %s", id)
		return nil, err
	}

	return r.GetByID(id)
}

//Move ...
func (r *folderRepo) Move(sourceID string, targetID string) (*entity.Folder, error) {
	sourceFolder, err := r.GetByID(sourceID)
	if err != nil {
		return nil, err
	}
	targetFolder, err := r.GetByID(targetID)
	if err != nil {
		return nil, err
	}

	source := sourceFolder.Path
	target := targetFolder.Path
	sourceParent := source[0:strings.LastIndex(source, ".")]

	descendants, err := r.GetDescendants(source)
	if err != nil {
		return nil, err
	}

	var stmt *sqlx.Stmt
	tx, _ := r.conn.Beginx()
	movingRoot := target + source[strings.LastIndex(source, "."):]
	query := "update folder set path = '" + movingRoot + "' where path = '" + source + "'"
	stmt, err = tx.Preparex(query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Errorf("failed to update folder during folder move %s", source)
		return nil, err
	}

	for _, descendant := range descendants {
		newPath := target + "." + descendant.Path[len(sourceParent)+1:len(descendant.Path)]
		query = "update folder set path = '" + newPath + "' where id = '" + descendant.ID + "'"
		stmt, err = tx.Preparex(query)
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Errorf("failed to update folder during moving folder to %s", descendant.ID)
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit transaction during moving folder")
		return nil, err
	}

	rtnVal, err := r.GetByPath(movingRoot)
	if err != nil {
		return nil, err
	}
	rtnVal, _ = r.GetByID(rtnVal.ID)

	return rtnVal, nil
}

func getCollection(query string, r *folderRepo) ([]entity.Folder, error) {
	rows, err := r.conn.Queryx(query)
	if err != nil {
		return nil, err
	}

	rtnVal := make([]entity.Folder, 0)
	for rows.Next() {
		data := entity.Folder{}

		err := rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		rtnVal = append(rtnVal, data)
	}
	return rtnVal, nil
}

func executeStatement(conn Preparable, query string) error {
	stmt, err := conn.Preparex(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

//Preparable ...
type Preparable interface {
	Preparex(query string) (*sqlx.Stmt, error)
}
