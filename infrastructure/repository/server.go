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

type serverRepo struct {
	db *gorm.DB
}

func NewServerRepo(db *gorm.DB) domain.ServerRepository {
	return &serverRepo{
		db,
	}
}

func (m *serverRepo) Add(server entity.Server) (*entity.Server, error) {
	logrus.Debugf("about to save a server %s", server.Name)
	if len(server.Status) == 0 {
		server.Status = "available"
	}

	if err := m.db.Create(&server).Error; err != nil {
		return nil, err
	}
	logrus.Debugf("server %s saved", server.Name)
	return &server, nil
}

func (m *serverRepo) Get(id string) (*entity.Server, error) {
	logrus.Debugf("about to get a server %s", id)
	var data entity.Server
	err := m.db.Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &representation.AppError{
				Code:    http.StatusFound,
				Message: fmt.Sprintf("Server %s is not found.", id),
			}
		}
		return &entity.Server{}, err
	}

	logrus.Debugf("server %s retrieved", id)
	return &data, err
}

func (m *serverRepo) GetAll(maps interface{}) ([]*entity.Server, error) {
	logrus.Debug("about to get all server")
	var servers []*entity.Server
	err := m.db.Where(maps).Find(&servers).Error
	if err != nil {
		return []*entity.Server{}, err
	}
	logrus.Debug("all server retrieved")
	return servers, nil
}

func (m *serverRepo) Update(server entity.Server) (*entity.Server, error) {
	logrus.Debugf("about to update a server %s", server.Name)
	err := m.db.Save(&server).Error
	logrus.Debugf("drawer %s updated", server.Name)
	return &server, err
}

func (m *serverRepo) Delete(id string) error {
	logrus.Debugf("about to delete a server %s", id)
	tx := m.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.Server{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	logrus.Debugf("server %s deleted", id)
	return tx.Commit().Error
}
