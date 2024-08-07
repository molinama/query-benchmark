package logging

import (
	"fmt"

	"github.com/molinama/timescale/src/utils"
)

type config struct {
	encoding string
	level    string
}

const (
	logEncodingEnvVar = "LOG_ENCODING" // available values: console (default), json
	logLevelEnvVar    = "LOG_LEVEL"    //  available values: trace, debug, info (default), warn, error, fatal

	logEncodingDefault = "console"
	logLevelDefault    = "info"
)

func loadConfig() (*config, error) {
	fmt.Println("Load Logging configurations")
	return &config{
		encoding: utils.GetStringEnv(logEncodingEnvVar, logEncodingDefault),
		level:    utils.GetStringEnv(logLevelEnvVar, logLevelDefault),
	}, nil
}
