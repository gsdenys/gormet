package gormet

type Config struct {
	Validate bool
}

func DefaultConfig() *Config {
	return &Config{
		Validate: true,
	}
}
