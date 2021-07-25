package entity

type Drawer struct {
	Base

	User string
}

// func (d *Drawer) BeforeCreate(tx *gorm.DB) error {
// 	d.ID = uuid.New().String()
// 	return nil
// }
