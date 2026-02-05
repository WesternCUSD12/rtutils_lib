package rtconfig

import (
	"fmt"
	"os"
	"strconv"
)

// RTConnection represents a connection to a Request Tracker instance
type RTConnection struct {
	URL      string
	Username string
	Password string
	Timeout  int
}

// LoadFromEnv loads RT configuration from environment variables
func LoadFromEnv() (*RTConnection, error) {
	url := os.Getenv("RT_URL")
	username := os.Getenv("RT_USERNAME")
	password := os.Getenv("RT_PASSWORD")
	timeout := 30
	if t := os.Getenv("RT_TIMEOUT"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil {
			timeout = parsed
		}
	}
	return &RTConnection{
		URL:      url,
		Username: username,
		Password: password,
		Timeout:  timeout,
	}, nil
}

// Validate checks if the connection has required fields
func (c *RTConnection) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("RT_URL not set")
	}
	if c.Username == "" {
		return fmt.Errorf("RT_USERNAME not set")
	}
	if c.Password == "" {
		return fmt.Errorf("RT_PASSWORD not set")
	}
	return nil
}
