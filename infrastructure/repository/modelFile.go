package repository

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/entity"
	"github.com/vcycyv/catwalk/representation"
	"gorm.io/gorm"
)

type modelFileRepo struct {
	db          *gorm.DB
	fileService domain.FileService
}

func NewModelFileRepo(db *gorm.DB, fileService domain.FileService) domain.ModelFileRepository {
	return &modelFileRepo{
		db,
		fileService,
	}
}

func (m *modelFileRepo) Add(modelFile entity.ModelFile, reader io.Reader) (*entity.ModelFile, error) {
	logrus.Debugf("about to save a modelFile %s", modelFile.Name)
	fID, err := m.fileService.Save(modelFile.Name, reader)
	if err != nil {
		return nil, err
	}
	modelFile.FileID = fID
	if err := m.db.Create(&modelFile).Error; err != nil {
		return nil, err
	}
	logrus.Debugf("modelFile %s saved", modelFile.Name)
	return &modelFile, nil
}

func (m *modelFileRepo) Get(id string) (*entity.ModelFile, error) {
	logrus.Debugf("about to get a modelFile %s", id)
	var data entity.ModelFile
	err := m.db.Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &representation.AppError{
				Code:    http.StatusFound,
				Message: fmt.Sprintf("ModelFile %s is not found.", id),
			}
		}
		return &entity.ModelFile{}, err
	}

	logrus.Debugf("modelFile %s retrieved", id)
	return &data, err
}

func (m *modelFileRepo) GetContent(id string, writer io.Writer) error {
	modelFile, err := m.Get(id)
	if err != nil {
		return err
	}
	err = m.fileService.DirectContentToWriter(modelFile.FileID, writer)
	if err != nil {
		return err
	}

	return nil
}

func (m *modelFileRepo) GetAll() ([]*entity.ModelFile, error) {
	logrus.Debug("about to get all modelFile")
	var modelFiles []*entity.ModelFile
	err := m.db.Find(&modelFiles).Error
	if err != nil {
		return []*entity.ModelFile{}, err
	}
	logrus.Debug("all modelFile retrieved")
	return modelFiles, nil
}

func (m *modelFileRepo) Update(modelFile entity.ModelFile) (*entity.ModelFile, error) {
	logrus.Debugf("about to update a modelFile %s", modelFile.Name)
	err := m.db.Save(&modelFile).Error
	logrus.Debugf("modelFile %s updated", modelFile.Name)
	return &modelFile, err
}

func (m *modelFileRepo) Delete(id string) error {
	logrus.Debugf("about to delete a modelFile %s", id)
	tx := m.db.Begin()
	if err := tx.Where("id = ?", id).Delete(&entity.ModelFile{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	logrus.Debugf("modelFile %s deleted", id)
	return tx.Commit().Error
}
