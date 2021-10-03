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

type dataSourceRepo struct {
	db          *gorm.DB
	fileService domain.FileService
}

func NewDataSourceRepo(db *gorm.DB, fileService domain.FileService) domain.DataSourceRepository {
	return &dataSourceRepo{
		db,
		fileService,
	}
}

func (m *dataSourceRepo) Add(dataSource entity.DataSource) (*entity.DataSource, error) {
	logrus.Debugf("about to save a dataSource %s", dataSource.Name)
	if err := m.db.Create(&dataSource).Error; err != nil {
		return nil, err
	}
	logrus.Debugf("dataSource %s saved", dataSource.Name)
	return &dataSource, nil
}

func (m *dataSourceRepo) Get(id string) (*entity.DataSource, error) {
	logrus.Debugf("about to get a dataSource %s", id)
	var data entity.DataSource
	err := m.db.Preload("Columns").Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &representation.AppError{
				Code:    http.StatusFound,
				Message: fmt.Sprintf("DataSource %s is not found.", id),
			}
		}
		return &entity.DataSource{}, err
	}

	logrus.Debugf("dataSource %s retrieved", id)
	return &data, err
}

func (m *dataSourceRepo) GetAll() ([]*entity.DataSource, error) {
	logrus.Debug("about to get all dataSource")
	var dataSources []*entity.DataSource
	err := m.db.Find(&dataSources).Error
	if err != nil {
		return []*entity.DataSource{}, err
	}
	logrus.Debug("all dataSource retrieved")
	return dataSources, nil
}

func (m *dataSourceRepo) Update(dataSource entity.DataSource) (*entity.DataSource, error) {
	logrus.Debugf("about to update a dataSource %s", dataSource.Name)
	err := m.db.Save(&dataSource).Error
	logrus.Debugf("dataSource %s updated", dataSource.Name)
	return &dataSource, err
}

func (m *dataSourceRepo) Delete(id string) error {
	logrus.Debugf("about to delete a dataSource %s", id)
	tx := m.db.Begin()

	data, err := m.Get(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("data_source_id = ?", id).Delete(&entity.Column{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&entity.DataSource{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	err = m.fileService.Delete(data.FileID)
	if err != nil {
		tx.Rollback()
		return err
	}
	logrus.Debugf("dataSource %s deleted", id)
	return tx.Commit().Error
}
