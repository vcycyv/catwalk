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

type modelRepo struct {
	db *gorm.DB
}

func NewModelRepo(db *gorm.DB) domain.ModelRepository {
	return &modelRepo{
		db,
	}
}

func (m *modelRepo) Add(model entity.Model) (*entity.Model, error) {
	logrus.Debugf("about to save a model %s", model.Name)
	if err := m.db.Create(&model).Error; err != nil {
		return nil, err
	}
	logrus.Debugf("model %s saved", model.Name)
	return &model, nil
}

func (m *modelRepo) Get(id string) (*entity.Model, error) {
	logrus.Debugf("about to get a model %s", id)
	var data entity.Model
	err := m.db.Where("id = ?", id).Preload("Files").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &representation.AppError{
				Code:    http.StatusFound,
				Message: fmt.Sprintf("Model %s is not found.", id),
			}
		}
		return &entity.Model{}, err
	}

	logrus.Debugf("model %s retrieved", id)
	return &data, err
}

func (m *modelRepo) GetAll() ([]*entity.Model, error) {
	logrus.Debug("about to get all model")
	var models []*entity.Model
	err := m.db.Find(&models).Error
	if err != nil {
		return []*entity.Model{}, err
	}
	logrus.Debug("all model retrieved")
	return models, nil
}

func (m *modelRepo) Update(model entity.Model) (*entity.Model, error) {
	logrus.Debugf("about to update a model %s", model.Name)
	err := m.db.Save(&model).Error
	logrus.Debugf("model %s updated", model.Name)
	return &model, err
}

func (m *modelRepo) Delete(id string) error {
	logrus.Debugf("about to delete a model %s", id)
	tx := m.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.Model{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	logrus.Debugf("model %s deleted", id)
	return tx.Commit().Error
}
