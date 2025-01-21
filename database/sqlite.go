package database

import (
	"database/sql"
	"fmt"
	"time"

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
			generated_at BIGINT NOT NULL
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
		INSERT INTO urls (original_url, generated_url, generated_at)
		VALUES (?, ?, ?);
	`

	_, err := repo.pool.Exec(query, tinyurl.OriginalURL, tinyurl.GeneratedURL, tinyurl.GeneratedAt.UnixMilli())
	if err != nil {
		return fmt.Errorf("failed to add URL: %w", err)
	}

	return nil
}

func (repo *SqliteRepository) GetOriginalUrl(tiny string) (string, error) {
	query := `
		SELECT original_url FROM urls
		WHERE generated_url = ?;
	`

	var originalUrl string
	err := repo.pool.QueryRow(query, tiny).Scan(&originalUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("there is no such URL: %s", tiny)
		}
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}

	return originalUrl, nil
}

func (repo *SqliteRepository) GetPool() *sql.DB {
	return repo.pool
}

func (repo *SqliteRepository) GetInfo() *common.DriverConfig {
	return repo.driverCfg
}

func NewSqliteRepository(driverCfg *common.DriverConfig) (*SqliteRepository, error) {
	driverCfg.DriverName = "sqlite3"
	connStr := driverCfg.DbURI

	if driverCfg.DbURI == "" {
		return nil, fmt.Errorf("database URI is required")
	}

	db, err := sql.Open(driverCfg.DriverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	db.SetMaxOpenConns(driverCfg.MaxConnections)
	db.SetConnMaxIdleTime(time.Duration(driverCfg.ConnectionIdleTime))

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
