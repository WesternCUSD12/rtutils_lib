package rtutils_lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func TestAssetService_SearchRaw_RejectsBroadQuery(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	result, err := client.Assets.SearchRaw(context.Background(), "Name LIKE '%'")
	assert.Error(t, err)
	assert.ErrorIs(t, err, errUnsafeBroadAssetQuery)
	assert.Nil(t, result)
}

func TestAssetService_SearchRaw_AllowsFilteredQuery(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/assets?query=Name+LIKE+%27W12-26%25%27",
		httpmock.NewStringResponder(200, `{"total": 1, "items": [{"id": "1", "Name": "W12-26001", "type": "asset"}]}`))

	result, err := client.Assets.SearchRaw(context.Background(), "Name LIKE 'W12-26%'")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
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

	searchURL := "http://rt.example.com/REST/2.0/assets?query=CustomField.%7BInternal+Name%7D+LIKE+%27%25COMP%25%27"
	httpmock.RegisterResponder("GET", searchURL,
		httpmock.NewStringResponder(200, `{"total": 1, "items": [{"id": "1", "_url": "http://rt.example.com/REST/2.0/asset/1", "type": "asset"}]}`))
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

func TestAssetService_ListCustomFieldValues(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"total":1,"items":[{"id":"9","Name":"Make","type":"customfield"}]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/9/values",
		httpmock.NewStringResponder(200, `{"total":2,"items":[{"id":"1","name":"Apple"},{"id":"2","name":"apple"},{"id":"3","name":"Lenovo"}]}`))

	values, err := client.Assets.ListCustomFieldValues(context.Background(), "Make")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Apple", "Lenovo"}, values)
}

func TestAssetService_ListCustomFieldValues_ReturnsErrorWhenPartialSearchFails(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"total":1,"items":[{"id":"7","Name":"Supplier","type":"customfield"}]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/7/values",
		httpmock.NewStringResponder(400, "custom field partial search unsupported"))

	values, err := client.Assets.ListCustomFieldValues(context.Background(), "Supplier")
	assert.Error(t, err)
	assert.Nil(t, values)
}

