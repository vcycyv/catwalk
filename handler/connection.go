package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type connectionHandler struct {
	connectionService domain.ConnectionInterface
}

func NewConnectionHandler(connectionService domain.ConnectionInterface) *connectionHandler {
	return &connectionHandler{
		connectionService,
	}
}

func (s *connectionHandler) Add(c *gin.Context) {
	connection := &rep.Connection{}
	if err := c.ShouldBind(connection); err != nil {
		appErr := &rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		}
		_ = c.Error(appErr)
		return
	}
	logger.Debugf("Received request to add a connection %s.", connection.Name)

	rtnVal, err := s.connectionService.Add(*connection)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The connection %s is added successfully.", rtnVal.Name)
	c.JSON(http.StatusCreated, rtnVal)
}

func (s *connectionHandler) Get(c *gin.Context) {
	id := c.Param("id")
	connection, err := s.connectionService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, connection)
}

func (s *connectionHandler) GetAll(c *gin.Context) {
	connections, err := s.connectionService.GetAll()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, connections)
}

func (s *connectionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	_, err := s.connectionService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	connection := &rep.Connection{}
	if err := c.ShouldBind(connection); err != nil {
		_ = c.Error(&rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		})
		return
	}
	logger.Debugf("Received request to add a connection %s", connection.Name)

	connection, err = s.connectionService.Update(*connection)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The connection %s is updated successfully.", connection.Name)
	c.JSON(http.StatusOK, connection)
}

func (s *connectionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	logger.Debugf("Received request to delete a connection %s.", id)
	err := s.connectionService.Delete(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The connection %s is deleted successfully.", id)
	c.JSON(http.StatusNoContent, nil)
}
