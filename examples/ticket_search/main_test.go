package ticketsearch
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






































}	}		t.Fatalf("expected %q, got %q", expected, query)	if query != expected {	expected := "Subject LIKE '%printer\\'s%'"	}		t.Fatalf("expected no error, got %v", err)	if err != nil {	query, err := buildTicketSQL(filters)	filters := searchFilters{Subject: "printer's"}func TestBuildTicketSQL_EscapesSingleQuotes(t *testing.T) {}	}		t.Fatalf("expected %q, got %q", expected, query)	if query != expected {	expected := "Queue = 'Helpdesk' AND Status = 'open' AND Subject LIKE '%printer%'"	}		t.Fatalf("expected no error, got %v", err)	if err != nil {	query, err := buildTicketSQL(filters)	}		Subject: "printer",		Status:  "open",		Queue:   "Helpdesk",	filters := searchFilters{func TestBuildTicketSQL_BuildsAndMatches(t *testing.T) {}	}		t.Fatalf("expected error when no filters are provided")	if err == nil {	_, err := buildTicketSQL(searchFilters{})func TestBuildTicketSQL_RequiresFilter(t *testing.T) {