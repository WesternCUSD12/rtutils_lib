package rtutils_lib

import (
	"context"
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
