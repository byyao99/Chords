package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func testRouter(t *testing.T) *gin.Engine {
	t.Helper()
	distFS := fstest.MapFS{
		"index.html":       {Data: []byte("<!DOCTYPE html><html><body>spa</body></html>")},
		"assets/index.js":  {Data: []byte("// js")},
		"assets/index.css": {Data: []byte("body{}")},
	}
	r, err := SetupRouter(distFS)
	if err != nil {
		t.Fatalf("SetupRouter: %v", err)
	}
	return r
}

func TestSPAFallback_Root(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("GET / status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "spa") {
		t.Error("GET / should return SPA index.html")
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html…", ct)
	}
}

func TestSPAFallback_PageRoutes(t *testing.T) {
	r := testRouter(t)
	for _, path := range []string{"/chords", "/diatonic", "/scales", "/metronome"} {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("GET %s status = %d, want 200", path, w.Code)
		}
		if !strings.Contains(w.Body.String(), "spa") {
			t.Errorf("GET %s should fall through to SPA index", path)
		}
	}
}

func TestAssetServing(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/assets/index.js", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET /assets/index.js status = %d, want 200", w.Code)
	}
}

func TestAPIUnknown_Returns404JSON(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/does-not-exist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json…", ct)
	}
}

func TestTranspose_Semitone(t *testing.T) {
	r := testRouter(t)
	body := `{"sheet":"Am  C  G  F\nHello world","mode":"semitone","semitones":2}`
	req := httptest.NewRequest("POST", "/api/transpose", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("POST /api/transpose status = %d, want 200", w.Code)
	}

	var resp TransposeResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if !strings.Contains(resp.Result, "Bm") {
		t.Errorf("result should contain Bm: %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "Hello world") {
		t.Errorf("result should preserve lyrics: %q", resp.Result)
	}
}

func TestTranspose_Key(t *testing.T) {
	r := testRouter(t)
	body := `{"sheet":"C  Am  F  G","mode":"key","from_key":"C","to_key":"G"}`
	req := httptest.NewRequest("POST", "/api/transpose", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp TransposeResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if !strings.Contains(resp.Result, "G") {
		t.Errorf("result should contain G: %q", resp.Result)
	}
	if !strings.Contains(resp.Result, "Em") {
		t.Errorf("result should contain Em: %q", resp.Result)
	}
}

func TestTranspose_InvalidMode(t *testing.T) {
	r := testRouter(t)
	body := `{"sheet":"C Am","mode":"invalid"}`
	req := httptest.NewRequest("POST", "/api/transpose", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestTranspose_EmptySheet(t *testing.T) {
	r := testRouter(t)
	body := `{"sheet":"","mode":"semitone","semitones":1}`
	req := httptest.NewRequest("POST", "/api/transpose", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestTranspose_InvalidKey(t *testing.T) {
	r := testRouter(t)
	body := `{"sheet":"C Am","mode":"key","from_key":"X","to_key":"G"}`
	req := httptest.NewRequest("POST", "/api/transpose", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestScalesAPI_Default(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/scales", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("GET /api/scales status = %d, want 200", w.Code)
	}

	var resp struct {
		Scales []ScaleCard `json:"scales"`
	}
	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp.Scales) != 16 {
		t.Fatalf("expected 16 scales, got %d", len(resp.Scales))
	}
	if resp.Scales[0].Name != "C Major" {
		t.Errorf("first scale name = %q, want %q", resp.Scales[0].Name, "C Major")
	}
}

func TestScalesAPI_WithRoot(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/scales?root=G", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp struct {
		Scales []ScaleCard `json:"scales"`
	}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Scales[0].Name != "G Major" {
		t.Errorf("first scale name = %q, want %q", resp.Scales[0].Name, "G Major")
	}
	if !strings.Contains(resp.Scales[0].Display, "F#") {
		t.Errorf("G major should contain F#: %q", resp.Scales[0].Display)
	}
}

func TestScalesAPI_InvalidRoot(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/scales?root=X", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestScalesAPI_TypeFilter(t *testing.T) {
	r := testRouter(t)

	t.Run("single type", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/scales?root=C&type=Major", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var resp struct {
			Scales []ScaleCard `json:"scales"`
		}
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp.Scales) != 1 {
			t.Fatalf("expected 1 scale, got %d", len(resp.Scales))
		}
		if resp.Scales[0].Name != "C Major" {
			t.Errorf("name = %q, want %q", resp.Scales[0].Name, "C Major")
		}
	})

	t.Run("multiple types", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/scales?root=C&type=Major,Dorian", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var resp struct {
			Scales []ScaleCard `json:"scales"`
		}
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp.Scales) != 2 {
			t.Fatalf("expected 2 scales, got %d", len(resp.Scales))
		}
	})

	t.Run("no type param returns all", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/scales?root=C", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var resp struct {
			Scales []ScaleCard `json:"scales"`
		}
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp.Scales) != 16 {
			t.Fatalf("expected 16 scales, got %d", len(resp.Scales))
		}
	})

	t.Run("nonexistent type returns empty", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/scales?root=C&type=Nonexistent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp struct {
			Scales []ScaleCard `json:"scales"`
		}
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp.Scales) != 0 {
			t.Fatalf("expected 0 scales, got %d", len(resp.Scales))
		}
	})
}

// ── Diatonic Tests ──

func TestDiatonicAPI_CMajor(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/diatonic?root=C&minor=false", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp struct {
		Key    string          `json:"key"`
		Chords json.RawMessage `json:"chords"`
	}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Key != "C Major" {
		t.Errorf("key = %q, want %q", resp.Key, "C Major")
	}
	if !strings.Contains(string(resp.Chords), "Dm") {
		t.Error("C major diatonic should contain Dm")
	}
}

func TestDiatonicAPI_AMinor(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/diatonic?root=A&minor=true", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp struct {
		Key string `json:"key"`
	}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Key != "A Minor" {
		t.Errorf("key = %q, want %q", resp.Key, "A Minor")
	}
}

func TestDiatonicAPI_InvalidRoot(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/diatonic?root=X", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

// ── Chord Diagram Tests ──

func TestChordDiagramAPI_Am(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/chord-diagram?name=Am", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp struct {
		Voicings []struct {
			Name  string `json:"name"`
			Frets [6]int `json:"frets"`
		} `json:"voicings"`
		Notes []string `json:"notes"`
	}
	json.NewDecoder(w.Body).Decode(&resp)

	if len(resp.Voicings) == 0 {
		t.Fatal("expected at least one voicing")
	}
	if resp.Voicings[0].Name != "Am" {
		t.Errorf("name = %q, want Am", resp.Voicings[0].Name)
	}
	if len(resp.Notes) == 0 {
		t.Error("expected notes to be present")
	}
}

func TestChordDiagramAPI_NotFound(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/chord-diagram?name=Xdim13", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestChordDiagramAPI_MissingName(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/chord-diagram", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

// ── Progressions Tests ──

func TestProgressionsAPI_CMajorPop(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/progressions?root=C&style=pop", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("GET /api/progressions status = %d, want 200", w.Code)
	}

	var resp struct {
		Key          string            `json:"key"`
		Progressions []json.RawMessage `json:"progressions"`
	}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Key != "C Major" {
		t.Errorf("key = %q, want C Major", resp.Key)
	}
	if len(resp.Progressions) == 0 {
		t.Error("expected at least one progression")
	}
}

func TestProgressionsAPI_InvalidRoot(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/progressions?root=X", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}
