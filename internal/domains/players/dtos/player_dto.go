package dtos

type PlayerRequest struct {
	TeamID       uint   `json:"teamId" validate:"required,gt=0"`
	Name         string `json:"name" validate:"required,min=2,max=120"`
	HeightCM     int    `json:"heightCm" validate:"required,gte=120,lte=250"`
	WeightKG     int    `json:"weightKg" validate:"required,gte=35,lte=180"`
	Position     string `json:"position" validate:"required,oneof=penyerang gelandang bertahan penjaga_gawang"`
	JerseyNumber int    `json:"jerseyNumber" validate:"required,gte=1,lte=99"`
}
