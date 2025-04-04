package application

import (
	"Backend/src/kits/domain"
	"Backend/src/kits/domain/repositories"
	"fmt"
)

type GetAllActives struct {
	repo repositories.IKitRepository
}

func NewGetAllActives(repo repositories.IKitRepository) *GetAllActives {
	return &GetAllActives{repo: repo}
}
func (uc *GetAllActives) Execute() ([]domain.Kit, error) {
	kits, err := uc.repo.GetActives()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo todos los kits: %w", err)
	}
	return kits, nil
}
