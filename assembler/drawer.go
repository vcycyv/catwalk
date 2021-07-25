package assembler

import (
	"fmt"
	"net/http"

	"github.com/vcycyv/blog/entity"
	rep "github.com/vcycyv/blog/representation"
)

type DrawerAssembler struct{}

func NewDrawerAssembler() DrawerAssembler {
	return DrawerAssembler{}
}

func (s *DrawerAssembler) ToData(rep rep.Drawer) *entity.Drawer {
	return &entity.Drawer{
		Base: entity.Base{
			ID:        rep.ID,
			CreatedAt: rep.CreatedAt,
			UpdatedAt: rep.UpdatedAt,
			Name:      rep.Name,
		},

		User: rep.User,
	}
}

func (s *DrawerAssembler) ToRepresentation(data entity.Drawer) *rep.Drawer {
	return &rep.Drawer{
		Base: rep.Base{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Name:      data.Name,

			Links: []rep.ResourceLink{
				{
					Rel:    "self",
					Method: http.MethodGet,
					Href:   fmt.Sprintf("/drawers/%s", data.ID),
				},
				{
					Rel:    "add-drawer",
					Method: http.MethodPost,
					Href:   "/drawers",
				},
				{
					Rel:    "edit-drawer",
					Method: http.MethodPut,
					Href:   fmt.Sprintf("/drawers/%s", data.ID),
				},
				{
					Rel:    "delete-drawer",
					Method: http.MethodDelete,
					Href:   fmt.Sprintf("/drawers/%s", data.ID),
				},
			},
		},

		User: data.User,
	}
}
