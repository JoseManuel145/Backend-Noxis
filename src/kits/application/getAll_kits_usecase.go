package application

import (
	"Backend/src/Kits/domain"
	"Backend/src/Kits/domain/repositories"
	"fmt"
)

type GetAllKits struct {
	repo repositories.IKitRepository
}

func NewGetAllKits(repo repositories.IKitRepository) *GetAllKits {
	return &GetAllKits{repo: repo}
}
func (uc *GetAllKits) Execute() ([]domain.Kit, error) {
	kits, err := uc.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo todos los kits: %w", err)
	}
	return kits, nil
}
