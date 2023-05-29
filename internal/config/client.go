package config

// ClientConfig struct for client config.
type ClientConfig struct {
	ServerAddress string
}

// NewClientConfig gets client config.
func NewClientConfig() ClientConfig {
	return ClientConfig{
		ServerAddress: ":3200",
	}
}
