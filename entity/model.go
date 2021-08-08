package entity

type Model struct {
	Base

	Description string
	Function    string
	Files       []ModelFile `gorm:"foreignkey:ModelID;AssociationForeignKey:id"`
}

func (Model) TableName() string {
	return "Model"
}
