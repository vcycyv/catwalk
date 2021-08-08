package representation

type Server struct {
	Base

	Host   string `json:"host"`
	Port   int    `json:"port"`
	Status string `json:"status"`
	Health bool   `json:"health"`
}
