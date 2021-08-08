package representation

type Model struct {
	Base

	Description string      `json:"description"`
	Function    string      `json:"function"`
	Files       []ModelFile `json:"files"`
}
