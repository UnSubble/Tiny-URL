package common

type DriverConfig struct {
	DriverName            string
	DbName                string `json:"db-uri"`
	MaxConnections        int    `json:"max-conn"`
	MaxConnectionIdleTime int    `json:"conn-idle-time"`
}

func NewDriverConfig(driverName string, databaseName string, maxConn int, connIdleTime int) *DriverConfig {
	return &DriverConfig{
		DriverName:            driverName,
		DbName:                databaseName,
		MaxConnections:        maxConn,
		MaxConnectionIdleTime: connIdleTime,
	}
}
