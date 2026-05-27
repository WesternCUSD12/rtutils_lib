package rtutils_lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"sync"
)

var errUnsafeBroadAssetQuery = errors.New("unsafe broad asset query is not allowed")

const customFieldValueFallbackMaxAssets = 8

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
	result, err := s.SearchRaw(ctx, query)
	if err != nil {
		return nil, err
	}

	if err := s.Expand(ctx, result); err != nil {
		return nil, err
	}

	return result, nil
}

// SearchRaw searches for assets using AssetSQL without expanding each result.
func (s *AssetService) SearchRaw(ctx context.Context, query string) (*SearchResult[Asset], error) {
	if isBroadAssetQuery(query) {
		return nil, errUnsafeBroadAssetQuery
	}
	path := fmt.Sprintf("/assets?query=%s", url.QueryEscape(query))
	var result SearchResult[Asset]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	result.Finalize()
	return &result, nil
}

// Expand fetches full details for each asset in the search result.
func (s *AssetService) Expand(ctx context.Context, result *SearchResult[Asset]) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Limit concurrency to 10
	errChan := make(chan error, len(result.Items))

	for i := range result.Items {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			var fullAsset Asset
			// Using the _url provided in the search result ensures we get the exact resource
			err := s.client.request(ctx, "GET", result.Items[i].URL, nil, &fullAsset)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch details for asset %s: %w", result.Items[i].ID, err)
				return
			}
			result.Items[i] = fullAsset
		}(i)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// SearchByNameExact searches for assets with an exact name match.
func (s *AssetService) SearchByNameExact(ctx context.Context, name string) (*SearchResult[Asset], error) {
	if err := requireNonEmpty("name", name); err != nil {
		return nil, err
	}
	criteria := []map[string]interface{}{
		{
			"field":    "Name",
			"operator": "=",
			"value":    name,
		},
	}
	return s.SearchWithCriteria(ctx, criteria)
}

// SearchByNamePartial searches for assets with a partial name match.
func (s *AssetService) SearchByNamePartial(ctx context.Context, query string) (*SearchResult[Asset], error) {
	if err := requireNonEmpty("query", query); err != nil {
		return nil, err
	}
	criteria := []map[string]interface{}{
		{
			"field":    "Name",
			"operator": "LIKE",
			"value":    query,
		},
	}
	return s.SearchWithCriteria(ctx, criteria)
}

// SearchByCustomFieldExact searches for assets with an exact custom field match.
func (s *AssetService) SearchByCustomFieldExact(ctx context.Context, fieldName string, value string) (*SearchResult[Asset], error) {
	if err := requireNonEmpty("fieldName", fieldName); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("value", value); err != nil {
		return nil, err
	}
	criteria := []map[string]interface{}{
		{
			"field":    fmt.Sprintf("CustomField.{%s}", fieldName),
			"operator": "=",
			"value":    value,
		},
	}
	result, err := s.SearchWithCriteria(ctx, criteria)
	if err != nil {
		return nil, mapCustomFieldNotFound(fieldName, err)
	}
	return result, nil
}

