package navitas

import (
	"database/sql"
)

// OpenDB opens a connection to a sql database. dbType must be one of postgres (or pgx).
// TODO: add support for mysql/mariadb
func (n *Navitas) OpenDB(dbType, dsn string) (*sql.DB, error) {
	if dbType == "postgres" || dbType == "postgresql" {
		dbType = "pgx"
	}

	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
