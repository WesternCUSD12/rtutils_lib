package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"

	"log"
	"os"
	"rtutils_lib"
	"strings"
	"text/tabwriter"
)

func loadEnv() {
	// Try loading from .env in current directory
	file, err := os.Open(".env")
	if err != nil {
		// .env file is optional
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Only set if not already set (allow override from real env)
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Warning: error reading .env file: %v", err)
	}
}

func printAssetDetails(asset *rtutils_lib.Asset) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Core Fields
	fmt.Fprintf(w, "ID:\t%s\n", asset.ID)
	fmt.Fprintf(w, "Name:\t%s\n", asset.Name)
	fmt.Fprintf(w, "Status:\t%s\n", asset.Status)

	// Specific Custom Fields (Priority)
	typeCF := asset.GetCustomField("Type")
	if typeCF != "" {
		fmt.Fprintf(w, "Type:\t%s\n", typeCF)
	}
	modelCF := asset.GetCustomField("Model")
	if modelCF != "" {
		fmt.Fprintf(w, "Model:\t%s\n", modelCF)
	}
	manufCF := asset.GetCustomField("Manufacturer")
	if manufCF != "" {
		fmt.Fprintf(w, "Manufacturer:\t%s\n", manufCF)
	}

	// Other Custom Fields
	// Map to track what we've already printed
	printedFields := map[string]bool{
		"Type":         true,
		"Model":        true,
		"Manufacturer": true,
	}

	for _, cf := range asset.CustomFields {
		if _, seen := printedFields[cf.Name]; !seen {
			if len(cf.Values) > 0 && cf.Values[0] != "" {
				fmt.Fprintf(w, "%s:\t%s\n", cf.Name, cf.Values[0])
			}
		}
	}
	w.Flush()
}

func runIDSearch(client *rtutils_lib.Client, id string) {
	ctx := context.Background()
	asset, err := client.Assets.Get(ctx, id)
	if err != nil {
		var apiErr *rtutils_lib.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			fmt.Fprintf(os.Stderr, "Asset ID %s not found.\n", id)
			os.Exit(2)
		}
		fmt.Fprintf(os.Stderr, "Error fetching asset %s: %v\n", id, err)
		os.Exit(1)
	}

	printAssetDetails(asset)
	os.Exit(0)
}

func printAssetSummary(assets []rtutils_lib.Asset) {
	fmt.Fprintf(os.Stderr, "Multiple assets found. Please refine your search.\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tURL")
	for _, asset := range assets {
		fmt.Fprintf(w, "%s\t%s\t%s\n", asset.ID, asset.Name, asset.URL)
	}
	w.Flush()
}

func handleSearchResults(results *rtutils_lib.SearchResult[rtutils_lib.Asset]) {
	if len(results.Items) == 0 {
		fmt.Fprintln(os.Stderr, "Asset not found.")
		os.Exit(2)
	}

	if len(results.Items) == 1 {
		printAssetDetails(&results.Items[0])
		os.Exit(0)
	}

	// Multiple matches
	printAssetSummary(results.Items)
	os.Exit(1) // Ambiguous match
}

func runNameSearch(client *rtutils_lib.Client, name string) {
	ctx := context.Background()
	// AssetSQL: Name = 'VALUE'
	query := fmt.Sprintf("Name = '%s'", name)
	results, err := client.Assets.Search(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching asset by name %s: %v\n", name, err)
		os.Exit(1)
	}
	handleSearchResults(results)
}

func runInternalNameSearch(client *rtutils_lib.Client, internalName string) {
	// Construct AssetSQL query for CF.{Internal Name}
	// Note: We need to handle potential quoting if name has single quotes, but for now assuming simple alphanumeric.
	// Spec defines Internal Name as a Custom Field.
	query := fmt.Sprintf("'CF.{Internal Name}' = '%s'", internalName)

	ctx := context.Background()
	results, err := client.Assets.Search(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching asset by internal name %s: %v\n", internalName, err)
		os.Exit(1)
	}
	handleSearchResults(results)
}

func main() {
	loadEnv()

	var id string
	var name string
	var internalName string

	flag.StringVar(&id, "id", "", "Search by Asset ID")
	flag.StringVar(&name, "name", "", "Search by Asset Name")
	flag.StringVar(&internalName, "internal-name", "", "Search by Internal Name (Custom Field)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Validation: Exactly one flag
	count := 0
	if id != "" {
		count++
	}
	if name != "" {
		count++
	}
	if internalName != "" {
		count++
	}

	if count == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if count > 1 {
		fmt.Fprintln(os.Stderr, "Error: Flags --id, --name, and --internal-name are mutually exclusive.")
		flag.Usage()
		os.Exit(1)
	}

	// Environment Variables
	baseURL := os.Getenv("RT_BASE_URL")
	token := os.Getenv("RT_TOKEN")

	if baseURL == "" || token == "" {
		fmt.Fprintln(os.Stderr, "Error: RT_BASE_URL and RT_TOKEN environment variables are required.")
		os.Exit(1)
	}

	client := rtutils_lib.NewClient(baseURL, token)
	_ = client // Suppress unused var error for now

	// Determine action based on flag (Placeholder for Phase 3/4/5)
	if id != "" {
		runIDSearch(client, id)
	} else if name != "" {
		runNameSearch(client, name)
	} else if internalName != "" {
		runInternalNameSearch(client, internalName)
	}
}