func TestAssetService_ListCustomFieldValuesMap(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/customfields", func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		payload := string(body)
		if strings.Contains(payload, "Make") {
			return httpmock.NewStringResponse(200, `{"total":1,"items":[{"id":"9","Name":"Make","type":"customfield"}]}`), nil
		}
		if strings.Contains(payload, "Model") {
			return httpmock.NewStringResponse(200, `{"total":1,"items":[{"id":"4","Name":"Model","type":"customfield"}]}`), nil
		}
		return httpmock.NewStringResponse(200, `{"total":0,"items":[]}`), nil
	})
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/9/values",
		httpmock.NewStringResponder(200, `{"total":1,"items":[{"id":"1","name":"Apple"}]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/4/values",
		httpmock.NewStringResponder(200, `{"total":1,"items":[{"id":"2","name":"iPad"}]}`))

	values, err := client.Assets.ListCustomFieldValuesMap(context.Background(), []string{"Make", "Model"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"Apple"}, values["Make"])
	assert.Equal(t, []string{"iPad"}, values["Model"])
}

func TestAssetService_ListCustomFieldValues_FallsBackToSampledAssets(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"total":1,"items":[{"id":"4","Name":"Model","type":"customfield"}]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/4/values",
		httpmock.NewStringResponder(200, `{"total":0,"items":[]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/assets?query=CustomField.%7BModel%7D+LIKE+%27%25%27",
		httpmock.NewStringResponder(200, `{"total":2,"items":[{"id":"101","_url":"http://rt.example.com/REST/2.0/asset/101","type":"asset"},{"id":"102","_url":"http://rt.example.com/REST/2.0/asset/102","type":"asset"}]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/101",
		httpmock.NewStringResponder(200, `{"id":"101","Name":"W12-26001","CustomFields":[{"name":"Model","values":["XPS 13"]}]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/102",
		httpmock.NewStringResponder(200, `{"id":"102","Name":"W12-26002","CustomFields":[{"name":"Model","values":["xps 13","ThinkPad"]}]}`))

	values, err := client.Assets.ListCustomFieldValues(context.Background(), "Model")
	assert.NoError(t, err)
	assert.Equal(t, []string{"ThinkPad", "XPS 13"}, values)
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

func TestAssetService_VerifyAssetCreated(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/42",
		httpmock.NewStringResponder(200, `{"id": "42", "Name": "Verified Laptop", "type": "asset"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	asset, err := client.Assets.VerifyAssetCreated(context.Background(), "42")
	assert.NoError(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "Verified Laptop", asset.Name)
}

func TestAssetService_SnapshotCustomFieldValues(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/customfields", func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		payload := string(body)
		if strings.Contains(payload, "Make") {
			return httpmock.NewStringResponse(200, `{"total":1,"items":[{"id":"9","Name":"Make","type":"customfield"}]}`), nil
		}
		if strings.Contains(payload, "Model") {
			return httpmock.NewStringResponse(200, `{"total":1,"items":[{"id":"4","Name":"Model","type":"customfield"}]}`), nil
		}
		return httpmock.NewStringResponse(200, `{"total":0,"items":[]}`), nil
	})
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/9/values",
		httpmock.NewStringResponder(200, `{"total":1,"items":[{"id":"1","name":"Apple"}]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/4/values",
		httpmock.NewStringResponder(200, `{"total":1,"items":[{"id":"2","name":"iPad"}]}`))

	snapshot, err := client.Assets.SnapshotCustomFieldValues(context.Background(), []string{"Make", "Model"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"Apple"}, snapshot["Make"])
	assert.Equal(t, []string{"iPad"}, snapshot["Model"])
}

func TestAssetService_BatchCreateAssets(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	createCalls := 0
	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/asset", func(req *http.Request) (*http.Response, error) {
		createCalls++
		if createCalls == 2 {
			return httpmock.NewStringResponse(500, `{"error":"failed"}`), nil
		}
		return httpmock.NewStringResponse(201, fmt.Sprintf(`{"id":"%d","type":"asset"}`, createCalls)), nil
	})

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	ids, errs := client.Assets.BatchCreateAssets(context.Background(), []*Asset{{Name: "A"}, {Name: "B"}, {Name: "C"}}, false)
	assert.Len(t, ids, 2)
	assert.Equal(t, []string{"1", "3"}, ids)
	assert.Len(t, errs, 1)
}

func TestAsset_MarshalJSON_CustomFieldsAsMap(t *testing.T) {
	asset := &Asset{
		Name:    "W12-26001",
		Catalog: "Technology",
		CustomFields: []AssetCustomField{
			{Name: "Internal Name", Values: []string{"Curly Bee"}},
			{Name: "Manufacturer", Values: []string{"Dell"}},
			{Name: "No Values"},
		},
	}
	data, err := json.Marshal(asset)
	assert.NoError(t, err)

	var out map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &out))

	cf, ok := out["CustomFields"].(map[string]interface{})
	assert.True(t, ok, "CustomFields should be a JSON object (map), not an array")
	assert.Equal(t, "Curly Bee", cf["Internal Name"])
	assert.Equal(t, "Dell", cf["Manufacturer"])
	_, hasEmpty := cf["No Values"]
	assert.False(t, hasEmpty, "field with no values should be omitted")
	assert.Equal(t, "Technology", out["Catalog"])
}

func TestAssetService_UpdateAsset(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", "http://rt.example.com/REST/2.0/asset/7",
		httpmock.NewStringResponder(200, `{"message": "Updated"}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/asset/7",
		httpmock.NewStringResponder(200, `{"id": "7", "Name": "Updated Asset", "type": "asset"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	asset, err := client.Assets.UpdateAsset(context.Background(), "7", &Asset{Name: "Updated Asset"})
	assert.NoError(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "Updated Asset", asset.Name)
}
