package repository

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/entity"
	"github.com/vcycyv/catwalk/representation"
	"gorm.io/gorm"
)

type connectionRepo struct {
	db *gorm.DB
}

func NewConnectionRepo(db *gorm.DB) domain.ConnectionRepository {
	return &connectionRepo{
		db,
	}
}

func (m *connectionRepo) Add(connection entity.Connection) (*entity.Connection, error) {
	logrus.Debugf("about to save a connection %s", connection.Name)
	if err := m.db.Create(&connection).Error; err != nil {
		return nil, err
	}
	logrus.Debugf("connection %s saved", connection.Name)
	return &connection, nil
}

func (m *connectionRepo) Get(id string) (*entity.Connection, error) {
	logrus.Debugf("about to get a connection %s", id)
	var data entity.Connection
	err := m.db.Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &representation.AppError{
				Code:    http.StatusFound,
				Message: fmt.Sprintf("Connection %s is not found.", id),
			}
		}
		return &entity.Connection{}, err
	}

	logrus.Debugf("connection %s retrieved", id)
	return &data, err
}

func (m *connectionRepo) GetAll() ([]*entity.Connection, error) {
	logrus.Debug("about to get all connection")
	var connections []*entity.Connection
	err := m.db.Find(&connections).Error
	if err != nil {
		return []*entity.Connection{}, err
	}
	logrus.Debug("all connection retrieved")
	return connections, nil
}

func (m *connectionRepo) Update(connection entity.Connection) (*entity.Connection, error) {
	logrus.Debugf("about to update a connection %s", connection.Name)
	err := m.db.Save(&connection).Error
	logrus.Debugf("connection %s updated", connection.Name)
	return &connection, err
}

func (m *connectionRepo) Delete(id string) error {
	logrus.Debugf("about to delete a connection %s", id)
	tx := m.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.Connection{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	logrus.Debugf("connection %s deleted", id)
	return tx.Commit().Error
}
