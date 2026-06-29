package problem

import (
	"net/http"
	"testing"
)

func TestPageNotFoundProblem(t *testing.T) {
	p := NewProblems()

	got := p.PageNotFoundProblem("/missing")

	if got.Status != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", got.Status, http.StatusNotFound)
	}
	if got.Instance != "/missing" {
		t.Errorf("Instance = %q, want %q", got.Instance, "/missing")
	}
	if got.Type != typeBlank {
		t.Errorf("Type = %q, want %q", got.Type, typeBlank)
	}
}

func TestProblemDefaultsTypeToAboutBlank(t *testing.T) {
	p := NewProblems()

	got := p.Problem("boom", "/x", http.StatusBadRequest, "Bad request", "")

	if got.Type != typeBlank {
		t.Errorf("Type = %q, want %q", got.Type, typeBlank)
	}
	if got.Status != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", got.Status, http.StatusBadRequest)
	}
}
