package entity

type Drawer struct {
	Base

	User string
}

func (Drawer) TableName() string {
	return "Drawer"
}
