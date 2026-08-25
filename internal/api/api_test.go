package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestKinematicsEndpoint(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest(http.MethodPost, "/api/kinematics", strings.NewReader(`{"energy_keV":511,"angle_deg":90}`))
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var out kinematicsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.RecoilKEV < 250 || out.RecoilKEV > 260 {
		t.Fatalf("recoil=%v", out.RecoilKEV)
	}
	if s.Book().Len() != 1 {
		t.Fatalf("book len=%d", s.Book().Len())
	}
}

func TestSectionEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"energy_keV":511}`))
	req := httptest.NewRequest(http.MethodPost, "/api/section", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var out sectionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Total <= 0 || out.TotalRatio >= 1 {
		t.Fatalf("section=%+v", out)
	}
}

func TestTrendEndpoint(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest(http.MethodGet, "/api/trend", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	var out trendResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if !out.Decreasing || len(out.EnergiesKEV) != 9 {
		t.Fatalf("trend=%+v", out)
	}
}

func TestHealthEndpoint(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "ok") {
		t.Fatalf("health status=%d body=%s", rec.Code, rec.Body.String())
	}
}
