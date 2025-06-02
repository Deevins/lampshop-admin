package infra

import (
	"database/sql"
)

// ConnectDB пытается подключиться к базе по строке connString.
// В дальнейшем вместо in-memory репозитория можно вернуть реальную БД.
func ConnectDB(connString string) (*sql.DB, error) {
	return nil, nil
}
