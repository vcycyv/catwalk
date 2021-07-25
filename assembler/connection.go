package assembler

import (
	"fmt"
	"net/http"

	"github.com/vcycyv/blog/entity"
	rep "github.com/vcycyv/blog/representation"
)

type ConnectionAssembler struct{}

func NewConnectionAssembler() ConnectionAssembler {
	return ConnectionAssembler{}
}

func (s *ConnectionAssembler) ToData(rep rep.Connection) *entity.Connection {
	return &entity.Connection{
		Base: entity.Base{
			ID:        rep.ID,
			CreatedAt: rep.CreatedAt,
			UpdatedAt: rep.UpdatedAt,
			Name:      rep.Name,
		},

		Type:     rep.Type,
		Host:     rep.Host,
		User:     rep.User,
		Password: rep.Password,
		DbName:   rep.DbName,
	}
}

func (s *ConnectionAssembler) ToRepresentation(data entity.Connection) *rep.Connection {
	return &rep.Connection{
		Base: rep.Base{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Name:      data.Name,

			Links: []rep.ResourceLink{
				{
					Rel:    "self",
					Method: http.MethodGet,
					Href:   fmt.Sprintf("/connections/%s", data.ID),
				},
				{
					Rel:    "add-connection",
					Method: http.MethodPost,
					Href:   "/connections",
				},
				{
					Rel:    "edit-connection",
					Method: http.MethodPut,
					Href:   fmt.Sprintf("/connections/%s", data.ID),
				},
				{
					Rel:    "delete-connection",
					Method: http.MethodDelete,
					Href:   fmt.Sprintf("/connections/%s", data.ID),
				},
			},
		},

		Type:     data.Type,
		Host:     data.Host,
		User:     data.User,
		Password: data.Password,
		DbName:   data.DbName,
	}
}
