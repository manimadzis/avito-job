package config

import "fmt"

func newErrCantParseConfig(err error) error {
	return fmt.Errorf("can't load config: %v", err)
}

func newErrCantLoadConfig(err error) error {
	return fmt.Errorf("can't load config: %v", err)
}
