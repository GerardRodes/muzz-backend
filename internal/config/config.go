package config

import (
	"log"
	"os"
	"reflect"
)

type Config struct {
	DBAddr     string `env:"MUZZ_DB_ADDR"`
	DBName     string `env:"MUZZ_DB_NAME"`
	DBUser     string `env:"MUZZ_DB_USER"`
	DBPassword string `env:"MUZZ_DB_PASSWORD"`
	HTTPPort   string `env:"MUZZ_HTTP_PORT"`
	KVAddr     string `env:"MUZZ_KV_ADDR"`
}

func New() Config {
	cfg := &Config{}
	cfg.Init()
	return *cfg
}

// Init populates config from env
func (c *Config) Init() {
	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()

	for idx := 0; idx < t.NumField(); idx++ {
		envVar, ok := t.Field(idx).Tag.Lookup("env")
		if !ok {
			log.Fatalf("config init: missing env struct tag for field %q", t.Field(idx).Name)
		}

		envVarVal := os.Getenv(envVar)
		if envVarVal == "" {
			log.Fatalf("config init: env missing var %q", envVar)
		}
		v.Field(idx).SetString(envVarVal)
	}
}
