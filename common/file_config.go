package common

type Config struct {
	Domain       string          `json:"domain"`
	BaseFormat   string          `json:"base-format"`
	DriverConfig []*DriverConfig `json:"databases"`
	EpochTime    uint64          `json:"epoch-time"`
}

func NewConfig(domain string, base string, driverConfigs []*DriverConfig, epoch uint64) *Config {
	return &Config{
		Domain:       domain,
		BaseFormat:   base,
		DriverConfig: driverConfigs,
		EpochTime:    epoch,
	}
}
