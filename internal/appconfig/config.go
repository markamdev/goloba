package appconfig

import (
	"fmt"
	"strings"

	"github.com/markamdev/goloba/pkg/balancer"
	"github.com/markamdev/goloba/pkg/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfig() (Config, error) {
	var currentConfig Config

	// default configuration values
	viper.SetDefault("listen-port", 8060)
	viper.SetDefault("targets", []string{})
	viper.SetDefault("algorithm", "round_robin")
	viper.SetDefault("log-level", "info")

	// environment variables loading
	viper.SetEnvPrefix("GOLOBA")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// command-line arguments loading
	pflag.Int("listen-port", 8060, "Port to listen on")
	pflag.StringSlice("targets", []string{}, "Comma-separated list of target servers")
	pflag.String("algorithm", "round_robin", "Load balancing algorithm (round_robin, least_connections)")
	pflag.String("log-level", "info", "Log level (debug, info, warn, error, fatal)")
	pflag.String("config-file", "", "Path to configuration file")
	// TODO in future replace '--help' flag with 'help' command
	pflag.Bool("help", false, "Print help screen")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// file-based configuration loading
	configFile := viper.GetString("config-file")
	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			return currentConfig, fmt.Errorf("error reading config file: %v", err)
		}

	}

	err := viper.Unmarshal(&currentConfig)
	if err != nil {
		return currentConfig, fmt.Errorf("unable to decode into struct: %v", err)
	}

	switch currentConfig.LogLevel {
	case "debug", "info", "warn", "error", "fatal":
		// valid log level
	default:
		return currentConfig, fmt.Errorf("invalid log level: %s", currentConfig.LogLevel)
	}

	return currentConfig, nil
}

func AppConfigToBalancerConfig(appCfg Config) balancer.Config {
	return balancer.Config{
		Port:      appCfg.ListenPort,
		Servers:   appCfg.Targets,
		Algorithm: appCfg.Algorithm,
		LogLevel:  string(logger.ParseLogLevel(appCfg.LogLevel)),
	}
}
