package nested_comments

import (
	"errors"
	"os"
	"strconv"
)

type WithConfig func(*Config) error

type Config struct {
	Port      uint64
	MysqlConn string
}

func LoadConfig(configOpt ...WithConfig) (Config, error) {
	// predefine config
	ac := &Config{
		Port:      3030,
		MysqlConn: "root:@tcp(127.0.0.1)/comment_service",
	}

	for _, opt := range configOpt {
		err := opt(ac)
		if err != nil {
			return Config{}, err
		}
	}

	return *ac, nil
}

var ErrorInvalidPort = errors.New("invalid port")

func WithPort(env string, defaultPort uint64) WithConfig {
	if port, exists := os.LookupEnv(env); exists {
		if p, err := strconv.ParseUint(port, 10, 64); err == nil {
			return func(c *Config) error {
				c.Port = p
				return nil
			}
		} else {
			return func(c *Config) error {
				return ErrorInvalidPort
			}
		}
	}

	return func(c *Config) error {
		c.Port = defaultPort
		return nil
	}
}

func WithMysqlConn(env string, defaultConn string) WithConfig {
	if conn, exists := os.LookupEnv(env); exists {
		return func(c *Config) error {
			c.MysqlConn = conn
			return nil
		}
	}

	return func(c *Config) error {
		c.MysqlConn = defaultConn
		return nil
	}
}
