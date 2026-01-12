package appconfig

type Config struct {
	ListenPort uint     `mapstructure:"listen_port"`
	Targets    []string `mapstructure:"targets"`
	Algorithm  string   `mapstructure:"algorithm"`
	LogLevel   string   `mapstructure:"log_level"`
	ConfigFile string   `mapstructure:"config_file"`
	Help       bool     `mapstructure:"help"`
}