// SearchByCustomFieldPartial searches for assets with a partial custom field match.
func (s *AssetService) SearchByCustomFieldPartial(ctx context.Context, fieldName string, query string) (*SearchResult[Asset], error) {
	result, err := s.SearchByCustomFieldPartialRaw(ctx, fieldName, query)
	if err != nil {
		return nil, err
	}
	if err := s.Expand(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

// SearchByCustomFieldPartialRaw searches for assets with a partial custom field match without expansion.
func (s *AssetService) SearchByCustomFieldPartialRaw(ctx context.Context, fieldName string, query string) (*SearchResult[Asset], error) {
	if err := requireNonEmpty("fieldName", fieldName); err != nil {
		return nil, err
	}
	if err := requireNonEmpty("query", query); err != nil {
		return nil, err
	}
	pattern := query
	if !strings.Contains(pattern, "%") {
		pattern = "%" + pattern + "%"
	}
	assetSQL := fmt.Sprintf("CustomField.{%s} LIKE '%s'", strings.TrimSpace(fieldName), escapeAssetSQLLiteral(pattern))
	result, err := s.SearchRaw(ctx, assetSQL)
	if err != nil {
		return nil, mapCustomFieldNotFound(fieldName, err)
	}
	return result, nil
}

// ListCustomFieldValues returns unique values currently present for a custom field across assets.
// Values are deduplicated case-insensitively and returned in ascending lexical order.
func (s *AssetService) ListCustomFieldValues(ctx context.Context, fieldName string) ([]string, error) {
	if err := requireNonEmpty("fieldName", fieldName); err != nil {
		return nil, err
	}

	fieldID, err := s.customFieldIDByName(ctx, fieldName)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/customfield/%s/values", url.PathEscape(fieldID))
	values := make([]string, 0, 32)
	seen := make(map[string]struct{})

	for path != "" {
		var result SearchResult[map[string]interface{}]
		if err := s.client.request(ctx, "GET", path, nil, &result); err != nil {
			return nil, err
		}
		result.Finalize()

		for _, item := range result.Items {
			value := strings.TrimSpace(parseCustomFieldValue(item))
			if value == "" {
				continue
			}
			key := strings.ToLower(value)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			values = append(values, value)
		}

		path = strings.TrimSpace(result.NextPage)
	}

	sort.Strings(values)
	if len(values) > 0 {
		return values, nil
	}

	// Some RT custom fields are freeform or have no configured value list.
	// Fall back to sampling a bounded number of matching assets and extract values.
	fallbackValues, err := s.listCustomFieldValuesFromAssets(ctx, fieldName, customFieldValueFallbackMaxAssets)
	if err != nil {
		return nil, err
	}
	if len(fallbackValues) > 0 {
		return fallbackValues, nil
	}

	return values, nil
}

func (s *AssetService) listCustomFieldValuesFromAssets(ctx context.Context, fieldName string, maxAssets int) ([]string, error) {
	if maxAssets <= 0 {
		return []string{}, nil
	}

	query := fmt.Sprintf("CustomField.{%s} LIKE '%%'", strings.TrimSpace(fieldName))
	path := fmt.Sprintf("/assets?query=%s", url.QueryEscape(query))
	seen := make(map[string]string)
	processed := 0

	for path != "" && processed < maxAssets {
		var result SearchResult[Asset]
		if err := s.client.request(ctx, "GET", path, nil, &result); err != nil {
			return nil, mapCustomFieldNotFound(fieldName, err)
		}
		result.Finalize()

		for _, item := range result.Items {
			if processed >= maxAssets {
				break
			}
			assetPath := strings.TrimSpace(item.URL)
			if assetPath == "" {
				id := parseIDField(item.ID)
				if id == "" {
					continue
				}
				assetPath = "/asset/" + id
			}

			var fullAsset Asset
			if err := s.client.request(ctx, "GET", assetPath, nil, &fullAsset); err != nil {
				continue
			}
			processed++

			for _, cf := range fullAsset.CustomFields {
				if !strings.EqualFold(strings.TrimSpace(cf.Name), strings.TrimSpace(fieldName)) {
					continue
				}
				for _, raw := range cf.Values {
					value := strings.TrimSpace(raw)
					if value == "" {
						continue
					}
					key := strings.ToLower(value)
					if _, ok := seen[key]; ok {
						continue
					}
					seen[key] = value
				}
			}
		}

		path = strings.TrimSpace(result.NextPage)
	}

	values := make([]string, 0, len(seen))
	for _, value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values, nil
}

// ListCustomFieldValuesMap returns unique value lists for each provided custom field name.
// The map key is the original field name passed by caller.
func (s *AssetService) ListCustomFieldValuesMap(ctx context.Context, fieldNames []string) (map[string][]string, error) {
	out := make(map[string][]string, len(fieldNames))
	for _, fieldName := range fieldNames {
		name := strings.TrimSpace(fieldName)
		if name == "" {
			continue
		}
		values, err := s.ListCustomFieldValues(ctx, name)
		if err != nil {
			return nil, err
		}
		out[fieldName] = values
	}
	return out, nil
}

// GetAllCustomFieldValues pages through all assets and aggregates the set of
// values for every custom field found. The returned map has keys equal to the
// custom field names as they appear on assets, and values are deduplicated
// (case-insensitive) and sorted.
func (s *AssetService) GetAllCustomFieldValues(ctx context.Context) (map[string][]string, error) {
	perPage := 100
	page := 1
	// map[fieldName]map[lowerValue]originalValue
	seen := make(map[string]map[string]string)

	// Use a safe asset search that returns all assets: "Name LIKE '%'".
	// Some RT servers do not return assets for a bare /assets request,
	// so include an explicit Name LIKE '%' query to ensure results.
	for {
		query := "Name LIKE '%'"
		path := fmt.Sprintf("/assets?query=%s&per_page=%d&page=%d", url.QueryEscape(query), perPage, page)
		var result SearchResult[Asset]
		if err := s.client.request(ctx, "GET", path, nil, &result); err != nil {
			return nil, err
		}
		result.Finalize()

		if len(result.Items) == 0 {
			break
		}

		for _, asset := range result.Items {
			for _, cf := range asset.CustomFields {
				name := strings.TrimSpace(cf.Name)
				if name == "" {
					continue
				}
				m, ok := seen[name]
				if !ok {
					m = make(map[string]string)
					seen[name] = m
				}
				for _, raw := range cf.Values {
					val := strings.TrimSpace(raw)
					if val == "" {
						continue
					}
					key := strings.ToLower(val)
					if _, exists := m[key]; !exists {
						m[key] = val
					}
				}
			}
		}

		// stop if we've reached the last page
		if result.Pages > 0 {
			if page >= result.Pages {
				break
			}
		} else if strings.TrimSpace(result.NextPage) == "" {
			break
		}
		page++
	}

	out := make(map[string][]string, len(seen))
	for name, m := range seen {
		vals := make([]string, 0, len(m))
		for _, v := range m {
			vals = append(vals, v)
		}
		sort.Strings(vals)
		out[name] = vals
	}

	return out, nil
}

// SearchWithCriteria searches for assets using JSON search syntax via POST.
func (s *AssetService) SearchWithCriteria(ctx context.Context, criteria []map[string]interface{}) (*SearchResult[Asset], error) {
	result, err := s.SearchWithCriteriaRaw(ctx, criteria)
	if err != nil {
		return nil, err
	}

	if err := s.Expand(ctx, result); err != nil {
		return nil, err
	}

	return result, nil
}

// SearchWithCriteriaRaw searches for assets using JSON search syntax via POST without expansion.
func (s *AssetService) SearchWithCriteriaRaw(ctx context.Context, criteria []map[string]interface{}) (*SearchResult[Asset], error) {
	path := "/assets"
	var result SearchResult[Asset]
	err := s.client.request(ctx, "POST", path, criteria, &result)
	if err != nil {
		return nil, err
	}
	result.Finalize()
	return &result, nil
}

// Update updates an asset.
func (s *AssetService) Update(ctx context.Context, id string, asset *Asset) error {
	path := "/asset/" + id
	var result ActionResult
	err := s.client.request(ctx, "PUT", path, asset, &result)
	if isActionResultArrayDecodeError(err) {
		return nil
	}
	return err
}

func isActionResultArrayDecodeError(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "cannot unmarshal array into Go value of type rtutils_lib.ActionResult")
}

// UpdateAsset updates an asset and returns the refreshed asset payload.
func (s *AssetService) UpdateAsset(ctx context.Context, id string, asset *Asset) (*Asset, error) {
	if err := s.Update(ctx, id, asset); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Delete deletes an asset.
func (s *AssetService) Delete(ctx context.Context, id string) error {
	path := "/asset/" + id
	var result ActionResult
	err := s.client.request(ctx, "DELETE", path, nil, &result)
	return err
}

// VerifyAssetCreated verifies that an asset with the given ID exists in RT.
// Returns the asset if found, or an error if not found or verification failed.
// This can be used as a read-back verification after asset creation.
func (s *AssetService) VerifyAssetCreated(ctx context.Context, id string) (*Asset, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("asset id cannot be empty")
	}
	asset, err := s.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to verify asset %q: %w", id, err)
	}
	if asset == nil || strings.TrimSpace(parseIDField(asset.ID)) == "" {
		return nil, fmt.Errorf("asset %q not found", id)
	}
	return asset, nil
}

