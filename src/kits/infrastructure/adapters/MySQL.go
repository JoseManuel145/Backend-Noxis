package adapters

import (
	"Backend/src/Kits/domain"
	"Backend/src/core"
	"fmt"
	"log"
	"strconv"
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
	query := "SELECT clave, sensores, username, userfk, status FROM kits"

	rows := mysql.conn.FetchRows(query)
	if rows == nil {
		return nil, fmt.Errorf("error al recuperar los kits")
	}
	defer rows.Close()

	var kits []domain.Kit
	for rows.Next() {
		var kit domain.Kit
		var sensoresStr string
		var userfkRaw []byte

		if err := rows.Scan(&kit.Clave, &sensoresStr, &kit.Username, &userfkRaw, &kit.Status); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %w", err)
		}
		kit.Sensores = strings.Split(sensoresStr, ",")
		if userfkRaw == nil {
			kit.Userfk = 0
		} else {
			userfkStr := string(userfkRaw)
			userfk, err := strconv.Atoi(userfkStr)
			if err != nil {
				return nil, fmt.Errorf("error al convertir userfk a int: %w", err)
			}
			kit.Userfk = userfk
		}
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
