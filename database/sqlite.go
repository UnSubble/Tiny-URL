package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/unsubble/tiny-url/common"
)

type SqliteRepository struct {
	pool      *sql.DB
	driverCfg *common.DriverConfig
}

func (repo *SqliteRepository) CreateTableIfNotExists(name string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original_url TEXT NOT NULL,
			generated_url TEXT NOT NULL UNIQUE,
			generated_at INT NOT NULL
		);
	`, name)

	_, err := repo.pool.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", name, err)
	}

	return nil
}

func (repo *SqliteRepository) AddUrl(tinyurl *TinyURL) error {
	query := `
		INSERT INTO urls (original_url, short_code)
		VALUES (?, ?, ?);
	`

	_, err := repo.pool.Exec(query, tinyurl.OriginalURL, tinyurl.GeneratedURL, tinyurl.GeneratedAt)
	if err != nil {
		return fmt.Errorf("failed to add URL: %w", err)
	}

	return nil
}

func (repo *SqliteRepository) GetPool() *sql.DB {
	return repo.pool
}

func (repo *SqliteRepository) GetInfo() *common.DriverConfig {
	return repo.driverCfg
}

func NewSqliteRepository(driverCfg *common.DriverConfig) (*SqliteRepository, error) {
	driverCfg.DriverName = "sqlite3"

	db, err := sql.Open("sqlite3", driverCfg.DbName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repo := &SqliteRepository{
		pool:      db,
		driverCfg: driverCfg,
	}

	err = repo.CreateTableIfNotExists("urls")
	if err != nil {
		return nil, err
	}

	return repo, nil
}
