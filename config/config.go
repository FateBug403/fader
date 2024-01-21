package config

type Config struct {
	FoFa    FoFa       `mapstructure:"fofa" yaml:"fofa"`
	Collect Collect    `mapstructure:"collect" yaml:"collect"`
	Proxy   string     `mapstructure:"proxy" yaml:"proxy"`
	OneForAll string   `mapstructure:"oneforall" yaml:"oneforall"`
}
