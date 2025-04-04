package adapters

import (
	"Backend/src/Kits/domain"
	"Backend/src/core"
	"fmt"
	"log"
	"strings"
)

type MySQL struct {
	conn *core.Conn_MySQL
}

func NewMySQL() *MySQL {
	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}
func (mysql *MySQL) CreateKit(clave string, sensores []string, username string) error {
	query := "INSERT INTO kits (clave, sensores, username, status) VALUES (?, ?, ?, ?)"
	sensoresStr := strings.Join(sensores, ",")

	_, err := mysql.conn.ExecutePreparedQuery(query, clave, sensoresStr, username, false)
	if err != nil {
		return fmt.Errorf("error al guardar el kit: %w", err)
	}
	return nil
}
func (mysql *MySQL) GetAll() ([]domain.Kit, error) {
	query := "SELECT clave, sensores, username, status FROM kits"

	rows := mysql.conn.FetchRows(query)
	if rows == nil {
		return nil, fmt.Errorf("error al recuperar los kits")
	}
	defer rows.Close()

	var kits []domain.Kit
	for rows.Next() {
		var kit domain.Kit
		var sensoresStr string

		if err := rows.Scan(&kit.Clave, &sensoresStr, &kit.Username, &kit.Status); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %w", err)
		}
		kit.Sensores = strings.Split(sensoresStr, ",")

		kits = append(kits, kit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando sobre las filas: %w", err)
	}
	return kits, nil
}
func (mysql *MySQL) UpdateKit(clave string, status bool, userFk int) error {
	query := "UPDATE kits SET status = ?, userfk = ? WHERE clave = ?"
	_, err := mysql.conn.ExecutePreparedQuery(query, status, userFk, clave)
	if err != nil {
		return fmt.Errorf("error al actualizar el kit: %w", err)
	}
	return nil
}
func (mysql *MySQL) GetInactives() ([]domain.Kit, error) {
	query := "SELECT clave, sensores, username, status FROM kits WHERE status = 0"

	rows := mysql.conn.FetchRows(query)
	if rows == nil {
		return nil, fmt.Errorf("error al recuperar los kits inactivos")
	}
	defer rows.Close()

	var kits []domain.Kit
	for rows.Next() {
		var kit domain.Kit
		var sensoresStr string

		if err := rows.Scan(&kit.Clave, &sensoresStr, &kit.Username, &kit.Status); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %w", err)
		}

		kit.Sensores = strings.Split(sensoresStr, ",")

		kits = append(kits, kit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando sobre las filas: %w", err)
	}
	return kits, nil
}
func (mysql *MySQL) GetActives() ([]domain.Kit, error) {
	query := "SELECT clave, sensores, username, status FROM kits WHERE status = 1"

	rows := mysql.conn.FetchRows(query)
	if rows == nil {
		return nil, fmt.Errorf("error al recuperar los kits activos")
	}
	defer rows.Close()

	var kits []domain.Kit
	for rows.Next() {
		var kit domain.Kit
		var sensoresStr string

		if err := rows.Scan(&kit.Clave, &sensoresStr, &kit.Username, &kit.Status); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %w", err)
		}

		kit.Sensores = strings.Split(sensoresStr, ",")

		kits = append(kits, kit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando sobre las filas: %w", err)
	}
	return kits, nil
}
