package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	GRPC        GRPCConfig `yaml:"grpc"`
}
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	cfgpath := fetchCfgPath()
	if cfgpath == "" {
		panic("config path is empty")
	}
	return MustLoadByPath(cfgpath)
}
func MustLoadByPath(cfgpath string) *Config {
	if _, err := os.Stat(cfgpath); err != nil {
		panic("config file does not exits" + err.Error())
	}
	var cfg Config

	if err := cleanenv.ReadConfig(cfgpath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}

func fetchCfgPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
