package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"rtutils_lib"
)

func main() {
	baseURL := os.Getenv("RT_BASE_URL")
	token := os.Getenv("RT_TOKEN")
	if baseURL == "" || token == "" {
		log.Fatal("Please set RT_BASE_URL and RT_TOKEN environment variables.")
	}

	client := rtutils_lib.NewClient(baseURL, token)
	ctx := context.Background()

	fieldName := "Model" // Change this to any custom field name you want to list values for
	values, err := client.Assets.ListCustomFieldValues(ctx, fieldName)
	if err != nil {
		log.Fatalf("Error retrieving custom field values for '%s': %v", fieldName, err)
	}

	fmt.Printf("Available values for custom field '%s':\n", fieldName)
	for _, v := range values {
		fmt.Println(v)
	}
}
