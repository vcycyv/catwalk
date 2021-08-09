package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
	rep "github.com/vcycyv/catwalk/representation"
)

type serverHandler struct {
	serverService domain.ServerInterface
}

func NewServerHandler(serverService domain.ServerInterface) serverHandler {
	return serverHandler{
		serverService,
	}
}

func (s *serverHandler) Add(c *gin.Context) {
	server := &rep.Server{}
	if err := c.ShouldBind(server); err != nil {
		appErr := &rep.AppError{
			Code:    http.StatusBadRequest,
			Message: Message.InvalidMessage,
		}
		_ = c.Error(appErr)
		return
	}
	logger.Debugf("Received request to add a server %s.", server.Name)

	rtnVal, err := s.serverService.Add(*server)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The server %s is added successfully.", rtnVal.Name)
	c.JSON(http.StatusCreated, rtnVal)
}

func (s *serverHandler) Get(c *gin.Context) {
	id := c.Param("id")
	server, err := s.serverService.Get(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, server)
}

func (s *serverHandler) GetAll(c *gin.Context) {
	servers, err := s.serverService.GetAll(make(map[string]interface{}))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, servers)
}

func (s *serverHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	logger.Debugf("Received request to delete a server %s.", id)
	err := s.serverService.Delete(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.Debugf("The server %s is deleted successfully.", id)
	c.JSON(http.StatusNoContent, nil)
}