// SnapshotCustomFieldValues returns a frozen snapshot of custom field values at call time.
// This is useful for batch operations where you want consistent field value options
// across multiple asset creations. The snapshot map key is the field name, value is the list.
func (s *AssetService) SnapshotCustomFieldValues(ctx context.Context, fieldNames []string) (map[string][]string, error) {
	if len(fieldNames) == 0 {
		return make(map[string][]string), nil
	}

	snapshot := make(map[string][]string, len(fieldNames))
	for _, fieldName := range fieldNames {
		name := strings.TrimSpace(fieldName)
		if name == "" {
			continue
		}
		values, err := s.ListCustomFieldValues(ctx, name)
		if err != nil {
			// Log the error but continue to populate snapshot with other fields
			// This allows batch operations to proceed even if some field snapshots fail
			snapshot[fieldName] = []string{}
			continue
		}
		snapshot[fieldName] = values
	}
	return snapshot, nil
}

// BatchCreateAssets creates multiple assets in sequence with optional error handling.
// If any creation fails and stopOnError is true, returns immediately with the error.
// Otherwise, continues with remaining assets.
// Returns a slice of created asset IDs and any accumulated errors.
func (s *AssetService) BatchCreateAssets(ctx context.Context, assets []*Asset, stopOnError bool) ([]string, []error) {
	var createdIDs []string
	var errs []error

	for _, asset := range assets {
		if asset == nil {
			errs = append(errs, errors.New("nil asset in batch"))
			if stopOnError {
				break
			}
			continue
		}

		id, err := s.Create(ctx, asset)
		if err != nil {
			errs = append(errs, err)
			if stopOnError {
				break
			}
			continue
		}

		createdIDs = append(createdIDs, id)
	}

	return createdIDs, errs
}

