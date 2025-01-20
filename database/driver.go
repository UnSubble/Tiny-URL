package database

import (
	"database/sql"

	"github.com/unsubble/tiny-url/common"
)

type TinyURL common.TinyURL

type Driver interface {
	CreateTableIfNotExists(string) error
	AddUrl(*TinyURL) error
	GetPool() *sql.DB
	GetInfo() *common.DriverConfig
}
