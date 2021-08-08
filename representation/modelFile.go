package representation

type ModelFile struct {
	Base

	Role    string `json:"role"`
	FileID  string `json:"fileId"`
	ModelID string `json:"modelId"`
}
