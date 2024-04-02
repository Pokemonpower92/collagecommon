package db

import "os"

type DBConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	DbName   string
}

func NewDBConfig(db_name string) DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		DbName:   db_name,
	}
}
