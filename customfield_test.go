package rtutils_lib

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCustomField_UnmarshalJSON_NumericID(t *testing.T) {
	var cf CustomField
	err := json.Unmarshal([]byte(`{"id": 17, "Name": "Damages", "Type": "Freeform"}`), &cf)
	assert.NoError(t, err)
	assert.Equal(t, "17", cf.ID)
	assert.Equal(t, "Damages", cf.Name)
	assert.Equal(t, "Freeform", cf.Type)
}

func TestCustomFieldService_Get(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/5",
		httpmock.NewStringResponder(200, `{"id":"5","Name":"Manufacturer","Type":"FreeformSingle","Description":"Device manufacturer","Disabled":"0"}`))

	cf, err := client.CustomFields.Get(context.Background(), "5")
	assert.NoError(t, err)
	assert.Equal(t, "5", cf.ID)
	assert.Equal(t, "Manufacturer", cf.Name)
	assert.Equal(t, "FreeformSingle", cf.Type)
}

func TestCustomFieldService_Get_NotFound(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/999",
		httpmock.NewStringResponder(404, `{"message":"Not found"}`))

	_, err := client.CustomFields.Get(context.Background(), "999")
	assert.Error(t, err)
}

func TestCustomFieldService_Get_EmptyID(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	_, err := client.CustomFields.Get(context.Background(), "")
	assert.Error(t, err)
}

func TestCustomFieldService_GetByCategory(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/7?category=Hardware",
		httpmock.NewStringResponder(200, `{"id":"7","Name":"Asset Type","Type":"SelectSingle","Values":["Laptop","Desktop","Tablet"]}`))

	cf, err := client.CustomFields.GetByCategory(context.Background(), "7", "Hardware")
	assert.NoError(t, err)
	assert.Equal(t, "Asset Type", cf.Name)
	assert.Equal(t, []string{"Laptop", "Desktop", "Tablet"}, cf.Values)
}

func TestCustomFieldService_ListAll(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"count":2,"total":2,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"1","Name":"Internal Name","Type":"FreeformSingle"},
			{"id":"2","Name":"Manufacturer","Type":"FreeformSingle"}
		]}`))

	fields, err := client.CustomFields.ListAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, fields, 2)
	assert.Equal(t, "Internal Name", fields[0].Name)
	assert.Equal(t, "Manufacturer", fields[1].Name)
}

func TestCustomFieldService_ListAll_AutoExpands(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	// Simulate RT returning stubs (ID + _url only, no Name)
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"count":2,"total":2,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"1","_url":"http://rt.example.com/REST/2.0/customfield/1"},
			{"id":"2","_url":"http://rt.example.com/REST/2.0/customfield/2"}
		]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/1",
		httpmock.NewStringResponder(200, `{"id":"1","Name":"Internal Name","Type":"FreeformSingle"}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/2",
		httpmock.NewStringResponder(200, `{"id":"2","Name":"Manufacturer","Type":"FreeformSingle"}`))

	fields, err := client.CustomFields.ListAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, fields, 2)
	// After expansion, names should be populated
	names := []string{fields[0].Name, fields[1].Name}
	assert.Contains(t, names, "Internal Name")
	assert.Contains(t, names, "Manufacturer")
}

func TestCustomFieldService_ListAllForCatalog_AutoExpands(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/catalog/2/customfields",
		httpmock.NewStringResponder(200, `{"count":2,"total":2,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"3","_url":"http://rt.example.com/REST/2.0/customfield/3"},
			{"id":"4","_url":"http://rt.example.com/REST/2.0/customfield/4"}
		]}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/3",
		httpmock.NewStringResponder(200, `{"id":"3","Name":"Serial Number","Type":"FreeformSingle"}`))
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfield/4",
		httpmock.NewStringResponder(200, `{"id":"4","Name":"Warranty Expiration","Type":"Date"}`))

	fields, err := client.CustomFields.ListAllForCatalog(context.Background(), "2")
	assert.NoError(t, err)
	assert.Len(t, fields, 2)
	names := []string{fields[0].Name, fields[1].Name}
	assert.Contains(t, names, "Serial Number")
	assert.Contains(t, names, "Warranty Expiration")
}

func TestCustomFieldService_List_WithCriteria(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"count":1,"total":1,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"3","Name":"Model","Type":"FreeformSingle"}
		]}`))

	criteria := []map[string]interface{}{
		{"field": "Name", "operator": "=", "value": "Model"},
	}
	result, err := client.CustomFields.List(context.Background(), criteria)
	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "Model", result.Items[0].Name)
}

