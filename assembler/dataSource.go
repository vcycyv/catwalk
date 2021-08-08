package assembler

import (
	"fmt"
	"net/http"

	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type DataSourceAssembler struct{}

func NewDataSourceAssebler() DataSourceAssembler {
	return DataSourceAssembler{}
}

func (s *DataSourceAssembler) ToData(rep rep.DataSource) *entity.DataSource {
	return &entity.DataSource{
		Base: entity.Base{
			ID:        rep.ID,
			CreatedAt: rep.CreatedAt,
			UpdatedAt: rep.UpdatedAt,
			Name:      rep.Name,
		},

		DrawerID:    rep.DrawerID,
		Description: rep.Description,
		User:        rep.User,
		FileID:      rep.FileID,
	}
}

func (s *DataSourceAssembler) ToRepresentation(data entity.DataSource) *rep.DataSource {
	var columns []string
	for _, columnData := range data.Columns {
		columns = append(columns, columnData.Name)
	}
	return &rep.DataSource{
		Base: rep.Base{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Name:      data.Name,

			Links: []rep.ResourceLink{
				{
					Rel:    "self",
					Method: http.MethodGet,
					Href:   fmt.Sprintf("/dataSource/%s", data.ID),
				},
				{
					Rel:    "add-dataSource",
					Method: http.MethodPost,
					Href:   "/dataSource",
				},
				{
					Rel:    "edit-dataSource",
					Method: http.MethodPut,
					Href:   fmt.Sprintf("/dataSource/%s", data.ID),
				},
				{
					Rel:    "delete-dataSource",
					Method: http.MethodDelete,
					Href:   fmt.Sprintf("/dataSource/%s", data.ID),
				},
			},
		},

		DrawerID:    data.DrawerID,
		Description: data.Description,
		User:        data.User,
		FileID:      data.FileID,
		Columns:     columns,
	}
}
