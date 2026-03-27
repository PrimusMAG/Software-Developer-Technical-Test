package dtos

type TeamRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=120"`
	FoundedYear int    `json:"foundedYear" validate:"required,gte=1850,lte=2100"`
	HQAddress   string `json:"hqAddress" validate:"required,max=255"`
	HQCity      string `json:"hqCity" validate:"required,max=120"`
}
