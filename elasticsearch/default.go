package elasticsearch

type (
	Config struct {
		username           string
		password           string
		isUseElasticSearch bool
		address            []string
	}
	Option func(c *Config)
)

func NewES(opts ...Option) *Config {
	c := &Config{
		isUseElasticSearch: false,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func Address(address []string) Option {
	return func(c *Config) {
		c.address = address
	}
}

func IsUseElasticSearch(isUseElasticSearch bool) Option {
	return func(c *Config) {
		c.isUseElasticSearch = isUseElasticSearch
	}
}

func UserName(username string) Option {
	return func(c *Config) {
		c.username = username
	}
}

func Password(password string) Option {
	return func(c *Config) {
		c.password = password
	}
}
