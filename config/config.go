package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

const (
	defaultConfigPath = "./config/config.yaml"
	defaultEnvPath    = ".env"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		PG   `yaml:"pg"`
		RMQ  `yaml:"rabbitmq"`
		Log  `yaml:"logger"`
	}

	App struct {
		Name           string        `env:"APP_NAME"            env-default:"wallet-rielta" yaml:"name"`
		Version        string        `env:"APP_VERSION"         env-default:"1.0.0"         yaml:"version"`
		CountWorkers   int           `env:"APP_WORKERS"         env-default:"24"            yaml:"workers"`
		Timeout        time.Duration `env:"APP_TIMEOUT"         env-default:"5s"            yaml:"timeout"`
		DefaultBalance uint          `env:"APP_DEFAULT_BALANCE" env-default:"100"           yaml:"defaultBalance"`
	}

	HTTP struct {
		Port    string        `env:"HTTP_PORT"    env-default:":8080" yaml:"port"`
		Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"5s"    yaml:"timeout"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX" env-default:"2"     yaml:"poolMax"`
		URL     string `env:"PG_URL"      env-required:"true" yaml:"url"`
	}

	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER" env-default:"rpc_server" yaml:"rpcServerExchange"`
		ClientExchange string `env:"RMQ_RPC_CLIENT" env-default:"rpc_client" yaml:"rpcClientExchange"`
		URL            string `env:"RMQ_URL"        env-required:"true"      yaml:"url"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug" yaml:"logLevel"`
	}
)

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath, defaultEnvPath)
}

func MustLoadPath(configPath, envPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	// try loading .env file
	_ = godotenv.Load(envPath)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = defaultConfigPath
	}

	return res
}
