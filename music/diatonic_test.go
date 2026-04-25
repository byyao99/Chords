package music

import (
	"strings"
	"testing"
)

func TestDiatonicChords_CMajor(t *testing.T) {
	chords := DiatonicChords("C", false)
	want := []struct {
		degree, name, quality string
	}{
		{"I", "C", "Major"},
		{"ii", "Dm", "minor"},
		{"iii", "Em", "minor"},
		{"IV", "F", "Major"},
		{"V", "G", "Major"},
		{"vi", "Am", "minor"},
		{"vii°", "Bdim", "dim"},
	}
	if len(chords) != 7 {
		t.Fatalf("expected 7 chords, got %d", len(chords))
	}
	for i, w := range want {
		if chords[i].Degree != w.degree {
			t.Errorf("[%d] degree = %q, want %q", i, chords[i].Degree, w.degree)
		}
		if chords[i].Name != w.name {
			t.Errorf("[%d] name = %q, want %q", i, chords[i].Name, w.name)
		}
		if chords[i].Quality != w.quality {
			t.Errorf("[%d] quality = %q, want %q", i, chords[i].Quality, w.quality)
		}
	}
}

func TestDiatonicChords_GMajor(t *testing.T) {
	chords := DiatonicChords("G", false)
	names := make([]string, 7)
	for i, c := range chords {
		names[i] = c.Name
	}
	got := strings.Join(names, " ")
	want := "G Am Bm C D Em F#dim"
	if got != want {
		t.Errorf("G major diatonic = %q, want %q", got, want)
	}
}

func TestDiatonicChords_FMajor(t *testing.T) {
	chords := DiatonicChords("F", false)
	// F major should use Bb not A#
	if chords[3].Name != "Bb" {
		t.Errorf("IV of F major = %q, want Bb", chords[3].Name)
	}
}

func TestDiatonicChords_AMinor(t *testing.T) {
	chords := DiatonicChords("A", true)
	want := []struct {
		degree, name string
	}{
		{"i", "Am"},
		{"ii°", "Bdim"},
		{"III", "C"},
		{"iv", "Dm"},
		{"v", "Em"},
		{"VI", "F"},
		{"VII", "G"},
	}
	for i, w := range want {
		if chords[i].Degree != w.degree {
			t.Errorf("[%d] degree = %q, want %q", i, chords[i].Degree, w.degree)
		}
		if chords[i].Name != w.name {
			t.Errorf("[%d] name = %q, want %q", i, chords[i].Name, w.name)
		}
	}
}

func TestDiatonicChords_CSharpMajor(t *testing.T) {
	chords := DiatonicChords("C#", false)
	names := make([]string, 7)
	for i, c := range chords {
		names[i] = c.Name
	}
	got := strings.Join(names, " ")
	want := "C# D#m E#m F# G# A#m B#dim"
	if got != want {
		t.Errorf("C# major diatonic = %q, want %q", got, want)
	}
}

func TestDiatonicChords_EbMajor(t *testing.T) {
	chords := DiatonicChords("Eb", false)
	names := make([]string, 7)
	for i, c := range chords {
		names[i] = c.Name
	}
	got := strings.Join(names, " ")
	want := "Eb Fm Gm Ab Bb Cm Ddim"
	if got != want {
		t.Errorf("Eb major diatonic = %q, want %q", got, want)
	}
}
