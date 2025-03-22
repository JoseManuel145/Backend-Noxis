package domain

type IAdmin interface {
	Register(email, password string) error
	LogIn(email string) (string, error)
}
