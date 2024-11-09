package config

import "time"

type Config struct {
	RunAddress           string
	AccrualSystemAddress string
	DatabaseURI          string
	AccrualFetchPeriod   time.Duration

	TokenSecret    []byte
	TokenDuration  time.Duration
	TokenIssuer    string
	AuthCookieName string
}

// New creates new config
func New(runAddress, accrualSystemAddress, databaseURI string, tokenSecret []byte, tokenDuration time.Duration) *Config {
	return &Config{
		RunAddress:           runAddress,
		AccrualSystemAddress: accrualSystemAddress,
		DatabaseURI:          databaseURI,
		AccrualFetchPeriod:   20 * time.Second,

		TokenSecret:    tokenSecret,
		TokenDuration:  tokenDuration,
		TokenIssuer:    "gophermart",
		AuthCookieName: "auth_token",
	}
}
