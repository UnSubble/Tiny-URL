package common

type Config struct {
	Url          string          `json:"url"`
	BaseFormat   string          `json:"base-format"`
	DriverConfig []*DriverConfig `json:"databases"`
	EpochTime    uint64          `json:"epoch-time"`
}

func NewConfig(url string, base string, driverConfigs []*DriverConfig, epoch uint64) *Config {
	return &Config{
		Url:          url,
		BaseFormat:   base,
		DriverConfig: driverConfigs,
		EpochTime:    epoch,
	}
}
