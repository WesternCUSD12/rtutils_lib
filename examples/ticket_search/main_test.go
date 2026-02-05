package main

import "testing"

func TestNormalizePagination_Defaults(t *testing.T) {
	page, perPage, err := normalizePagination(0, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if page != 1 {
		t.Fatalf("expected page=1, got %d", page)
	}
	if perPage != 20 {
		t.Fatalf("expected perPage=20, got %d", perPage)
	}
}

func TestNormalizePagination_MaxPerPage(t *testing.T) {
	_, _, err := normalizePagination(1, 101)
	if err == nil {
		t.Fatalf("expected error for perPage>100")
	}
}

func TestBuildTicketSQL_RequiresFilter(t *testing.T) {
	_, err := buildTicketSQL(searchFilters{})
	if err == nil {
		t.Fatalf("expected error when no filters are provided")
	}
}

func TestBuildTicketSQL_BuildsAndMatches(t *testing.T) {
	filters := searchFilters{
		Queue:   "Helpdesk",
		Status:  "open",
		Subject: "printer",
	}
	query, err := buildTicketSQL(filters)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "Queue = 'Helpdesk' AND Status = 'open' AND Subject LIKE '%printer%'"
	if query != expected {
		t.Fatalf("expected %q, got %q", expected, query)
	}
}

func TestBuildTicketSQL_EscapesSingleQuotes(t *testing.T) {
	filters := searchFilters{Subject: "printer's"}
	query, err := buildTicketSQL(filters)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "Subject LIKE '%printer\\'s%'"
	if query != expected {
		t.Fatalf("expected %q, got %q", expected, query)
	}
}
