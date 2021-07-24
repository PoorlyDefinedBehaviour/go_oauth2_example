package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type oAuth2Config struct {
	ClientID     string
	ClientSecret string
}

type configs struct {
	OAuth2 oAuth2Config
}

var Configs configs

func getFromEnvOrFail(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return value, errors.WithStack(fmt.Errorf("env key %s must have a value", key))
	}

	return value, nil
}

func ReadConfigsFromEnv() (err error) {
	Configs.OAuth2.ClientID, err = getFromEnvOrFail("OAUTH2_CLIENT_ID")
	if err != nil {
		return errors.WithStack(err)
	}

	Configs.OAuth2.ClientSecret, err = getFromEnvOrFail("OAUTH2_CLIENT_SECRET")
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
