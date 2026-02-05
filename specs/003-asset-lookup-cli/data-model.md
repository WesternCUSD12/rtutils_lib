# Data Model: Asset Lookup CLI

**Feature**: `003-asset-lookup-cli`

## Overview

This tool does not introduce new database schemas. However, it defines a specific **View Model** for displaying Asset details to the user.

## Asset Detail View

The CLI will map the raw `Asset` struct to a display-optimized view.

| Display Label | Source Field | Type | Notes |
|---------------|--------------|------|-------|
| **ID** | `Asset.ID` | String | |
| **Name** | `Asset.Name` | String | |
| **Status** | `Asset.Status` | String | |
| **Type** | `Asset.CustomFields["Type"]` | String | Fetched via helper |
| **Model** | `Asset.CustomFields["Model"]` | String | Fetched via helper |
| **Manufacturer** | `Asset.CustomFields["Manufacturer"]` | String | Fetched via helper |
| *[Other CF Names]* | `Asset.CustomFields` | String | Dynamically listed |

## Logic

1. **Hydration**: The standard `AssetService.Search` returns summary objects. The tool must iterate results and fetch full details (Search Hydration Pattern) if `Search` doesn't return everything (library handles this now). _Wait, strict `Search` might trigger hydration overhead. For individual lookup by ID, use `Get`. For Name search, use `Search`, then if 1 result, use `Get` (or rely on library hydration if efficient)._
2. **Ambiguity**: If `len(results) > 1`, the view switches to "Summary List".

## Summary List View

| Column | Source Field |
|--------|--------------|
| **ID** | `Asset.ID` |
| **Name** | `Asset.Name` |
| **URL** | `Asset.URL` |
