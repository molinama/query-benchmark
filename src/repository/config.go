package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	logEncodingEnvVar = "LOG_ENCODING" // available values: console (default), json
	logLevelEnvVar    = "LOG_LEVEL"    //  available values: trace, debug, info (default), warn, error, fatal

	logEncodingDefault = "console"
	logLevelD
	efault = "info"
)

type config struct {
	dbConnString string
}

func loadConfig() (*config, error) {
	err := godotenv.Load("./timescaledb/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("TIMESCALE_USER")
	pwd := os.Getenv("TIMESCALE_PASSWORD")
	host := os.Getenv("TIMESCALE_HOST")
	db := os.Getenv("TIMESCALES_DB")
	port := os.Getenv("TIMESCALE_PORT")

	dbConnString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		pwd,
		host,
		port,
		db)

	return &config{
		dbConnString: dbConnString,
	}, nil
}
