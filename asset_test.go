package rtutils_lib

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestAssetService_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/asset",
		httpmock.NewStringResponder(201, `{"id": "1", "type": "asset"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	asset := &Asset{Name: "Laptop", Catalog: "IT"}
	id, err := client.Assets.Create(context.Background(), asset)
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
}

func TestAssetService_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"id": "1", "Name": "Laptop", "type": "asset"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	asset, err := client.Assets.Get(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, "Laptop", asset.Name)
}

func TestAssetService_Search(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initial Search Response
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/assets",
		httpmock.NewStringResponder(200, `{"total": 1, "items": [{"id": "1", "_url": "http://rt.example.com/REST/2.0/asset/1", "type": "asset"}]}`))

	// Follow-up Detail Fetch
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"id": "1", "Name": "Laptop", "type": "asset"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	result, err := client.Assets.Search(context.Background(), "Name like 'Laptop'")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, "Laptop", result.Items[0].Name)
}

func TestAssetService_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"message": "Updated"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	err := client.Assets.Update(context.Background(), "1", &Asset{Name: "NewName"})
	assert.NoError(t, err)
}

func TestAssetService_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"message": "Deleted"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	err := client.Assets.Delete(context.Background(), "1")
	assert.NoError(t, err)
}
