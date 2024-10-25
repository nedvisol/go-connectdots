package config

import (
	"os"
	"time"
)

type GraphDbConfig struct {
	Uri      string
	Username string
	Password string
}

type Config struct {
	CacheDir         string
	CacheTtl         time.Duration
	MongoUrl         string
	MongoDb          string
	CongressGovToken string
	GraphDb          *GraphDbConfig
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
		GraphDb: &GraphDbConfig{
			Uri:      "neo4j://nedlinux:7687",
			Username: "neo4j",
			Password: "neo4jpassword",
		},
	}
}
