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

type searchFilters struct {
	Queue     string
	Status    string
	Owner     string
	Requestor string
	Subject   string
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Warning: error reading .env file: %v", err)
	}
}

func normalizePagination(page, perPage int) (int, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		return 0, 0, fmt.Errorf("per-page cannot exceed 100")
	}
	return page, perPage, nil
}

func escapeTicketSQL(value string) string {
	return strings.ReplaceAll(value, "'", "\\'")
}

func buildTicketSQL(filters searchFilters) (string, error) {
	clauses := make([]string, 0, 5)

	if filters.Queue != "" {
		clauses = append(clauses, fmt.Sprintf("Queue = '%s'", escapeTicketSQL(filters.Queue)))
	}
	if filters.Status != "" {
		clauses = append(clauses, fmt.Sprintf("Status = '%s'", escapeTicketSQL(filters.Status)))
	}
	if filters.Owner != "" {
		clauses = append(clauses, fmt.Sprintf("Owner = '%s'", escapeTicketSQL(filters.Owner)))
	}
	if filters.Requestor != "" {
		clauses = append(clauses, fmt.Sprintf("Requestor = '%s'", escapeTicketSQL(filters.Requestor)))
	}
	if filters.Subject != "" {
		subject := escapeTicketSQL(filters.Subject)
		clauses = append(clauses, fmt.Sprintf("Subject LIKE '%%%s%%'", subject))
	}

	if len(clauses) == 0 {
		return "", errors.New("at least one filter must be provided")
	}

	return strings.Join(clauses, " AND "), nil
}

func printSearchResults(results *rtutils_lib.SearchResult[rtutils_lib.Ticket]) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSubject\tStatus\tQueue\tOwner\tRequestor\tURL")
	for _, ticket := range results.Items {
		requestor := ""
		if len(ticket.Requestor) > 0 {
			requestor = ticket.Requestor[0]
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", ticket.ID, ticket.Subject, ticket.Status, ticket.Queue, ticket.Owner, requestor, ticket.URL)
	}
	w.Flush()
}

func printPaginationSummary(results *rtutils_lib.SearchResult[rtutils_lib.Ticket], page int, perPage int) {
	pages := results.Pages
	if pages == 0 && perPage > 0 {
		pages = (results.Total + perPage - 1) / perPage
	}

	fmt.Fprintf(os.Stdout, "\nPage %d of %d (per-page %d, total %d)\n", page, pages, perPage, results.Total)
	if results.NextPage != "" {
		fmt.Fprintf(os.Stdout, "Next page: %s\n", results.NextPage)
	}
}

func printTicketDetails(ticket *rtutils_lib.Ticket) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", ticket.ID)
	fmt.Fprintf(w, "Subject:\t%s\n", ticket.Subject)
	fmt.Fprintf(w, "Status:\t%s\n", ticket.Status)
	fmt.Fprintf(w, "Queue:\t%s\n", ticket.Queue)
	fmt.Fprintf(w, "Owner:\t%s\n", ticket.Owner)
	if len(ticket.Requestor) > 0 {
		fmt.Fprintf(w, "Requestor:\t%s\n", strings.Join(ticket.Requestor, ", "))
	}
	if ticket.URL != "" {
		fmt.Fprintf(w, "URL:\t%s\n", ticket.URL)
	}
	if ticket.Created != "" {
		fmt.Fprintf(w, "Created:\t%s\n", ticket.Created)
	}
	if ticket.Resolved != "" {
		fmt.Fprintf(w, "Resolved:\t%s\n", ticket.Resolved)
	}
	w.Flush()
}

func findTicketByID(results *rtutils_lib.SearchResult[rtutils_lib.Ticket], id string) (*rtutils_lib.Ticket, bool) {
	for _, ticket := range results.Items {
		if ticket.ID == id {
			return &ticket, true
		}
	}
	return nil, false
}

func main() {
	loadEnv()

	var queue string
	var status string
	var owner string
	var requestor string
	var subject string
	var page int
	var perPage int
	var detailsID string

	flag.StringVar(&queue, "queue", "", "Exact match for Queue")
	flag.StringVar(&status, "status", "", "Exact match for Status")
	flag.StringVar(&owner, "owner", "", "Exact match for Owner")
	flag.StringVar(&requestor, "requestor", "", "Exact match for Requestor")
	flag.StringVar(&subject, "subject", "", "Keyword match for Subject")
	flag.IntVar(&page, "page", 0, "Page number (default 1)")
	flag.IntVar(&perPage, "per-page", 0, "Results per page (default 20, max 100)")
	flag.StringVar(&detailsID, "details-id", "", "Ticket ID to fetch full details")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	page, perPage, err := normalizePagination(page, perPage)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	filters := searchFilters{
		Queue:     queue,
		Status:    status,
		Owner:     owner,
		Requestor: requestor,
		Subject:   subject,
	}

	query, err := buildTicketSQL(filters)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		flag.Usage()
		os.Exit(1)
	}

	baseURL := os.Getenv("RT_BASE_URL")
	token := os.Getenv("RT_TOKEN")
	if baseURL == "" || token == "" {
		fmt.Fprintln(os.Stderr, "Error: RT_BASE_URL and RT_TOKEN environment variables are required.")
		os.Exit(1)
	}

	client := rtutils_lib.NewClient(baseURL, token)
	ctx := context.Background()

	results, err := client.Tickets.Search(ctx, query, page, perPage)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error searching tickets:", err)
		os.Exit(1)
	}

	if results.Total == 0 {
		fmt.Fprintln(os.Stderr, "No tickets found.")
		os.Exit(2)
	}
	if len(results.Items) == 0 {
		fmt.Fprintf(os.Stderr, "Page %d is out of range.\n", page)
		os.Exit(2)
	}

	printSearchResults(results)
	printPaginationSummary(results, page, perPage)

	if detailsID != "" {
		match, ok := findTicketByID(results, detailsID)
		if !ok {
			fmt.Fprintf(os.Stderr, "Ticket %s not found in current results.\n", detailsID)
			os.Exit(2)
		}

		var full *rtutils_lib.Ticket
		if match.URL != "" {
			full, err = client.Tickets.GetByURL(ctx, match.URL)
		} else {
			full, err = client.Tickets.Get(ctx, match.ID)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching ticket details:", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, "\nTicket Details:")
		printTicketDetails(full)
	}
}
