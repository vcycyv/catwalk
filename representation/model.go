package representation

type Model struct {
	Base

	Description string      `json:"description"`
	Function    string      `json:"function"`
	Algorithm   string      `json:"algorithm"`
	Files       []ModelFile `json:"files"`
}