func mapCustomFieldNotFound(fieldName string, err error) error {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		message := strings.ToLower(apiErr.Message)
		if strings.Contains(message, "custom field") && strings.Contains(message, "not found") {
			return &CustomFieldNotFoundError{FieldName: fieldName}
		}
	}
	return err
}

// Asset represents a Request Tracker asset.
type Asset struct {
	ID           json.Number        `json:"id,omitempty"`
	URL          string             `json:"_url,omitempty"`
	Name         string             `json:"Name,omitempty"`
	Catalog      string             `json:"-"`
	Content      string             `json:"Content,omitempty"`
	Status       string             `json:"Status,omitempty"`
	Owner        string             `json:"-"`
	HeldBy       string             `json:"-"`
	LastUpdated  string             `json:"-"`
	CustomFields []AssetCustomField `json:"CustomFields,omitempty"`
}

// MarshalJSON implements custom marshaling for Asset so that Catalog is included
// in the JSON payload when set (the field tag is json:"-" to handle RT's polymorphic
// response format, but we need it serialized on creation/update).
// CustomFields is serialized as a map[string]string (RT write format: {"FieldName": "value"})
// rather than the array-of-structs format returned by RT on reads.
func (a *Asset) MarshalJSON() ([]byte, error) {
	type AssetAlias Asset
	// Build CustomFields as a map for the RT write API: {"FieldName": "value"}
	cfMap := make(map[string]string)
	for _, cf := range a.CustomFields {
		if len(cf.Values) > 0 {
			cfMap[cf.Name] = cf.Values[0]
		}
	}
	out := struct {
		*AssetAlias
		Catalog      string            `json:"Catalog,omitempty"`
		Owner        string            `json:"Owner,omitempty"`
		HeldBy       string            `json:"HeldBy,omitempty"`
		CustomFields map[string]string `json:"CustomFields,omitempty"`
	}{
		AssetAlias:   (*AssetAlias)(a),
		Catalog:      a.Catalog,
		Owner:        a.Owner,
		HeldBy:       a.HeldBy,
		CustomFields: cfMap,
	}
	return json.Marshal(out)
}

// CatalogInfo represents a Request Tracker asset catalog.
type CatalogInfo struct {
	ID          string `json:"id"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
}

// UnmarshalJSON handles custom unmarshaling for CatalogInfo to handle numeric IDs.
func (c *CatalogInfo) UnmarshalJSON(data []byte) error {
	type Alias CatalogInfo
	aux := &struct {
		ID interface{} `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.ID = parseIDField(aux.ID)
	return nil
}

