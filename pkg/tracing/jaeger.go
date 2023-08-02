package tracing

import (
	"disrupt/pkg/constants"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

var configPath string

type Config struct {
	Environment string `mapstructure:"environment"`
	ID          int64  `mapstructure:"id"`
	Url         string `mapstructure:"url"`
	HostPort    string `mapstructure:"hostPort"`
	Enable      bool   `mapstructure:"enable"`
	LogSpans    bool   `mapstructure:"logSpans"`
}

type Tracer trace.Tracer

func TracerProvider(url string, serviceName string) (*tracesdk.TracerProvider, error) {
	config, err := InitConfig()
	if err != nil {
		return nil, err
	}
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			attribute.String("environment", config.Environment),
			attribute.Int64("ID", config.ID),
		)),
	)
	return tp, nil
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/pkg/tracing/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}
