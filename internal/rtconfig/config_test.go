package rtconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnvWithAllVars(t *testing.T) {
	t.Setenv("RT_URL", "https://rt.example.com")
	t.Setenv("RT_USERNAME", "testuser")
	t.Setenv("RT_PASSWORD", "testpass")
	t.Setenv("RT_TIMEOUT", "60")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "https://rt.example.com", conn.URL)
	assert.Equal(t, "testuser", conn.Username)
	assert.Equal(t, "testpass", conn.Password)
	assert.Equal(t, 60, conn.Timeout)
}

func TestLoadFromEnvWithDefaults(t *testing.T) {
	t.Setenv("RT_URL", "https://rt.example.com")
	t.Setenv("RT_USERNAME", "testuser")
	t.Setenv("RT_PASSWORD", "testpass")
	os.Unsetenv("RT_TIMEOUT")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, 30, conn.Timeout)
}

func TestLoadFromEnvMissingCredentials(t *testing.T) {
	os.Unsetenv("RT_URL")
	os.Unsetenv("RT_USERNAME")
	os.Unsetenv("RT_PASSWORD")
	os.Unsetenv("RT_TIMEOUT")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, "", conn.URL)
	assert.Equal(t, "", conn.Username)
	assert.Equal(t, "", conn.Password)
}

func TestValidateWithValidConnection(t *testing.T) {
	conn := &RTConnection{
		URL:      "https://rt.example.com",
		Username: "testuser",
		Password: "testpass",
		Timeout:  30,
	}
	err := conn.Validate()
	assert.NoError(t, err)
}

func TestValidateMissingURL(t *testing.T) {
	conn := &RTConnection{
		Username: "testuser",
		Password: "testpass",
		Timeout:  30,
	}
	err := conn.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RT_URL")
}

func TestValidateMissingUsername(t *testing.T) {
	conn := &RTConnection{
		URL:      "https://rt.example.com",
		Password: "testpass",
		Timeout:  30,
	}
	err := conn.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RT_USERNAME")
}

func TestValidateMissingPassword(t *testing.T) {
	conn := &RTConnection{
		URL:      "https://rt.example.com",
		Username: "testuser",
		Timeout:  30,
	}
	err := conn.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "RT_PASSWORD")
}

func TestLoadFromEnvInvalidTimeout(t *testing.T) {
	t.Setenv("RT_URL", "https://rt.example.com")
	t.Setenv("RT_USERNAME", "testuser")
	t.Setenv("RT_PASSWORD", "testpass")
	t.Setenv("RT_TIMEOUT", "not_a_number")

	conn, err := LoadFromEnv()
	require.NoError(t, err)
	assert.Equal(t, 30, conn.Timeout)
}
