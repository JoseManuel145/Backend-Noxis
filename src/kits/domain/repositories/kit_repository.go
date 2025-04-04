package repositories

import "Backend/src/Kits/domain"

type IKitRepository interface {
	GetAll() ([]domain.Kit, error)
	CreateKit(clave string, sensores []string, username string) error
	UpdateKit(clave string, status bool, userfk int) error
	GetInactives() ([]domain.Kit, error)
	GetActives() ([]domain.Kit, error)
}
