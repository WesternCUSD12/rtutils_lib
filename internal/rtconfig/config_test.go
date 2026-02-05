package rtconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnvWithAllVars(t *testing.T) {
	t.Setenv("RT_BASE_URL", "https://rt.example.com")
	t.Setenv("RT_TOKEN", "test-token-abc123")
	t.Setenv("RT_TIMEOUT", "60")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "https://rt.example.com", conn.URL)
	assert.Equal(t, "test-token-abc123", conn.Token)
	assert.Equal(t, 60, conn.Timeout)
}

func TestLoadFromEnvWithDefaults(t *testing.T) {
	t.Setenv("RT_BASE_URL", "https://rt.example.com")
	t.Setenv("RT_TOKEN", "test-token")
	os.Unsetenv("RT_TIMEOUT")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, 30, conn.Timeout)
}

func TestLoadFromEnvMissingCredentials(t *testing.T) {
	os.Unsetenv("RT_BASE_URL")
	os.Unsetenv("RT_TOKEN")
	os.Unsetenv("RT_TIMEOUT")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "", conn.URL)
	assert.Equal(t, "", conn.Token)
}

func TestValidateWithValidConnection(t *testing.T) {
	conn := &RTConnection{
		URL:     "https://rt.example.com",
		Token:   "test-token",
		Timeout: 30,
	}
	err := conn.Validate()
	assert.NoError(t, err)
}

func TestValidateMissingURL(t *testing.T) {
	conn := &RTConnection{
		Token:   "test-token",
		Timeout: 30,
	}
	err := conn.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RT_BASE_URL")
}

func TestValidateMissingToken(t *testing.T) {
	conn := &RTConnection{
		URL:     "https://rt.example.com",
		Timeout: 30,
	}
	err := conn.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RT_TOKEN")
}

func TestLoadFromEnvInvalidTimeout(t *testing.T) {
	t.Setenv("RT_BASE_URL", "https://rt.example.com")
	t.Setenv("RT_TOKEN", "test-token")
	t.Setenv("RT_TIMEOUT", "not_a_number")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, 30, conn.Timeout)
}
