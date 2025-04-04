package application

import (
	"Backend/src/Kits/domain/repositories"
	"fmt"
)

type UpdateKit struct {
	repo repositories.IKitRepository
}

func NewUpdateKit(repo repositories.IKitRepository) *UpdateKit {
	return &UpdateKit{repo: repo}
}
func (uc *UpdateKit) Execute(clave string, status bool, userfk int) error {
	err := uc.repo.UpdateKit(clave, status, userfk)
	if err != nil {
		return fmt.Errorf("error actualizando kit: %w", err)
	}
	return nil
}
