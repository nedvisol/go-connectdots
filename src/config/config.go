package config

import (
	"os"
	"time"
)

type Config struct {
	CacheDir         string
	CacheTtl         time.Duration
	MongoUrl         string
	MongoDb          string
	CongressGovToken string
}

func NewConfig() *Config {
	congressApiToken, err := os.ReadFile("../.tmp/congressApiToken.txt")
	if err != nil {
		panic(err)
	}
	return &Config{
		CacheDir:         "../.tmp/cache",
		CacheTtl:         time.Hour * 240,
		MongoUrl:         "mongodb://nedlinux:27017",
		MongoDb:          "go_connectdots",
		CongressGovToken: string(congressApiToken),
	}
}
