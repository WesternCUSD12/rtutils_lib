package rtutils_lib

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// UserService handles communication with the user related methods of the
// Request Tracker API.
type UserService struct {
	client *Client
}

// Create creates a new user.
func (s *UserService) Create(ctx context.Context, user *User) (string, error) {
	var result ActionResult
	err := s.client.request(ctx, "POST", "/user", user, &result)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Get fetches a user by ID or Name.
func (s *UserService) Get(ctx context.Context, id string) (*User, error) {
	path := "/user/" + id
	var user User
	err := s.client.request(ctx, "GET", path, nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

// Search searches for users.
func (s *UserService) Search(ctx context.Context, query string) (*SearchResult[User], error) {
	path := fmt.Sprintf("/users?query=%s", url.QueryEscape(query))
	var result SearchResult[User]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchByUsernameExact searches for users by exact username.
func (s *UserService) SearchByUsernameExact(ctx context.Context, username string) (*SearchResult[User], error) {
	if err := requireNonEmpty("username", username); err != nil {
		return nil, err
	}
	return s.Search(ctx, buildUserQuery("Name", username, false))
}

// SearchByUsernamePartial searches for users by partial username.
func (s *UserService) SearchByUsernamePartial(ctx context.Context, query string) (*SearchResult[User], error) {
	if err := requireNonEmpty("query", query); err != nil {
		return nil, err
	}
	return s.Search(ctx, buildUserQuery("Name", query, true))
}

// SearchByEmailExact searches for users by exact email address.
func (s *UserService) SearchByEmailExact(ctx context.Context, email string) (*SearchResult[User], error) {
	if err := requireNonEmpty("email", email); err != nil {
		return nil, err
	}
	return s.Search(ctx, buildUserQuery("EmailAddress", email, false))
}

// SearchByEmailPartial searches for users by partial email address.
func (s *UserService) SearchByEmailPartial(ctx context.Context, query string) (*SearchResult[User], error) {
	if err := requireNonEmpty("query", query); err != nil {
		return nil, err
	}
	return s.Search(ctx, buildUserQuery("EmailAddress", query, true))
}

// SearchByNameExact searches for users by exact full name.
func (s *UserService) SearchByNameExact(ctx context.Context, name string) (*SearchResult[User], error) {
	if err := requireNonEmpty("name", name); err != nil {
		return nil, err
	}
	return s.Search(ctx, buildUserQuery("RealName", name, false))
}

// SearchByNamePartial searches for users by partial full name.
func (s *UserService) SearchByNamePartial(ctx context.Context, query string) (*SearchResult[User], error) {
	if err := requireNonEmpty("query", query); err != nil {
		return nil, err
	}
	return s.Search(ctx, buildUserQuery("RealName", query, true))
}

// Update updates a user.
func (s *UserService) Update(ctx context.Context, id string, user *User) error {
	path := "/user/" + id
	var result ActionResult
	err := s.client.request(ctx, "PUT", path, user, &result)
	return err
}

// Disable disables a user.
func (s *UserService) Disable(ctx context.Context, id string) error {
	path := "/user/" + id
	// Assuming DELETE disables. Or PUT with Disabled=1
	var result ActionResult
	err := s.client.request(ctx, "DELETE", path, nil, &result)
	return err
}

// GetHistory fetches user history.
func (s *UserService) GetHistory(ctx context.Context, id string) ([]Transaction, error) {
	path := "/user/" + id + "/history"
	var result SearchResult[Transaction]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// GetGroupMemberships fetches groups the user belongs to.
func (s *UserService) GetGroupMemberships(ctx context.Context, id string) ([]string, error) {
	path := "/user/" + id + "/groups"
	// Response might be a list of groups.
	// Structure: SearchResult[Group]? But we don't have Group struct.
	// We can define a temporary struct or just use map[string]interface{}
	// Interface says returns []string (IDs or Names?)
	// Task T026: GetGroupMemberships
	// Let's assume it returns standard SearchResult of objects with ID/Name
	type groupStub struct {
		ID   string `json:"id"`
		Name string `json:"Name"`
	}
	var result SearchResult[groupStub]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(result.Items))
	for i, g := range result.Items {
		names[i] = g.Name
	}
	return names, nil
}

// AddToGroup adds a user to a group.
func (s *UserService) AddToGroup(ctx context.Context, userID, groupID string) error {
	path := "/group/" + groupID + "/member"
	payload := map[string]string{"id": userID}
	var result ActionResult
	err := s.client.request(ctx, "POST", path, payload, &result)
	return err
}

// RemoveFromGroup removes a user from a group.
func (s *UserService) RemoveFromGroup(ctx context.Context, userID, groupID string) error {
	path := "/group/" + groupID + "/member/" + userID
	var result ActionResult
	err := s.client.request(ctx, "DELETE", path, nil, &result)
	return err
}

// User represents a Request Tracker user.
type User struct {
	ID           string `json:"id,omitempty"`
	URL          string `json:"_url,omitempty"`
	Name         string `json:"Name,omitempty"`
	RealName     string `json:"RealName,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
	Disabled     int    `json:"Disabled,omitempty"`
	Privileged   int    `json:"Privileged,omitempty"`
}

func buildUserQuery(field string, value string, partial bool) string {
	escaped := strings.ReplaceAll(value, "'", "''")
	if partial {
		return fmt.Sprintf("%s LIKE '%%%s%%'", field, escaped)
	}
	return fmt.Sprintf("%s = '%s'", field, escaped)
}
