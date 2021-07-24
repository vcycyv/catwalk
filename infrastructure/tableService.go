package infrastructure

import (
	"fmt"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/vcycyv/blog/domain"
	rep "github.com/vcycyv/blog/representation"
)

type tableService struct{}

func NewTableService() domain.TableServiceInterface {
	return &tableService{}
}

func (s *tableService) GetTables(connection rep.Connection) ([]string, error) {
	var tables []string

	if strings.EqualFold(connection.Type, "mysql") {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", connection.User, connection.Password, connection.Host, connection.DbName))
		if err != nil {
			logrus.Errorf("failed to connect to db %s", connection.Name)
			return nil, err
		}
		defer db.Close()

		res, _ := db.Query("SHOW TABLES")
		var table string
		for res.Next() {
			err := res.Scan(&table)
			if err != nil {
				logrus.Errorf("failed to get table list for %s", connection.Name)
				return nil, err
			}
			tables = append(tables, table)
		}
	}

	return tables, nil
}

func (s *tableService) GetTableData(connection rep.Connection, table string) ([][]string, error) {
	var ret [][]string

	db, err := s.getDBConn(connection)
	if err != nil {
		logrus.Errorf("failed to open connetion %s", connection.Name)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("failed to get sqlDB connetion %s", connection.Name)
		return nil, err
	}
	defer sqlDB.Close()

	rows, err := db.Raw("select * from " + table).Rows()
	if err != nil {
		logrus.Errorf("failed to get table data %s", table)
		return nil, err
	}
	cols, _ := rows.Columns()
	ret = append(ret, cols)

	count := len(cols)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make([]string, count)

		for i := 0; i < count; i++ {
			var value interface{}
			rawValue := values[i]

			byteArray, ok := rawValue.([]byte)
			if ok {
				value = string(byteArray)
			} else {
				value = rawValue
			}

			if value == nil {
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", value)
			}
		}
		ret = append(ret, row)
	}

	return ret, nil
}

func (s *tableService) getDBConn(connection rep.Connection) (*gorm.DB, error) {
	if "mysql" == connection.Type {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			connection.User,
			connection.Password,
			connection.Host,
			connection.DbName)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	return nil, fmt.Errorf("unrecognized db type: %s", connection.Type)
}
