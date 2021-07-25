package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Config struct {
	apiPort           string
	ServerReadTimeout int

	JwtJwksUrl     string
	jwtKeycloakUrl string

	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string

	LogFormat string

	Migrate string
}

func New() *Config {
	conf := &Config{}

	flag.StringVar(&conf.apiPort, "apiport", os.Getenv("API_PORT"), "Which port for the API server to listen on")

	serverReadTimeout, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		zap.S().Errorf("Could not read SERVER_READ_TIMEOUT: %v", err)
		serverReadTimeout = 60
	}
	flag.IntVar(&conf.ServerReadTimeout, "serverreadtimeout", serverReadTimeout, "Server read timeout")

	flag.StringVar(&conf.JwtJwksUrl, "jwtjwksurl", os.Getenv("JWT_JWKS_URL"), "JWT JWKS URL")
	flag.StringVar(&conf.jwtKeycloakUrl, "jwtkeycloakurl", os.Getenv("JWK_KEYCLOAK_URL"), "JWT Keycloak URL")

	flag.StringVar(&conf.dbHost, "dbhost", os.Getenv("DB_HOST"), "Database host")
	flag.StringVar(&conf.dbPort, "dbport", os.Getenv("DB_PORT"), "Database port")
	flag.StringVar(&conf.dbName, "dbname", os.Getenv("DB_NAME"), "Database name")
	flag.StringVar(&conf.dbUser, "dbuser", os.Getenv("DB_USER"), "Database user")
	flag.StringVar(&conf.dbPassword, "dbpassword", os.Getenv("DB_PASSWORD"), "Database password")

	flag.StringVar(&conf.LogFormat, "logformat", os.Getenv("LOG_FORMAT"), "Log format (json or console)")

	flag.StringVar(&conf.Migrate, "migrate", "up", "Specify if we should migrate DB 'up' or 'down'")

	flag.Parse()

	return conf
}

func (c *Config) GetDBConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.dbUser, c.dbPassword, c.dbHost, c.dbPort, c.dbName,
	)
}

func (c *Config) GetServerUrl() string {
	return fmt.Sprintf("0.0.0.0:%s", c.apiPort)
}
