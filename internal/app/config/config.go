package config

type Config struct {
	RunAddress           string
	AccrualSystemAddress string
	DatabaseURI          string
}

func New(runAddress, accrualSystemAddress, databaseURI string) *Config {
	return &Config{
		RunAddress:           runAddress,
		AccrualSystemAddress: accrualSystemAddress,
		DatabaseURI:          databaseURI,
	}
}
