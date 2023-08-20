package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
    BucketID   string `envconfig:"BUCKET_ID" default:"bucket"`
    ObjectID   string `envconfig:"OBJECT_ID" default:"object"`
    BackupDir  string `envconfig:"BACKUP_DIR" default:"backup_dir"`
    BackupCron string `envconfig:"BACKUP_CRON" default:"*/5 * * * *"`

    HttpPort int `envconfig:"HTTP_PORT" default:"8080"`
}

func LoadConfigFromEnv() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("", c)
	return c, err
}

func MustLoadConfigFromEnv() *Config {
	c, err := LoadConfigFromEnv()
	if err != nil {
        panic(fmt.Errorf("failed to load config: %w", err))
	}
	return c
}

