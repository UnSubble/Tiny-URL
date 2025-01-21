package common

type DriverConfig struct {
	DriverName         string `json:"driver"`
	DbName             string `json:"db-name"`
	DbURI              string `json:"db-uri"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	MaxConnections     int    `json:"max-conn"`
	ConnectionIdleTime int    `json:"conn-idle-time"`
}

func NewDriverConfig(driverName string, databaseName string, uri string, username string,
	password string, maxConn int, connIdleTime int) *DriverConfig {
	return &DriverConfig{
		DriverName:         driverName,
		DbName:             databaseName,
		DbURI:              uri,
		Username:           username,
		Password:           password,
		MaxConnections:     maxConn,
		ConnectionIdleTime: connIdleTime,
	}
}
