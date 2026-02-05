package rtutils_lib

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/user",
		httpmock.NewStringResponder(201, `{"id": "u1", "type": "user"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	user := &User{Name: "jdoe", EmailAddress: "jdoe@example.com"}
	id, err := client.Users.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, "u1", id)
}

func TestUserService_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/user/jdoe",
		httpmock.NewStringResponder(200, `{"id": "u1", "Name": "jdoe", "type": "user"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	user, err := client.Users.Get(context.Background(), "jdoe")
	assert.NoError(t, err)
	assert.Equal(t, "jdoe", user.Name)
}

func TestUserService_Search(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/users",
		httpmock.NewStringResponder(200, `{"total": 1, "items": [{"Name": "jdoe"}]}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	result, err := client.Users.Search(context.Background(), "Name like 'jdoe'")
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Total)
}

func TestUserService_SearchByUsernameExactAndPartial(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/users",
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query().Get("query")
			switch query {
			case "Name = 'jsmith'":
				return httpmock.NewStringResponse(200, `{"total": 1, "items": [{"Name": "jsmith"}]}`), nil
			case "Name LIKE '%smith%'":
				return httpmock.NewStringResponse(200, `{"total": 2, "items": [{"Name": "jsmith"}, {"Name": "ajsmith"}]}`), nil
			default:
				return httpmock.NewStringResponse(400, "Bad Request"), nil
			}
		})

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	exact, err := client.Users.SearchByUsernameExact(context.Background(), "jsmith")
	assert.NoError(t, err)
	assert.Equal(t, 1, exact.Total)

	partial, err := client.Users.SearchByUsernamePartial(context.Background(), "smith")
	assert.NoError(t, err)
	assert.Equal(t, 2, partial.Total)
}

func TestUserService_SearchByEmailExactAndPartial(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/users",
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query().Get("query")
			switch query {
			case "EmailAddress = 'jsmith@example.com'":
				return httpmock.NewStringResponse(200, `{"total": 1, "items": [{"Name": "jsmith"}]}`), nil
			case "EmailAddress LIKE '%@example.com%'":
				return httpmock.NewStringResponse(200, `{"total": 2, "items": [{"Name": "jsmith"}, {"Name": "jsmith2"}]}`), nil
			default:
				return httpmock.NewStringResponse(400, "Bad Request"), nil
			}
		})

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	exact, err := client.Users.SearchByEmailExact(context.Background(), "jsmith@example.com")
	assert.NoError(t, err)
	assert.Equal(t, 1, exact.Total)

	partial, err := client.Users.SearchByEmailPartial(context.Background(), "@example.com")
	assert.NoError(t, err)
	assert.Equal(t, 2, partial.Total)
}

func TestUserService_SearchByNameExactAndPartial(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/users",
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query().Get("query")
			switch query {
			case "RealName = 'John Smith'":
				return httpmock.NewStringResponse(200, `{"total": 1, "items": [{"Name": "jsmith"}]}`), nil
			case "RealName LIKE '%Smith%'":
				return httpmock.NewStringResponse(200, `{"total": 2, "items": [{"Name": "jsmith"}, {"Name": "asmith"}]}`), nil
			default:
				return httpmock.NewStringResponse(400, "Bad Request"), nil
			}
		})

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	exact, err := client.Users.SearchByNameExact(context.Background(), "John Smith")
	assert.NoError(t, err)
	assert.Equal(t, 1, exact.Total)

	partial, err := client.Users.SearchByNamePartial(context.Background(), "Smith")
	assert.NoError(t, err)
	assert.Equal(t, 2, partial.Total)
}

func TestUserService_Groups(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Get Memberships
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/user/u1/groups",
		httpmock.NewStringResponder(200, `{"items": [{"id": "g1", "Name": "Admins"}]}`))
	// Add
	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/group/g1/member",
		httpmock.NewStringResponder(200, `{"message": "Added"}`))
	// Remove
	httpmock.RegisterResponder("DELETE", "http://rt.example.com/REST/2.0/group/g1/member/u1",
		httpmock.NewStringResponder(200, `{"message": "Removed"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	groups, err := client.Users.GetGroupMemberships(context.Background(), "u1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(groups))

	err = client.Users.AddToGroup(context.Background(), "u1", "g1")
	assert.NoError(t, err)

	err = client.Users.RemoveFromGroup(context.Background(), "u1", "g1")
	assert.NoError(t, err)
}

func TestUserService_History(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/user/u1/history",
		httpmock.NewStringResponder(200, `{"items": [{"Type": "Create"}]}`))
	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)
	h, err := client.Users.GetHistory(context.Background(), "u1")
	assert.NoError(t, err)
	assert.NotEmpty(t, h)
}
