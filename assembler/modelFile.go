package assembler

import (
	"fmt"
	"net/http"

	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type ModelFileAssembler struct{}

func NewModelFileAssembler() ModelFileAssembler {
	return ModelFileAssembler{}
}

func (s *ModelFileAssembler) ToData(rep rep.ModelFile) *entity.ModelFile {
	return &entity.ModelFile{
		Base: entity.Base{
			ID:        rep.ID,
			CreatedAt: rep.CreatedAt,
			UpdatedAt: rep.UpdatedAt,
			Name:      rep.Name,
		},

		Role:    rep.Role,
		FileID:  rep.FileID,
		ModelID: rep.ModelID,
	}
}

func (s *ModelFileAssembler) ToRepresentation(data entity.ModelFile) *rep.ModelFile {
	return &rep.ModelFile{
		Base: rep.Base{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Name:      data.Name,

			Links: []rep.ResourceLink{
				{
					Rel:    "self",
					Method: http.MethodGet,
					Href:   fmt.Sprintf("/modelFiles/%s", data.ID),
				},
				{
					Rel:    "add-modelFile",
					Method: http.MethodPost,
					Href:   "/modelFiles",
				},
				{
					Rel:    "edit-modelFile",
					Method: http.MethodPut,
					Href:   fmt.Sprintf("/modelFiles/%s", data.ID),
				},
				{
					Rel:    "delete-modelFile",
					Method: http.MethodDelete,
					Href:   fmt.Sprintf("/modelFiles/%s", data.ID),
				},
			},
		},

		Role:    data.Role,
		FileID:  data.FileID,
		ModelID: data.ModelID,
	}
}
