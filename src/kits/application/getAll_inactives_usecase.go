package application

import (
	"Backend/src/Kits/domain"
	"Backend/src/Kits/domain/repositories"
	"fmt"
)

type GetAllInactives struct {
	repo repositories.IKitRepository
}

func NewGetAllInactives(repo repositories.IKitRepository) *GetAllInactives {
	return &GetAllInactives{repo: repo}
}
func (uc *GetAllInactives) Execute() ([]domain.Kit, error) {
	kits, err := uc.repo.GetInactives()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo todos los kits inactivos: %w", err)
	}
	return kits, nil
}
