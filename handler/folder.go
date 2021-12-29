package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type folderHandler struct {
	folderService domain.FolderService
}

func NewFolderHandler(folderService domain.FolderService) folderHandler {
	return folderHandler{
		folderService,
	}
}

func (s *folderHandler) Add(c *gin.Context) {
	folder := &entity.Folder{}
	if err := c.ShouldBind(folder); err != nil {
		appErr := &rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		}
		_ = c.Error(appErr)
		return
	}
	logger.Debugf("Received request to add a folder %s.", folder.Path)

	if len(folder.ParentID) > 0 {
		parent, err := s.folderService.GetByID(folder.ParentID)
		if err != nil {
			appErr := &rep.AppError{
				Code:    http.StatusBadRequest,
				Message: "parent folder id is invalid",
			}
			_ = c.Error(appErr)
			return
		}
		folder.Path = parent.Path + "." + folder.Name
	} else {
		folder.Path = folder.Name
	}

	rtnVal, err := s.folderService.Create(*folder)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The folder %s is added successfully.", rtnVal.Name)
	c.JSON(http.StatusCreated, rtnVal)
}

func (s *folderHandler) Get(c *gin.Context) {
	id := c.Param("id")
	folder, err := s.folderService.GetByID(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, folder)
}

func (s *folderHandler) GetAll(c *gin.Context) {
	folders, err := s.folderService.GetAll()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, folders)
}

func (s *folderHandler) Update(c *gin.Context) {
	id := c.Param("id")
	_, err := s.folderService.GetByID(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	name := c.Query("name")

	folder, err := s.folderService.Rename(id, name)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The folder %s is updated successfully.", folder.Name)
	c.JSON(http.StatusOK, folder)
}

func (s *folderHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	logger.Debugf("Received request to delete a drawer %s.", id)
	err := s.folderService.Delete(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The folder %s is deleted successfully.", id)
	c.JSON(http.StatusNoContent, nil)
}