func TestCustomFieldService_ListAllForCatalog(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/catalog/2/customfields",
		httpmock.NewStringResponder(200, `{"count":3,"total":3,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"1","Name":"Internal Name","Type":"FreeformSingle"},
			{"id":"2","Name":"Manufacturer","Type":"FreeformSingle"},
			{"id":"4","Name":"Warranty Expiration","Type":"Date"}
		]}`))

	fields, err := client.CustomFields.ListAllForCatalog(context.Background(), "2")
	assert.NoError(t, err)
	assert.Len(t, fields, 3)
	assert.Equal(t, "Warranty Expiration", fields[2].Name)
}

func TestCustomFieldService_ListForCatalog_EmptyID(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	_, err := client.CustomFields.ListAllForCatalog(context.Background(), "")
	assert.Error(t, err)
}

func TestCustomFieldService_ListAllForQueue(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/queue/1/customfields",
		httpmock.NewStringResponder(200, `{"count":1,"total":1,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"10","Name":"Severity","Type":"SelectSingle","Values":["Low","Medium","High"]}
		]}`))

	fields, err := client.CustomFields.ListAllForQueue(context.Background(), "1")
	assert.NoError(t, err)
	assert.Len(t, fields, 1)
	assert.Equal(t, "Severity", fields[0].Name)
}

func TestCustomFieldService_ListAllForClass(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/class/4/customfields",
		httpmock.NewStringResponder(200, `{"count":1,"total":1,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"11","Name":"Department","Type":"FreeformSingle"}
		]}`))

	fields, err := client.CustomFields.ListAllForClass(context.Background(), "4")
	assert.NoError(t, err)
	assert.Len(t, fields, 1)
	assert.Equal(t, "Department", fields[0].Name)
}

func TestCustomFieldService_FindByName(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"count":2,"total":2,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"1","Name":"Internal Name","Type":"FreeformSingle"},
			{"id":"2","Name":"Manufacturer","Type":"FreeformSingle"}
		]}`))

	cf, err := client.CustomFields.FindByName(context.Background(), "manufacturer")
	assert.NoError(t, err)
	assert.Equal(t, "Manufacturer", cf.Name)
	assert.Equal(t, "2", cf.ID)
}

func TestCustomFieldService_FindByName_NotFound(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/customfields",
		httpmock.NewStringResponder(200, `{"count":1,"total":1,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"1","Name":"Internal Name","Type":"FreeformSingle"}
		]}`))

	_, err := client.CustomFields.FindByName(context.Background(), "NonExistent")
	assert.Error(t, err)
	var cfErr *CustomFieldNotFoundError
	assert.ErrorAs(t, err, &cfErr)
	assert.Equal(t, "NonExistent", cfErr.FieldName)
}

func TestCustomFieldService_NamesForCatalog(t *testing.T) {
	client, cleanup := testClient(t, "http://rt.example.com/REST/2.0")
	defer cleanup()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/catalog/2/customfields",
		httpmock.NewStringResponder(200, `{"count":3,"total":3,"page":1,"pages":1,"per_page":20,"items":[
			{"id":"1","Name":"Internal Name"},
			{"id":"2","Name":"Manufacturer"},
			{"id":"3","Name":"Serial Number"}
		]}`))

	names, err := client.CustomFields.NamesForCatalog(context.Background(), "2")
	assert.NoError(t, err)
	assert.Equal(t, []string{"Internal Name", "Manufacturer", "Serial Number"}, names)
}
