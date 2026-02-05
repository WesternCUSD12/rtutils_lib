package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"rtutils_lib"
	"text/tabwriter"
)

func main() {
	// Flags
	var username string
	var email string
	var realName string

	flag.StringVar(&username, "username", "", "Search by User Login Name")
	flag.StringVar(&email, "email", "", "Search by Email Address")
	flag.StringVar(&realName, "name", "", "Search by Real Name")
	flag.Parse()

	// Validation
	count := 0
	if username != "" {
		count++
	}
	if email != "" {
		count++
	}
	if realName != "" {
		count++
	}

	if count != 1 {
		log.Fatal("Exactly one search flag (--username, --email, --name) must be provided")
	}

	// Environment Variables
	baseURL := os.Getenv("RT_BASE_URL")
	token := os.Getenv("RT_TOKEN")

	if baseURL == "" || token == "" {
		log.Fatal("RT_BASE_URL and RT_TOKEN environment variables are required")
	}

	// Query Builder
	var query string
	if username != "" {
		query = fmt.Sprintf("Owner.Name = '%s' OR HeldBy.Name = '%s'", username, username)
	} else if email != "" {
		query = fmt.Sprintf("Owner.EmailAddress = '%s' OR HeldBy.EmailAddress = '%s'", email, email)
	} else if realName != "" {
		query = fmt.Sprintf("Owner.RealName = '%s' OR HeldBy.RealName = '%s'", realName, realName)
	}

	client := rtutils_lib.NewClient(baseURL, token)

	// Search Assets
	results, err := client.Assets.Search(context.Background(), query)
	if err != nil {
		log.Fatalf("Error searching assets: %v", err)
	}

	if len(results.Items) == 0 {
		fmt.Println("No assets found matching the criteria.")
		return
	}

	// Output Table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tType\tModel\tManufacturer")
	for _, asset := range results.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			asset.ID,
			asset.Name,
			asset.GetCustomField("Type"),
			asset.GetCustomField("Model"),
			asset.GetCustomField("Manufacturer"),
		)
	}
	w.Flush()
}
