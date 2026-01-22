package appconfig

type Config struct {
	ListenPort uint     `mapstructure:"listen-port"`
	Targets    []string `mapstructure:"targets"`
	Algorithm  string   `mapstructure:"algorithm"`
	LogLevel   string   `mapstructure:"log-level"`
	ConfigFile string   `mapstructure:"config-file"`
	Help       bool     `mapstructure:"help"`
}
