package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vcycyv/blog/entity"
	"github.com/vcycyv/blog/representation"
)

func TestDrawer_Add(t *testing.T) {
	name := "sqlite"
	drawer := representation.Drawer{
		Base: representation.Base{
			Name: name,
		},
		User: "user_a",
	}
	newDrawer, _ := drawerSvc.Add(drawer)
	existingDrawer, _ := drawerSvc.Get(newDrawer.ID)
	assert.Equal(t, name, existingDrawer.Name)

	_ = db.Migrator().DropTable(&entity.Drawer{})
	_ = db.Migrator().CreateTable(&entity.Drawer{})
}
