package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vcycyv/catwalk/entity"
	rep "github.com/vcycyv/catwalk/representation"
)

func TestDrawer_Add(t *testing.T) {
	var drawerName = "sqlite"
	drawer := addDrawer(drawerName)
	assert.Equal(t, drawerName, drawer.Name)

	_ = db.Migrator().DropTable(&entity.Drawer{})
	_ = db.Migrator().CreateTable(&entity.Drawer{})
}

func TestDrawer_Get(t *testing.T) {
	w := httptest.NewRecorder()

	var drawerName = "sqlite"
	drawer := addDrawer(drawerName)
	drawerName = drawer.Name

	uri := fmt.Sprintf("/drawers/%s", drawer.ID)
	r := gin.Default()
	req := httptest.NewRequest(
		http.MethodGet,
		uri,
		nil,
	)
	r.GET("/drawers/:id", drawerHdlr.Get)
	r.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)
	got := &rep.Drawer{}
	_ = json.Unmarshal(body, &got)
	assert.True(t, len(got.ID) > 0)
	assert.Equal(t, drawerName, got.Name)

	_ = db.Migrator().DropTable(&entity.Drawer{})
	_ = db.Migrator().CreateTable(&entity.Drawer{})
}

func TestDrawer_GetAll(t *testing.T) {
	var drawerName = "sqlite"
	drawer := addDrawer(drawerName)

	drawers := getDrawers()

	assert.Equal(t, len(drawers), 1)
	assert.Equal(t, drawer.Name, drawers[0].Name)

	_ = db.Migrator().DropTable(&entity.Drawer{})
	_ = db.Migrator().CreateTable(&entity.Drawer{})
}

func TestDrawer_Update(t *testing.T) {
	w := httptest.NewRecorder()

	var drawerName = "sqlite"
	var renamed = "sqlite_renamed"
	drawer := addDrawer(drawerName)

	drawer.Name = renamed
	reqBody, _ := json.Marshal(drawer)

	uri := fmt.Sprintf("/drawers/%s", drawer.ID)
	r := gin.Default()
	req := httptest.NewRequest(
		http.MethodPut,
		uri,
		bytes.NewReader(reqBody),
	)
	req.Header.Set("Content-Type", "application/json")

	r.PUT("/drawers/:id", drawerHdlr.Update)
	r.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)
	got := &rep.Drawer{}
	_ = json.Unmarshal(body, &got)
	assert.Equal(t, renamed, got.Name)

	_ = db.Migrator().DropTable(&entity.Drawer{})
	_ = db.Migrator().CreateTable(&entity.Drawer{})
}

func TestDrawer_Delete(t *testing.T) {
	w := httptest.NewRecorder()

	var drawerName = "sqlite"
	drawer := addDrawer(drawerName)

	uri := fmt.Sprintf("/drawers/%s", drawer.ID)
	r := gin.Default()
	req := httptest.NewRequest(
		http.MethodDelete,
		uri,
		nil,
	)

	r.DELETE("/drawers/:id", drawerHdlr.Delete)
	r.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)
	got := &rep.Drawer{}
	_ = json.Unmarshal(body, &got)

	_ = db.Migrator().DropTable(&entity.Drawer{})
	_ = db.Migrator().CreateTable(&entity.Drawer{})
}

func addDrawer(name string) *rep.Drawer {
	w := httptest.NewRecorder()

	uri := "/drawers"
	r := gin.Default()

	drawer := rep.Drawer{
		Base: rep.Base{Name: name},
	}
	reqBody, _ := json.Marshal(drawer)
	req := httptest.NewRequest(
		http.MethodPost,
		uri,
		bytes.NewReader(reqBody),
	)
	req.Header.Set("Content-Type", "application/json")

	r.POST(uri, drawerHdlr.Add)
	//r.GET(uri, drawerController.GetAll)
	r.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)
	_ = json.Unmarshal(body, &drawer)
	return &drawer
}

func getDrawers() []*rep.Drawer {
	w := httptest.NewRecorder()

	uri := "/drawers"
	r := gin.Default()
	req := httptest.NewRequest(
		http.MethodGet,
		uri,
		nil,
	)
	r.GET(uri, drawerHdlr.GetAll)
	r.ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Body)
	var drawers []*rep.Drawer
	_ = json.Unmarshal(body, &drawers)
	return drawers
}
