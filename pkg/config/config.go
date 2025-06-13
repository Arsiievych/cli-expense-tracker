package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ExpensesFilePath string `mapstructure:"expenses_file_path"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./pkg/config")

	//v.SetDefault("data_file_path", "./expenses.json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}

	return &config, nil
}
