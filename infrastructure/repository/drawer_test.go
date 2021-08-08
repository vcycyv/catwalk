package repository

import (
	"strings"
	"testing"

	"github.com/vcycyv/catwalk/entity"
	"github.com/vcycyv/catwalk/infrastructure/mock"
)

func TestCreate(t *testing.T) {
	db := mock.CreateDB()
	InitDB(db)
	repo := NewDrawerRepo(db)
	name := "test"
	drawer := entity.Drawer{
		Base: entity.Base{Name: name},
		User: "tester",
	}
	newDrawer, err := repo.Add(drawer)
	if err != nil {
		t.Errorf("failed to add a drawer: %v", err)
		return
	}

	if !strings.EqualFold(newDrawer.Name, name) {
		t.Errorf("The added drawer is not correct.")
		return
	}
}
