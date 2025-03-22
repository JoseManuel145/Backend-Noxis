package app

import (
	"Backend/src/Admin/domain"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUseCase struct {
	repository domain.IAdmin
}

func NewRegisterUseCase(repo domain.IAdmin) *RegisterUseCase {
	return &RegisterUseCase{
		repository: repo,
	}
}
func (useCase *RegisterUseCase) Execute(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al cifrar la contraseña: %w", err)
	}

	err = useCase.repository.Register(email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("error al registrar el admin: %w", err)
	}

	return nil
}
