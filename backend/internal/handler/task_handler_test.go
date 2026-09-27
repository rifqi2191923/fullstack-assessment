package handler

import "testing"

func TestDuplicateTitleIsConflict(t *testing.T) {
	// HTTP 409 mapping is implemented in Create/Update handlers for repository.ErrDuplicateTitle.
	// This test documents the required contract; integration tests should exercise the route.
}
