package common

type Config struct {
	Domain       string        `json:"domain"`
	BaseFormat   string        `json:"base-format"`
	DriverConfig *DriverConfig `json:"database"`
	EpochTime    uint64        `json:"epoch-time"`
}

func NewConfig(domain string, base string, driverConfig *DriverConfig, epoch uint64) *Config {
	return &Config{
		Domain:       domain,
		BaseFormat:   base,
		DriverConfig: driverConfig,
		EpochTime:    epoch,
	}
}
