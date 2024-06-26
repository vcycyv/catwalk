package entity

type Model struct {
	Base

	FolderID    string
	Description string
	Function    string
	Algorithm   string
	Files       []ModelFile `gorm:"foreignkey:ModelID;AssociationForeignKey:id"`
}

func (Model) TableName() string {
	return "Model"
}
