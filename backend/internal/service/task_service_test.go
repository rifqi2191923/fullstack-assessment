package service

import "testing"

func TestCacheKeyIncludesQueryParameters(t *testing.T) {
	key := cacheKey("todo", "web", nil, 1, 10, "created_at_desc")
	if key == cacheKey("done", "web", nil, 1, 10, "created_at_desc") {
		t.Fatal("status must affect cache key")
	}
	if key == cacheKey("todo", "mobile", nil, 1, 10, "created_at_desc") {
		t.Fatal("keyword must affect cache key")
	}
	if key == cacheKey("todo", "web", nil, 2, 10, "created_at_desc") {
		t.Fatal("page must affect cache key")
	}
}
