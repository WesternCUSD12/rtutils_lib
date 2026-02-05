package rtutils_lib

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
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

func TestAssetService_SearchByNameExact(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	searchURL := "http://rt.example.com/REST/2.0/assets"
	httpmock.RegisterResponder("POST", searchURL, func(req *http.Request) (*http.Response, error) {
		var criteria []map[string]interface{}
		body, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(body, &criteria)
		expected := []map[string]interface{}{{
			"field":    "Name",
			"operator": "=",
			"value":    "MacBook-1234",
		}}
		assert.Equal(t, expected, criteria)
		return httpmock.NewStringResponse(200, `{"total": 1, "items": [{"id": "1", "_url": "http://rt.example.com/REST/2.0/asset/1", "type": "asset"}]}`), nil
	})
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"id": "1", "Name": "MacBook-1234", "type": "asset"}`))

	result, err := client.Assets.SearchByNameExact(context.Background(), "MacBook-1234")
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, "MacBook-1234", result.Items[0].Name)
}

func TestAssetService_SearchByNamePartial(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	searchURL := "http://rt.example.com/REST/2.0/assets"
	httpmock.RegisterResponder("POST", searchURL, func(req *http.Request) (*http.Response, error) {
		var criteria []map[string]interface{}
		body, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(body, &criteria)
		expected := []map[string]interface{}{{
			"field":    "Name",
			"operator": "LIKE",
			"value":    "MacBook",
		}}
		assert.Equal(t, expected, criteria)
		return httpmock.NewStringResponse(200, `{"total": 2, "items": [{"id": "1", "_url": "http://rt.example.com/REST/2.0/asset/1", "type": "asset"}, {"id": "2", "_url": "http://rt.example.com/REST/2.0/asset/2", "type": "asset"}]}`), nil
	})
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"id": "1", "Name": "MacBook-1234", "type": "asset"}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/2",
		httpmock.NewStringResponder(200, `{"id": "2", "Name": "MacBook Pro", "type": "asset"}`))

	result, err := client.Assets.SearchByNamePartial(context.Background(), "MacBook")
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, "MacBook-1234", result.Items[0].Name)
}

func TestAssetService_SearchByCustomFieldExact(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	searchURL := "http://rt.example.com/REST/2.0/assets"
	httpmock.RegisterResponder("POST", searchURL, func(req *http.Request) (*http.Response, error) {
		var criteria []map[string]interface{}
		body, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(body, &criteria)
		expected := []map[string]interface{}{{
			"field":    "CustomField.{Internal Name}",
			"operator": "=",
			"value":    "COMP-12345",
		}}
		assert.Equal(t, expected, criteria)
		return httpmock.NewStringResponse(200, `{"total": 1, "items": [{"id": "1", "_url": "http://rt.example.com/REST/2.0/asset/1", "type": "asset"}]}`), nil
	})
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"id": "1", "Name": "Laptop", "type": "asset"}`))

	result, err := client.Assets.SearchByCustomFieldExact(context.Background(), "Internal Name", "COMP-12345")
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Total)
}

func TestAssetService_SearchByCustomFieldPartial(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	searchURL := "http://rt.example.com/REST/2.0/assets"
	httpmock.RegisterResponder("POST", searchURL, func(req *http.Request) (*http.Response, error) {
		var criteria []map[string]interface{}
		body, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(body, &criteria)
		expected := []map[string]interface{}{{
			"field":    "CustomField.{Internal Name}",
			"operator": "LIKE",
			"value":    "COMP",
		}}
		assert.Equal(t, expected, criteria)
		return httpmock.NewStringResponse(200, `{"total": 1, "items": [{"id": "1", "_url": "http://rt.example.com/REST/2.0/asset/1", "type": "asset"}]}`), nil
	})
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/1",
		httpmock.NewStringResponder(200, `{"id": "1", "Name": "Laptop", "type": "asset"}`))

	result, err := client.Assets.SearchByCustomFieldPartial(context.Background(), "Internal Name", "COMP")
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Total)
}

func TestAssetService_SearchByCustomFieldExact_NotFound(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	searchURL := "http://rt.example.com/REST/2.0/assets"
	httpmock.RegisterResponder("POST", searchURL,
		httpmock.NewStringResponder(400, "Custom field 'Internal Name' not found"))

	_, err := client.Assets.SearchByCustomFieldExact(context.Background(), "Internal Name", "COMP-404")
	assert.Error(t, err)
	assert.IsType(t, &CustomFieldNotFoundError{}, err)
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
