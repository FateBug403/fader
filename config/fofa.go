package config

type FoFa struct {
	Api string `mapstructure:"api" yaml:"api"`
	Mail string `mapstructure:"mail" yaml:"mail"`
	Key string  `mapstructure:"key" yaml:"key"`
}
