package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type drawerHandler struct {
	drawerService domain.DrawerInterface
	authService   domain.AuthInterface
}

func NewDrawerHandler(drawerService domain.DrawerInterface, authService domain.AuthInterface) drawerHandler {
	return drawerHandler{
		drawerService,
		authService,
	}
}

func (s *drawerHandler) Add(c *gin.Context) {
	drawer := &rep.Drawer{}
	if err := c.ShouldBind(drawer); err != nil {
		appErr := &rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		}
		_ = c.Error(appErr)
		return
	}
	logger.Debugf("Received request to add a drawer %s.", drawer.Name)

	token := s.authService.ExtractToken(c)
	drawer.User, _ = s.authService.GetUserFromToken(token)

	rtnVal, err := s.drawerService.Add(*drawer)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The drawer %s is added successfully.", rtnVal.Name)
	c.JSON(http.StatusCreated, rtnVal)
}

func (s *drawerHandler) Get(c *gin.Context) {
	id := c.Param("id")
	drawer, err := s.drawerService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, drawer)
}

func (s *drawerHandler) GetAll(c *gin.Context) {
	drawers, err := s.drawerService.GetAll()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, drawers)
}

func (s *drawerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	_, err := s.drawerService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	drawer := &rep.Drawer{}
	if err := c.ShouldBind(drawer); err != nil {
		_ = c.Error(&rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		})
		return
	}
	logger.Debugf("Received request to add a drawer %s", drawer.Name)

	drawer, err = s.drawerService.Update(*drawer)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The drawer %s is updated successfully.", drawer.Name)
	c.JSON(http.StatusOK, drawer)
}

func (s *drawerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	logger.Debugf("Received request to delete a drawer %s.", id)
	err := s.drawerService.Delete(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The drawer %s is deleted successfully.", id)
	c.JSON(http.StatusNoContent, nil)
}
