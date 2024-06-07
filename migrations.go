package navitas

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (n *Navitas) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+n.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}
	return nil
}

func (n *Navitas) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+n.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}
	return nil
}

func (n *Navitas) Steps(val int, dsn string) error {
	m, err := migrate.New("file://"+n.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(val); err != nil {
		return err
	}
	return nil
}

func (n *Navitas) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+n.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		return err
	}
	return nil
}