// ListCatalogs returns all available asset catalogs.
func (s *AssetService) ListCatalogs(ctx context.Context) ([]CatalogInfo, error) {
	var result SearchResult[CatalogInfo]
	if err := s.client.request(ctx, "GET", "/catalogs/all", nil, &result); err != nil {
		return nil, err
	}
	result.Finalize()
	// If the returned catalog items do not include names, attempt to expand
	// each catalog by fetching its details.
	needsExpansion := false
	for _, c := range result.Items {
		if strings.TrimSpace(c.Name) == "" {
			needsExpansion = true
			break
		}
	}
	if needsExpansion && len(result.Items) > 0 {
		log.Printf("DEBUG: Catalogs: Expanding %d catalog results...", len(result.Items))
		for i := range result.Items {
			var full CatalogInfo
			if err := s.client.request(ctx, "GET", "/catalog/"+result.Items[i].ID, nil, &full); err != nil {
				log.Printf("DEBUG: Catalogs: Failed to expand catalog %s: %v", result.Items[i].ID, err)
				continue
			}
			result.Items[i] = full
		}
	}

	return result.Items, nil
}

// UnmarshalJSON handles custom unmarshaling for Asset to support RT's format
func (a *Asset) UnmarshalJSON(data []byte) error {
	type AssetAlias Asset
	aux := struct {
		Catalog     interface{} `json:"Catalog"`
		Owner       interface{} `json:"Owner"`
		HeldBy      interface{} `json:"HeldBy"`
		LastUpdated string      `json:"LastUpdated"`
		*AssetAlias
	}{
		AssetAlias: (*AssetAlias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	a.Catalog = parseStringOrObject(aux.Catalog)
	a.Owner = parseStringOrObject(aux.Owner)
	a.HeldBy = parseStringOrObject(aux.HeldBy)
	a.LastUpdated = aux.LastUpdated

	return nil
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

func uniqueCustomFieldValues(assets []Asset, fieldName string) []string {
	seen := make(map[string]string)
	for _, asset := range assets {
		for _, cf := range asset.CustomFields {
			if !strings.EqualFold(strings.TrimSpace(cf.Name), strings.TrimSpace(fieldName)) {
				continue
			}
			for _, value := range cf.Values {
				trimmed := strings.TrimSpace(value)
				if trimmed == "" {
					continue
				}
				key := strings.ToLower(trimmed)
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = trimmed
			}
		}
	}

	values := make([]string, 0, len(seen))
	for _, value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}

func (s *AssetService) customFieldIDByName(ctx context.Context, fieldName string) (string, error) {
	criteria := []map[string]interface{}{
		{
			"field":    "Name",
			"operator": "=",
			"value":    strings.TrimSpace(fieldName),
		},
	}

	var result SearchResult[map[string]interface{}]
	if err := s.client.request(ctx, "POST", "/customfields", criteria, &result); err != nil {
		return "", err
	}
	result.Finalize()

	for _, item := range result.Items {
		id := parseIDField(item["id"])
		if id != "" {
			return id, nil
		}
	}

	return "", &CustomFieldNotFoundError{FieldName: fieldName}
}

func parseCustomFieldValue(item map[string]interface{}) string {
	if item == nil {
		return ""
	}
	if value := parseFieldString(item, "name", "Name", "value", "Value", "Content"); value != "" {
		return value
	}
	return parseIDField(item["id"])
}

func parseFieldString(item map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if raw, ok := item[key]; ok {
			switch v := raw.(type) {
			case string:
				if strings.TrimSpace(v) != "" {
					return v
				}
			}
		}
	}
	return ""
}

func isBroadAssetQuery(query string) bool {
	normalized := strings.ToLower(strings.TrimSpace(query))
	if normalized == "" {
		return true
	}
	if normalized == "%" || normalized == "*" {
		return true
	}
	if normalized == "name like '%'" || normalized == "name like \"%\"" {
		return true
	}
	if normalized == "name like *" {
		return true
	}
	return false
}

func escapeAssetSQLLiteral(input string) string {
	return strings.ReplaceAll(input, "'", "\\'")
}
