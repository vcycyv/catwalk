package repository

import (
	"database/sql"
	"reflect"

	logger "github.com/sirupsen/logrus"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/vcycyv/catwalk/entity"
)

func InitDB(db *gorm.DB) {
	var createCallback = func(db *gorm.DB) {
		idField := db.Statement.Schema.LookUpField("id")
		if idField != nil {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					_ = idField.Set(db.Statement.ReflectValue.Index(i), uuid.New().String())
				}
			case reflect.Struct:
				_ = idField.Set(db.Statement.ReflectValue, uuid.New().String())
			}
		}
	}

	err := db.Callback().Create().Before("gorm:create").Register("uuid", createCallback)
	if err != nil {
		logger.Fatal("failed to register uuid hook")
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)

	migrate(db)

	migrateForFolderService(sqlDB)
}

func migrate(db *gorm.DB) {
	_ = db.AutoMigrate(&entity.Drawer{})
	_ = db.AutoMigrate(&entity.Connection{})
	_ = db.AutoMigrate(&entity.DataSource{})
	_ = db.AutoMigrate(&entity.Column{})
	_ = db.AutoMigrate(&entity.Server{})
	_ = db.AutoMigrate(&entity.ModelFile{})
	_ = db.AutoMigrate(&entity.Model{})
}

func migrateForFolderService(sqlDB *sql.DB) {
	_, err := sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")
	if err != nil {
		logger.Fatal("failed to create extension")
	}

	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS ltree")
	if err != nil {
		logger.Fatal("failed to create extension")
	}

	_, err = sqlDB.Exec("CREATE TABLE IF NOT EXISTS Folder (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), path ltree)")
	if err != nil {
		logger.Fatal("failed to create folder table")
	}

	_, err = sqlDB.Exec("CREATE INDEX IF NOT EXISTS path_gist_idx ON folder USING gist(path)")
	if err != nil {
		logger.Warn("failed to create gist index")
	}

	_, err = sqlDB.Exec("CREATE INDEX IF NOT EXISTS path_idx ON folder USING btree(path)")
	if err != nil {
		logger.Warn("failed to create btree index")
	}
}
