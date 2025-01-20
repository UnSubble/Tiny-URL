package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/unsubble/tiny-url/common"
)

type PostgresRepository struct {
	pool      *sql.DB
	driverCfg *common.DriverConfig
}

func (repo *PostgresRepository) CreateTableIfNotExists(name string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
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

func (repo *PostgresRepository) AddUrl(tinyurl *TinyURL) error {
	query := `
		INSERT INTO urls (original_url, generated_url, generated_at)
		VALUES ($1, $2, $3);
	`

	_, err := repo.pool.Exec(query, tinyurl.OriginalURL, tinyurl.GeneratedURL, tinyurl.GeneratedAt.UnixMilli())

	if err != nil {
		return fmt.Errorf("failed to add URL: %w", err)
	}

	return nil
}

func (repo *PostgresRepository) GetPool() *sql.DB {
	return repo.pool
}

func (repo *PostgresRepository) GetInfo() *common.DriverConfig {
	return repo.driverCfg
}

func NewPostgresRepository(driverCfg *common.DriverConfig) (*PostgresRepository, error) {
	driverCfg.DriverName = "postgres"
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s",
		driverCfg.Username, driverCfg.Password, driverCfg.DbURI, driverCfg.DbName)

	db, err := sql.Open(driverCfg.DriverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	db.SetMaxOpenConns(driverCfg.MaxConnections)
	db.SetConnMaxIdleTime(time.Duration(driverCfg.ConnectionIdleTime))

	repo := &PostgresRepository{
		pool:      db,
		driverCfg: driverCfg,
	}

	err = repo.CreateTableIfNotExists("urls")
	if err != nil {
		return nil, err
	}

	return repo, nil
}
