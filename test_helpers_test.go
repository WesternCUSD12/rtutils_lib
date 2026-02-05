package rtutils_lib

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

// testClient activates httpmock on a client and returns a cleanup function.
func testClient(t *testing.T, baseURL string) (*Client, func()) {
	t.Helper()
	client := NewClient(baseURL, "test-token")
	httpmock.ActivateNonDefault(client.client)
	cleanup := func() {
		httpmock.DeactivateAndReset()
	}
	return client, cleanup
}

// activateHTTPMock activates httpmock for non-default clients.
func activateHTTPMock() func() {
	httpmock.Activate()
	return func() {
		httpmock.DeactivateAndReset()
	}
}
