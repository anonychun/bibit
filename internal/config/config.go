package config

import (
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

var cfg *Config

type Config struct {
	Http struct {
		Port int `envconfig:"port"`
	} `envconfig:"http"`

	Grpc struct {
		Port int `envconfig:"port"`
	} `envconfig:"grpc"`

	DB struct {
		Sql struct {
			Host     string `envconfig:"host"`
			Port     int    `envconfig:"port"`
			User     string `envconfig:"user"`
			Password string `envconfig:"password"`
			Name     string `envconfig:"name"`
		} `envconfig:"sql"`
	} `envconfig:"db"`

	Storage struct {
		S3 struct {
			Endpoint        string        `envconfig:"endpoint"`
			Bucket          string        `envconfig:"bucket"`
			AccessKeyId     string        `envconfig:"access_key_id"`
			SecretAccessKey string        `envconfig:"secret_access_key"`
			UrlExpiration   time.Duration `envconfig:"url_expiration"`
		} `envconfig:"s3"`
	} `envconfig:"storage"`

	Broker struct {
		Backend string `envconfig:"backend"`

		Kafka struct {
			Brokers      []string      `envconfig:"brokers"`
			GroupId      string        `envconfig:"group_id"`
			BatchSize    int           `envconfig:"batch_size"`
			BatchTimeout time.Duration `envconfig:"batch_timeout"`
		} `envconfig:"kafka"`

		Rabbitmq struct {
			Url          string        `envconfig:"url"`
			Exchange     string        `envconfig:"exchange"`
			GroupId      string        `envconfig:"group_id"`
			BatchSize    int           `envconfig:"batch_size"`
			BatchTimeout time.Duration `envconfig:"batch_timeout"`
		} `envconfig:"rabbitmq"`
	} `envconfig:"broker"`

	OTLP struct {
		Endpoint string `envconfig:"endpoint"`
	} `envconfig:"otlp"`
}

func Setup() error {
	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		return err
	}

	cfg = config
	return nil
}

func Get() *Config {
	if cfg == nil {
		err := Setup()
		if err != nil {
			panic(err)
		}
	}

	return cfg
}
