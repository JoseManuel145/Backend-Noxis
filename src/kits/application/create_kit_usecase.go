package application

import (
	"Backend/src/kits/domain"
	"Backend/src/kits/domain/repositories"
	"fmt"
)

type CreateKit struct {
	repo repositories.IKitRepository
}

func NewCreateKit(repo repositories.IKitRepository) *CreateKit {
	return &CreateKit{repo: repo}
}
func (uc *CreateKit) Execute(kit *domain.Kit) error {
	clave := GenerateClave(5)
	kit.Clave = clave
	err := uc.repo.CreateKit(kit.Clave, kit.Sensores, kit.Username)
	if err != nil {
		return fmt.Errorf("error creando el kit: %w", err)
	}
	return nil
}
