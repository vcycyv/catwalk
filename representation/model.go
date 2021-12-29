package representation

type Model struct {
	Base

	FolderID    string      `json:"folderId"`
	Description string      `json:"description"`
	Function    string      `json:"function"`
	Algorithm   string      `json:"algorithm"`
	Target      string      `json:"target"`
	Variables   []string    `json:"variables"`
	Files       []ModelFile `json:"files"`
}
