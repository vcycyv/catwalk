package repository

import (
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
