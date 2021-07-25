package representation

import "time"

type Base struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`

	Links []ResourceLink `json:"links"`
}

type ResourceLink struct {
	Rel    string `json:"rel"`
	Method string `json:"method"`
	Href   string `json:"href"`
	//Uri    string `json:"uri"`
	//Type   string `json:"type"`
}
