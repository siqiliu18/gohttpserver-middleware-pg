package app

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DBHost string
	DBName string
	DBUser string
	DBPort string
	DBPass string
	DB     *sql.DB
	JwtKey string
}

var (
	dbhost = "localhost"
	dbname = "postgres"
	dbuser = "postgres"
	dbport = "5432"
	dbpass = "1234"
	jwtKey = "my-jwt-key"
)

func InitConfig() (*Config, error) {
	cfg := &Config{
		DBHost: dbhost,
		DBName: dbname,
		DBUser: dbuser,
		DBPort: dbport,
		DBPass: dbpass,
		JwtKey: jwtKey,
	}

	cfg.GETENVs()

	connStr := fmt.Sprintf(`host=%s dbname=%s user=%s port=%s password=%s sslmode=disable`, cfg.DBHost, cfg.DBName, cfg.DBUser, cfg.DBPort, cfg.DBPass)

	db, err := sql.Open("postgres", connStr)

	cfg.DB = db
	return cfg, err
}

func (c *Config) GETENVs() {
	if val, found := os.LookupEnv("CONFIG_DBHOST"); found {
		c.DBHost = val
	}
	if val, found := os.LookupEnv("CONFIG_DBPASS"); found {
		c.DBPass = val
	}
	if val, found := os.LookupEnv("JWT_KEY"); found {
		c.JwtKey = val
	}
}
