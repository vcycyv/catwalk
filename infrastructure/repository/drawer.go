package repository

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vcycyv/blog/domain"
	"github.com/vcycyv/blog/entity"
	"github.com/vcycyv/blog/representation"
	"gorm.io/gorm"
)

type drawerRepo struct {
	db *gorm.DB
}

func NewDrawerRepo(db *gorm.DB) domain.DrawerRepository {
	return &drawerRepo{
		db,
	}
}

func (m *drawerRepo) Add(drawer entity.Drawer) (*entity.Drawer, error) {
	logrus.Debugf("about to save a drawer %s", drawer.Name)
	if err := m.db.Create(&drawer).Error; err != nil {
		return nil, err
	}
	logrus.Debugf("drawer %s saved", drawer.Name)
	return &drawer, nil
}

func (m *drawerRepo) Get(id string) (*entity.Drawer, error) {
	logrus.Debugf("about to get a drawer %s", id)
	var data entity.Drawer
	err := m.db.Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &representation.AppError{
				Code:    http.StatusFound,
				Message: fmt.Sprintf("Drawer %s is not found.", id),
			}
		}
		return &entity.Drawer{}, err
	}

	logrus.Debugf("drawer %s retrieved", id)
	return &data, err
}

func (m *drawerRepo) GetAll() ([]*entity.Drawer, error) {
	logrus.Debug("about to get all drawer")
	var drawers []*entity.Drawer
	err := m.db.Find(&drawers).Error
	if err != nil {
		return []*entity.Drawer{}, err
	}
	logrus.Debug("all drawer retrieved")
	return drawers, nil
}

func (m *drawerRepo) Update(drawer entity.Drawer) (*entity.Drawer, error) {
	logrus.Debugf("about to update a drawer %s", drawer.Name)
	err := m.db.Save(&drawer).Error
	logrus.Debugf("drawer %s updated", drawer.Name)
	return &drawer, err
}

func (m *drawerRepo) Delete(id string) error {
	logrus.Debugf("about to delete a drawer %s", id)
	tx := m.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.Drawer{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	logrus.Debugf("drawer %s deleted", id)
	return tx.Commit().Error
}
