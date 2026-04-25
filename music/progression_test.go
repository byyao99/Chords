package music

import (
	"testing"
)

func chordNames(chords []DiatonicChord) []string {
	names := make([]string, len(chords))
	for i, c := range chords {
		names[i] = c.Name
	}
	return names
}

func containsName(chords []DiatonicChord, name string) bool {
	for _, c := range chords {
		if c.Name == name {
			return true
		}
	}
	return false
}

func TestGenerateProgressions_CMajorPop(t *testing.T) {
	results, err := GenerateProgressions("C", false, "", "pop", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	var found bool
	for _, r := range results {
		if r.PatternName == "Pop Classic" {
			found = true
			if !containsName(r.Chords, "C") {
				t.Errorf("Pop Classic should contain C: %v", chordNames(r.Chords))
			}
			if !containsName(r.Chords, "G") {
				t.Errorf("Pop Classic should contain G: %v", chordNames(r.Chords))
			}
			if !containsName(r.Chords, "Am") {
				t.Errorf("Pop Classic should contain Am: %v", chordNames(r.Chords))
			}
			if !containsName(r.Chords, "F") {
				t.Errorf("Pop Classic should contain F: %v", chordNames(r.Chords))
			}
		}
	}
	if !found {
		names := make([]string, len(results))
		for i, r := range results {
			names[i] = r.PatternName
		}
		t.Errorf("expected Pop Classic in results, got: %v", names)
	}
}

func TestGenerateProgressions_CSharpMajor(t *testing.T) {
	results, err := GenerateProgressions("C#", false, "", "pop", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	for _, r := range results {
		if r.PatternName == "Pop Classic" {
			for _, ch := range r.Chords {
				if ch.Name == "Fm" {
					t.Errorf("C# major should not contain 'Fm': %v", chordNames(r.Chords))
				}
			}
			if !containsName(r.Chords, "A#m") {
				t.Errorf("C# major pop classic should contain A#m: %v", chordNames(r.Chords))
			}
		}
	}
}

func TestGenerateProgressions_DiatonicChordFields(t *testing.T) {
	results, err := GenerateProgressions("C", false, "", "pop", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.PatternName == "Pop Classic" {
			// Degrees: I V vi IV
			wantDegrees := []string{"I", "V", "vi", "IV"}
			wantQualities := []string{"Major", "Major", "minor", "Major"}
			for i, ch := range r.Chords {
				if ch.Degree != wantDegrees[i] {
					t.Errorf("chord[%d] degree = %q, want %q", i, ch.Degree, wantDegrees[i])
				}
				if ch.Quality != wantQualities[i] {
					t.Errorf("chord[%d] quality = %q, want %q", i, ch.Quality, wantQualities[i])
				}
			}
		}
	}
}

func TestGenerateProgressions_GMinorSad(t *testing.T) {
	results, err := GenerateProgressions("G", true, "sad", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one suggestion for G minor sad")
	}
	for _, r := range results {
		if r.Mood != "sad" && r.Mood != "dark" {
			t.Errorf("expected sad/dark mood, got %q for pattern %q", r.Mood, r.PatternName)
		}
	}
}

func TestGenerateProgressions_NoFilter(t *testing.T) {
	results, err := GenerateProgressions("C", false, "", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) < 4 {
		t.Errorf("expected at least 4 suggestions with no filter, got %d", len(results))
	}
}

func TestGenerateProgressions_Length4(t *testing.T) {
	results, err := GenerateProgressions("C", false, "", "", 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if len(r.Chords) != 4 {
			t.Errorf("pattern %q: expected 4 chords, got %d", r.PatternName, len(r.Chords))
		}
	}
}

func TestGenerateProgressions_InvalidRoot(t *testing.T) {
	_, err := GenerateProgressions("X", false, "", "", 0)
	if err == nil {
		t.Error("expected error for invalid root 'X'")
	}
}

func TestDegreeToChordName_bVII_CMajor(t *testing.T) {
	diatonic := DiatonicChords("C", false)
	degMap := make(map[string]string)
	for _, dc := range diatonic {
		degMap[dc.Degree] = dc.Name
	}
	name := degreeToChordName("bVII", "C", false, degMap)
	if name != "Bb" {
		t.Errorf("bVII in C major = %q, want Bb", name)
	}
}
