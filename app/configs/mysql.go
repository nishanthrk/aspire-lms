package configs

import (
	"fmt"
	"os"
	"strconv"
)

// MysqlConfig object
type MysqlConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

// Dialect returns "mysql"
func (c MysqlConfig) Dialect() string {
	return "mysql"
}

// GetMysqlConnectionInfo returns Mysql URL string
func (c MysqlConfig) GetMysqlConnectionInfo() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FKolkata",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func (c MysqlConfig) GetMysqlConnectionForMigrate() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

// GetMysqlConfig returns MysqlConfig object
func GetMysqlConfig() MysqlConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return MysqlConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
