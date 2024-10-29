package main

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
	runAddress           = "localhost:8080"
	accrualSystemAddress = "localhost:8081"

	databaseURI string
)

func parseFlags() error {
	flag.Func("a", "address and port to run server in the form of host:port", func(flagValue string) error {
		if err := validateAddress(flagValue); err != nil {
			return fmt.Errorf("invalid server address: %s", err)
		}

		runAddress = flagValue
		return nil
	})

	flag.Func("d", "database URI", func(flagValue string) error {
		if flagValue == "" {
			return errors.New("invalid database URI")
		}

		databaseURI = flagValue
		return nil
	})

	flag.Func("r", "accrual system address", func(flagValue string) error {
		u, err := url.Parse(flagValue)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return errors.New("invalid URL format")
		}

		accrualSystemAddress = flagValue
		return nil
	})

	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		runAddress = envRunAddress
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		databaseURI = envDatabaseURI
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		accrualSystemAddress = envAccrualSystemAddress
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
