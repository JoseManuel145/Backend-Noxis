package adapters

import (
	"Backend/src/Admin/domain"
	"Backend/src/core"
	"fmt"
	"log"
)

type MySQL struct {
	conn *core.Conn_MySQL
}

var _ domain.IAdmin = (*MySQL)(nil)

func NewMySQL() *MySQL {
	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}
func (mysql *MySQL) Register(email, password string) error {
	query := "INSERT INTO admins (email, password) VALUES (?, ?)"

	_, err := mysql.conn.ExecutePreparedQuery(query, email, password)
	if err != nil {
		return fmt.Errorf("error al guardar el admin: %w", err)
	}
	return nil
}
func (mysql *MySQL) LogIn(email string) (string, error) {
	query := "SELECT password FROM admins WHERE email = ?"

	rows := mysql.conn.FetchRows(query, email)
	defer rows.Close()

	var storedPassword string
	if rows.Next() {
		if err := rows.Scan(&storedPassword); err != nil {
			return "", fmt.Errorf("error al escanear la contraseña: %w", err)
		}
	} else {
		return "", fmt.Errorf("no se encontró ningún admin con el email: %s", email)
	}

	return storedPassword, nil
}
