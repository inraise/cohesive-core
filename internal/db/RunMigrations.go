package db

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migrations
var migrationFiles embed.FS

func runMigrations(connStr string) error {
	driver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("Ошибка инициализации: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", driver, connStr)
	if err != nil {
		return fmt.Errorf("Ошибка создания мигратора: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Схема базы данных актуальна. Новых миграций нет.")
			return nil
		}
		return fmt.Errorf("Ошибка применения миграций: %w", err)
	}

	fmt.Println("Миграции успешно применены!")
	return nil
}
