package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	//"github.com/voznikaetnepriyazn/Autorization_service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func appendDSNParam(dsn, key, value string) string {
	sep := "?"
	if strings.ContainsRune(dsn, '?') {
		sep = "&"
	}
	return dsn + sep + key + "=" + value
}

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "PostgreSQL DSN")
	flag.StringVar(&migrationsPath, "migration-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migration")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required(PostgreSQl DSN)")
	}

	if migrationsPath == "" {
		panic("migrations-path is requred")
	}

	dsnWithTable := appendDSNParam(storagePath, "x-migrations-table", migrationsTable)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsnWithTable,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied succesfully")
}
