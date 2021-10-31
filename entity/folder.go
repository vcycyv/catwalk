package entity

type Folder struct {
	ID       string
	ParentID string `db:"-" json:",omitempty"`
	Path     string
	Name     string `db:"-" json:",omitempty"`
}
