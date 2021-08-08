package assembler

import (
	"fmt"
	"net/http"

	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type ServerAssembler struct{}

func NewServerAssembler() ServerAssembler {
	return ServerAssembler{}
}

func (s *ServerAssembler) ToData(rep rep.Server) *entity.Server {
	return &entity.Server{
		Base: entity.Base{
			ID:        rep.ID,
			CreatedAt: rep.CreatedAt,
			UpdatedAt: rep.UpdatedAt,
			Name:      rep.Name,
		},

		Host:   rep.Host,
		Port:   rep.Port,
		Status: rep.Status,
		Health: rep.Health,
	}
}

func (s *ServerAssembler) ToRepresentation(data entity.Server) *rep.Server {
	return &rep.Server{
		Base: rep.Base{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Name:      data.Name,

			Links: []rep.ResourceLink{
				{
					Rel:    "self",
					Method: http.MethodGet,
					Href:   fmt.Sprintf("/servers/%s", data.ID),
				},
				{
					Rel:    "add-server",
					Method: http.MethodPost,
					Href:   "/servers",
				},
				{
					Rel:    "edit-server",
					Method: http.MethodPut,
					Href:   fmt.Sprintf("/servers/%s", data.ID),
				},
				{
					Rel:    "delete-server",
					Method: http.MethodDelete,
					Href:   fmt.Sprintf("/servers/%s", data.ID),
				},
			},
		},

		Host:   data.Host,
		Port:   data.Port,
		Status: data.Status,
		Health: data.Health,
	}
}
