package entity

type ModelFile struct {
	Base

	Role    string
	FileID  string
	ModelID string
}

func (ModelFile) TableName() string {
	return "ModelFile"
}
