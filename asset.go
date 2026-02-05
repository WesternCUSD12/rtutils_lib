package rtutils_lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// AssetService handles communication with the asset related methods of the
// Request Tracker API.
type AssetService struct {
	client *Client
}

// Create creates a new asset.
func (s *AssetService) Create(ctx context.Context, asset *Asset) (string, error) {
	var result ActionResult
	err := s.client.request(ctx, "POST", "/asset", asset, &result)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Get fetches an asset by ID.
func (s *AssetService) Get(ctx context.Context, id string) (*Asset, error) {
	path := "/asset/" + id
	var asset Asset
	err := s.client.request(ctx, "GET", path, nil, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// Search searches for assets using AssetSQL.
func (s *AssetService) Search(ctx context.Context, query string) (*SearchResult[Asset], error) {
	path := fmt.Sprintf("/assets?query=%s", url.QueryEscape(query))
	var result SearchResult[Asset]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	// RT Search results often only contain summary data (ID, URL).
	// We need to fetch the full details for each asset using its _url.
	for i := range result.Items {
		var fullAsset Asset
		// Using the _url provided in the search result ensures we get the exact resource
		err := s.client.request(ctx, "GET", result.Items[i].URL, nil, &fullAsset)
		if err != nil {
			// If we fail to fetch details, we log/return error or just continue?
			// Returning error seems safer for data integrity.
			return nil, fmt.Errorf("failed to fetch details for asset %s: %w", result.Items[i].ID, err)
		}
		result.Items[i] = fullAsset
	}

	return &result, nil
}

// Update updates an asset.
func (s *AssetService) Update(ctx context.Context, id string, asset *Asset) error {
	path := "/asset/" + id
	var result ActionResult
	err := s.client.request(ctx, "PUT", path, asset, &result)
	return err
}

// Delete deletes an asset.
func (s *AssetService) Delete(ctx context.Context, id string) error {
	path := "/asset/" + id
	var result ActionResult
	err := s.client.request(ctx, "DELETE", path, nil, &result)
	return err
}

// Asset represents a Request Tracker asset.
type Asset struct {
	ID           json.Number        `json:"id,omitempty"`
	URL          string             `json:"_url,omitempty"`
	Name         string             `json:"Name,omitempty"`
	Catalog      interface{}        `json:"Catalog,omitempty"`
	Content      string             `json:"Content,omitempty"`
	Status       string             `json:"Status,omitempty"`
	CustomFields []AssetCustomField `json:"CustomFields,omitempty"`
}

// GetCustomField returns the first value of a custom field by name, or an empty string if not found.
func (a *Asset) GetCustomField(name string) string {
	for _, cf := range a.CustomFields {
		if cf.Name == name {
			if len(cf.Values) > 0 {
				return cf.Values[0]
			}
			return ""
		}
	}
	return ""
}

type AssetCustomField struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []string `json:"values"`
}
