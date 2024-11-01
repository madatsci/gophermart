package flags

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	RunAddress           = "localhost:8080"
	AccrualSystemAddress = "localhost:8081"

	DatabaseURI string
)

func Parse() error {
	flag.Func("a", "address and port to run server in the form of host:port", func(flagValue string) error {
		if err := validateAddress(flagValue); err != nil {
			return fmt.Errorf("invalid server address: %s", err)
		}

		RunAddress = flagValue
		return nil
	})

	flag.Func("d", "database URI", func(flagValue string) error {
		if flagValue == "" {
			return errors.New("invalid database URI")
		}

		DatabaseURI = flagValue
		return nil
	})

	flag.Func("r", "accrual system address", func(flagValue string) error {
		u, err := url.Parse(flagValue)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return errors.New("invalid URL format")
		}

		AccrualSystemAddress = flagValue
		return nil
	})

	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		RunAddress = envRunAddress
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		DatabaseURI = envDatabaseURI
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		AccrualSystemAddress = envAccrualSystemAddress
	}

	return nil
}

func validateAddress(value string) error {
	hp := strings.Split(value, ":")
	if len(hp) != 2 {
		return errors.New("wrong address format, must be host:port")
	}

	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return errors.New("invalid port")
	}

	return nil
}
