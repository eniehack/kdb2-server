package config

type ElasticSearchConfig struct {
	Host string `toml:"host"`
}

type Config struct {
	ElasticSearchConfig `toml:"elasticsearch"`
}
