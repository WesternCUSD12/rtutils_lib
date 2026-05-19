//go:build ignore
// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"rtutils_lib"
)

func main() {
	base := os.Getenv("RT_BASE_URL")
	token := os.Getenv("RT_TOKEN")
	if base == "" || token == "" {
		log.Fatal("Please set RT_BASE_URL and RT_TOKEN in the environment")
	}

	client := rtutils_lib.NewClient(base, token)
	ctx := context.Background()

	m, err := client.Assets.GetAllCustomFieldValues(ctx)
	if err != nil {
		log.Fatalf("failed to collect custom field values: %v", err)
	}

	outPath := "/tmp/rt_all_custom_field_values.json"
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal result: %v", err)
	}
	if err := os.WriteFile(outPath, b, 0644); err != nil {
		log.Fatalf("failed to write %s: %v", outPath, err)
	}

	fmt.Printf("Wrote %d fields to %s\n", len(m), outPath)
	// Print a small sample
	printed := 0
	for k, vals := range m {
		fmt.Printf("- %s: %d values\n", k, len(vals))
		printed++
		if printed >= 10 {
			break
		}
	}
}
