package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRulesAPI(t *testing.T) {
	// 1. Setup in-memory DB for tests
	initDB("file::memory:?cache=shared")
	router := setupRouter()

	// 2. Test GET /rules
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rules", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %v", w.Code)
	}

	// Because of auto-seed in initDB, we expect 1 rule by default
	var rules []Rule
	json.Unmarshal(w.Body.Bytes(), &rules)
	if len(rules) != 1 {
		t.Fatalf("Expected 1 rule seeded by default, got %v", len(rules))
	}

	// 3. Test POST /rules
	newRule := Rule{
		DeviceID:  "dev-2",
		Condition: "<",
		Threshold: 20.0,
		Message:   "Zbyt niska temperatura",
	}
	body, _ := json.Marshal(newRule)
	wPost := httptest.NewRecorder()
	reqPost, _ := http.NewRequest("POST", "/rules", bytes.NewBuffer(body))
	reqPost.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 created, got %v", wPost.Code)
	}

	// 4. Verify insertion
	wGet2 := httptest.NewRecorder()
	reqGet2, _ := http.NewRequest("GET", "/rules", nil)
	router.ServeHTTP(wGet2, reqGet2)
	
	var rulesAfter []Rule
	json.Unmarshal(wGet2.Body.Bytes(), &rulesAfter)
	if len(rulesAfter) != 2 {
		t.Fatalf("Expected 2 rules after insert, got %v", len(rulesAfter))
	}
}
