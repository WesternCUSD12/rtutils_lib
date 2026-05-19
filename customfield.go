package rtutils_lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
)

// CustomFieldService handles communication with the custom field related methods
// of the Request Tracker API.
type CustomFieldService struct {
	client *Client
}

// CustomField represents a Request Tracker custom field definition.
type CustomField struct {
	ID          string   `json:"id"`
	URL         string   `json:"_url,omitempty"`
	Name        string   `json:"Name"`
	Type        string   `json:"Type,omitempty"`
	Description string   `json:"Description,omitempty"`
	Disabled    string   `json:"Disabled,omitempty"`
	MaxValues   string   `json:"MaxValues,omitempty"`
	Pattern     string   `json:"Pattern,omitempty"`
	Values      []string `json:"Values,omitempty"` // populated for Select-type fields
}

// UnmarshalJSON handles RT's inconsistent id type (string or number).
func (cf *CustomField) UnmarshalJSON(data []byte) error {
	type Alias CustomField
	aux := &struct {
		ID interface{} `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(cf),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	cf.ID = parseIDField(aux.ID)
	return nil
}

// Get retrieves a custom field by ID.
func (s *CustomFieldService) Get(ctx context.Context, id string) (*CustomField, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return nil, err
	}
	var cf CustomField
	if err := s.client.request(ctx, "GET", "/customfield/"+id, nil, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// GetByCategory retrieves a custom field by ID with values filtered by category
// (for Select-type fields).
func (s *CustomFieldService) GetByCategory(ctx context.Context, id, category string) (*CustomField, error) {
	if err := requireNonEmpty("id", id); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("category", category); err != nil {
		return nil, err
	}
	path := "/customfield/" + id + "?category=" + url.QueryEscape(category)
	var cf CustomField
	if err := s.client.request(ctx, "GET", path, nil, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// expand fetches full details for each custom field stub (ID-only items) in place.
func (s *CustomFieldService) expand(ctx context.Context, result *SearchResult[CustomField]) {
	needsExpansion := false
	for _, cf := range result.Items {
		if strings.TrimSpace(cf.Name) == "" {
			needsExpansion = true
			break
		}
	}
	if !needsExpansion || len(result.Items) == 0 {
		return
	}

	log.Printf("DEBUG: CustomFields: Expanding %d results...", len(result.Items))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)
	for i := range result.Items {
		if strings.TrimSpace(result.Items[i].Name) != "" {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			var full CustomField
			fetchPath := "/customfield/" + result.Items[i].ID
			if result.Items[i].URL != "" {
				fetchPath = result.Items[i].URL
			}
			if err := s.client.request(ctx, "GET", fetchPath, nil, &full); err != nil {
				log.Printf("DEBUG: CustomFields: Failed to expand field %s: %v", result.Items[i].ID, err)
				return
			}
			result.Items[i] = full
		}(i)
	}
	wg.Wait()
}

// List searches for all custom fields using optional JSON criteria.
// Pass nil criteria to list all custom fields.
// Results are automatically expanded when the search endpoint returns stubs.
func (s *CustomFieldService) List(ctx context.Context, criteria []map[string]interface{}) (*SearchResult[CustomField], error) {
	var result SearchResult[CustomField]
	if criteria == nil {
		if err := s.client.request(ctx, "GET", "/customfields", nil, &result); err != nil {
			return nil, err
		}
	} else {
		if err := s.client.request(ctx, "POST", "/customfields", criteria, &result); err != nil {
			return nil, err
		}
	}
	result.Finalize()
	s.expand(ctx, &result)
	return &result, nil
}

// ListAll retrieves every custom field defined in the RT instance.
func (s *CustomFieldService) ListAll(ctx context.Context) ([]CustomField, error) {
	result, err := s.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// ListForCatalog returns all custom fields attached to the given catalog ID.
// Results are automatically expanded when the endpoint returns stubs.
func (s *CustomFieldService) ListForCatalog(ctx context.Context, catalogID string, criteria []map[string]interface{}) (*SearchResult[CustomField], error) {
	if err := requireNonEmpty("catalogID", catalogID); err != nil {
		return nil, err
	}
	var result SearchResult[CustomField]
	path := fmt.Sprintf("/catalog/%s/customfields", catalogID)
	if criteria == nil {
		if err := s.client.request(ctx, "GET", path, nil, &result); err != nil {
			return nil, err
		}
	} else {
		if err := s.client.request(ctx, "POST", path, criteria, &result); err != nil {
			return nil, err
		}
	}
	result.Finalize()
	s.expand(ctx, &result)
	return &result, nil
}

// ListAllForCatalog returns every custom field attached to the given catalog ID.
func (s *CustomFieldService) ListAllForCatalog(ctx context.Context, catalogID string) ([]CustomField, error) {
	result, err := s.ListForCatalog(ctx, catalogID, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// ListForQueue returns all custom fields attached to the given queue ID.
// Results are automatically expanded when the endpoint returns stubs.
func (s *CustomFieldService) ListForQueue(ctx context.Context, queueID string, criteria []map[string]interface{}) (*SearchResult[CustomField], error) {
	if err := requireNonEmpty("queueID", queueID); err != nil {
		return nil, err
	}
	var result SearchResult[CustomField]
	path := fmt.Sprintf("/queue/%s/customfields", queueID)
	if criteria == nil {
		if err := s.client.request(ctx, "GET", path, nil, &result); err != nil {
			return nil, err
		}
	} else {
		if err := s.client.request(ctx, "POST", path, criteria, &result); err != nil {
			return nil, err
		}
	}
	result.Finalize()
	s.expand(ctx, &result)
	return &result, nil
}

// ListAllForQueue returns every custom field attached to the given queue ID.
func (s *CustomFieldService) ListAllForQueue(ctx context.Context, queueID string) ([]CustomField, error) {
	result, err := s.ListForQueue(ctx, queueID, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// ListForClass returns all custom fields attached to the given class ID.
// Results are automatically expanded when the endpoint returns stubs.
func (s *CustomFieldService) ListForClass(ctx context.Context, classID string, criteria []map[string]interface{}) (*SearchResult[CustomField], error) {
	if err := requireNonEmpty("classID", classID); err != nil {
		return nil, err
	}
	var result SearchResult[CustomField]
	path := fmt.Sprintf("/class/%s/customfields", classID)
	if criteria == nil {
		if err := s.client.request(ctx, "GET", path, nil, &result); err != nil {
			return nil, err
		}
	} else {
		if err := s.client.request(ctx, "POST", path, criteria, &result); err != nil {
			return nil, err
		}
	}
	result.Finalize()
	s.expand(ctx, &result)
	return &result, nil
}

// ListAllForClass returns every custom field attached to the given class ID.
func (s *CustomFieldService) ListAllForClass(ctx context.Context, classID string) ([]CustomField, error) {
	result, err := s.ListForClass(ctx, classID, nil)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FindByName looks up a custom field by exact name from the full list.
// Returns CustomFieldNotFoundError if no match is found.
func (s *CustomFieldService) FindByName(ctx context.Context, name string) (*CustomField, error) {
	if err := requireNonEmpty("name", name); err != nil {
		return nil, err
	}
	all, err := s.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	for i, cf := range all {
		if strings.EqualFold(cf.Name, name) {
			return &all[i], nil
		}
	}
	return nil, &CustomFieldNotFoundError{FieldName: name}
}

// NamesForCatalog returns a sorted slice of all custom field names attached to a catalog.
// Useful for discovering exact field names required by the RT instance.
func (s *CustomFieldService) NamesForCatalog(ctx context.Context, catalogID string) ([]string, error) {
	fields, err := s.ListAllForCatalog(ctx, catalogID)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(fields))
	for _, cf := range fields {
		if cf.Name != "" {
			names = append(names, cf.Name)
		}
	}
	return names, nil
}
