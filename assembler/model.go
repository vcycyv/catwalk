package assembler

import (
	"fmt"
	"net/http"

	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

type ModelAssembler struct{}

func NewModelAssembler() ModelAssembler {
	return ModelAssembler{}
}

func (s *ModelAssembler) ToData(rep rep.Model) *entity.Model {
	var modelFiles []entity.ModelFile
	for _, modelFile := range rep.Files {
		modelFiles = append(modelFiles, *ModelFileAss.ToData(modelFile))
	}
	return &entity.Model{
		Base: entity.Base{
			ID:        rep.ID,
			CreatedAt: rep.CreatedAt,
			UpdatedAt: rep.UpdatedAt,
			Name:      rep.Name,
		},

		Description: rep.Description,
		Function:    rep.Function,
		Files:       modelFiles,
	}
}

func (s *ModelAssembler) ToRepresentation(data entity.Model) *rep.Model {
	var modelFiles []rep.ModelFile
	for _, modelFile := range data.Files {
		modelFiles = append(modelFiles, *ModelFileAss.ToRepresentation(modelFile))
	}
	return &rep.Model{
		Base: rep.Base{
			ID:        data.ID,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Name:      data.Name,

			Links: []rep.ResourceLink{
				{
					Rel:    "self",
					Method: http.MethodGet,
					Href:   fmt.Sprintf("/models/%s", data.ID),
				},
				{
					Rel:    "add-model",
					Method: http.MethodPost,
					Href:   "/models",
				},
				{
					Rel:    "edit-model",
					Method: http.MethodPut,
					Href:   fmt.Sprintf("/models/%s", data.ID),
				},
				{
					Rel:    "delete-model",
					Method: http.MethodDelete,
					Href:   fmt.Sprintf("/models/%s", data.ID),
				},
			},
		},

		Description: data.Description,
		Function:    data.Function,
		Files:       modelFiles,
	}
}
