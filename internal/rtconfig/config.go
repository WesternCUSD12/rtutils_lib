package rtconfig

import (
	"fmt"
	"os"
	"strconv"
)

// RTConnection represents a connection to a Request Tracker instance
type RTConnection struct {
	URL     string
	Token   string
	Timeout int
}

// LoadFromEnv loads RT configuration from environment variables
func LoadFromEnv() (*RTConnection, error) {
	url := os.Getenv("RT_BASE_URL")
	token := os.Getenv("RT_TOKEN")
	timeout := 30
	if t := os.Getenv("RT_TIMEOUT"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil {
			timeout = parsed
		}
	}
	return &RTConnection{
		URL:     url,
		Token:   token,
		Timeout: timeout,
	}, nil
}

// Validate checks if the connection has required fields
func (c *RTConnection) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("RT_BASE_URL not set")
	}
	if c.Token == "" {
		return fmt.Errorf("RT_TOKEN not set")
	}
	return nil
}
