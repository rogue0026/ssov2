package ssoconfig

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type SSOConfig struct {
	RunningEnv string        `yaml:"running_env"`
	Address    string        `yaml:"address"`
	DSN        string        `yaml:"dsn"`
	TokenTTL   time.Duration `yaml:"token_ttl"`
}

func (c SSOConfig) String() string {
	return fmt.Sprintf("%s %s %s %s", c.RunningEnv, c.Address, c.DSN, c.TokenTTL.String())
}

func MustLoad(ssoCfgPath string) SSOConfig {
	_, err := os.Stat(ssoCfgPath)
	if err != nil {
		panic(err.Error())
	}
	cfgData := SSOConfig{}
	if err = cleanenv.ReadConfig(ssoCfgPath, &cfgData); err != nil {
		panic(err.Error())
	}
	return cfgData
}
