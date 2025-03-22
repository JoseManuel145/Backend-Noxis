package app

import (
	"Backend/src/Admin/domain"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type LogInUseCase struct {
	repository domain.IAdmin
}

func NewLogInUseCase(repo domain.IAdmin) *LogInUseCase {
	return &LogInUseCase{
		repository: repo,
	}
}
func (useCase *LogInUseCase) Execute(email string, password string) (string, error) {
	SecretKey := os.Getenv("SECRET_KEY")

	// Obtener la contraseña almacenada desde el repositorio
	storedPassword, err := useCase.repository.LogIn(email)
	if err != nil {
		return "", fmt.Errorf("error al obtener la contraseña: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("credenciales incorrectas")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    email, // Se usa el email como issuer
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("error al generar el token: %w", err)
	}

	return token, nil
}
