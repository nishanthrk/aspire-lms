package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

const (
	prod = "production"
	dev  = "development"
	stg  = "staging"
)

// Config object
type Config struct {
	Tenant           string      `env:"TENANT" required:"true"`
	Env              string      `env:"ENV" required:"true"`
	Mysql            MysqlConfig `json:"mysql"`
	JWTAccessSecret  string      `env:"JWT_ACCESS_SIGN_KEY" required:"true"`
	JWTRefreshSecret string      `env:"JWT_REFRESH_SIGN_KEY" required:"true"`
	Host             string      `env:"APP_HOST" required:"true"`
	Port             string      `env:"APP_PORT" required:"true"`
	DbHost           string      `env:"DB_HOST" required:"true"`
	DbPort           string      `env:"DB_PORT" required:"true"`
	DbDriver         string      `env:"DB_DRIVER" required:"true"`
	DbUser           string      `env:"DB_USER" required:"true"`
	DbPassword       string      `env:"DB_PASSWORD" required:"true"`
	DbName           string      `env:"DB_NAME" required:"true"`
	NewRelicLicense  string      `env:"NEW_RELIC_LICENSE"`
}

// IsProd Checks if env is production
func (c Config) IsProd() bool {
	return c.Env == prod
}

func (c Config) IsDev() bool {
	return c.Env == dev
}

func (c Config) IsStg() bool {
	return c.Env == stg
}

func (c Config) GetTenantName() string {
	return strings.ToUpper(fmt.Sprintf("%v_%v", c.Tenant, c.Env))
}

// LoadLocalConfig gets config from .env
func LoadLocalConfig() {
	requiredEnvVars := getRequiredEnvVars(Config{})

	var missingVars []string

	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	environmentPath := filepath.Join(currentPath, ".env")

	if err := godotenv.Load(environmentPath); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	for _, envVar := range requiredEnvVars {
		if getEnv(envVar.tag) == "" && envVar.required {
			missingVars = append(missingVars, envVar.tag)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("missing environment variables: %s", strings.Join(missingVars, ", "))
	}
}

func getRequiredEnvVars(cfgStruct interface{}) []ConfigLoad {
	var envVars []ConfigLoad
	val := reflect.ValueOf(cfgStruct)

	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		_load := ConfigLoad{}
		tag := field.Tag.Get("env")
		required := field.Tag.Get("required")

		_load.tag = tag
		_load.required = tag != "" && required == "true"

		if tag != "" {
			envVars = append(envVars, _load)
		}
	}

	return envVars
}

type ConfigLoad struct {
	tag      string
	required bool
}

// GetConfig gets all config for the application
func GetConfig() Config {
	return Config{
		Tenant:           getEnv("TENANT"),
		Env:              getEnv("ENV"),
		Mysql:            GetMysqlConfig(),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SIGN_KEY"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SIGN_KEY"),
		Host:             getEnv("APP_HOST"),
		Port:             getEnv("APP_PORT"),
		DbHost:           getEnv("DB_HOST"),
		DbPort:           getEnv("DB_PORT"),
		DbDriver:         getEnv("DB_DRIVER"),
		DbUser:           getEnv("DB_USER"),
		DbPassword:       getEnv("DB_PASSWORD"),
		DbName:           getEnv("DB_NAME"),
		NewRelicLicense:  getEnv("NEW_RELIC_LICENSE"),
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}
