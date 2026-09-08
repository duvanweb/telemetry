package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Configuration holds the application configuration loaded from environment.
type Configuration struct {
	HTTPPort          string `env:"HTTP_PORT" envDefault:"8081"`
	RedisAddr         string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	PostgresURL       string `env:"POSTGRES_URL" envDefault:"postgres://postgres:postgres@localhost:5432/geo_service?sslmode=disable"`
	CacheTTL          string `env:"GPS_CACHE_TTL" envDefault:"60s"`
	VehicleServiceURL string `env:"VEHICLE_SERVICE_URL" envDefault:"http://localhost:8080"`
	RabbitMQURL       string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
}

// LoadEnv loads configuration from environment variables into a struct of type T.
// Struct fields must be tagged with `env:"VAR_NAME" envDefault:"default_value"`.
func LoadEnv[T any]() (*T, error) {
	var cfg T
	v := reflect.ValueOf(&cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		value := os.Getenv(envTag)
		if value == "" {
			value = field.Tag.Get("envDefault")
		}

		if err := setField(v.Field(i), value); err != nil {
			return nil, fmt.Errorf("failed to set field %s from env %s: %w", field.Name, envTag, err)
		}
	}

	return &cfg, nil
}

// setField sets a reflect.Value from a string value based on its kind.
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		n, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(n))
	default:
		return fmt.Errorf("unsupported type %s", field.Kind())
	}

	return nil
}
