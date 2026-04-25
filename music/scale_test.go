package music

import (
	"fmt"
	"testing"
)

func TestGenerateScale_CMajor(t *testing.T) {
	result := GenerateScale(C, ScaleTypes[0], KeyFromNote(C))
	want := "C - D - E - F - G - A - B"
	if result.Display != want {
		t.Errorf("C major scale = %q, want %q", result.Display, want)
	}
	if result.Name != "C Major" {
		t.Errorf("name = %q, want %q", result.Name, "C Major")
	}
}

func TestGenerateScale_GMajor(t *testing.T) {
	result := GenerateScale(G, ScaleTypes[0], KeyFromNote(G))
	want := "G - A - B - C - D - E - F#"
	if result.Display != want {
		t.Errorf("G major scale = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_FMajor(t *testing.T) {
	result := GenerateScale(F, ScaleTypes[0], KeyFromNote(F))
	want := "F - G - A - Bb - C - D - E"
	if result.Display != want {
		t.Errorf("F major scale = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_AMinor(t *testing.T) {
	key, _ := KeyFromName("Am")
	result := GenerateScale(A, ScaleTypes[1], key)
	want := "A - B - C - D - E - F - G"
	if result.Display != want {
		t.Errorf("A minor scale = %q, want %q", result.Display, want)
	}
	if result.Name != "A Natural Minor (Aeolian)" {
		t.Errorf("name = %q, want %q", result.Name, "A Natural Minor (Aeolian)")
	}
}

func TestGenerateScale_CMinor(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[1], key)
	want := "C - D - Eb - F - G - Ab - Bb"
	if result.Display != want {
		t.Errorf("C minor scale = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CMajorPentatonic(t *testing.T) {
	result := GenerateScale(C, ScaleTypes[2], KeyFromNote(C))
	want := "C - D - E - G - A"
	if result.Display != want {
		t.Errorf("C major pentatonic = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CMinorPentatonic(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[3], key)
	want := "C - Eb - F - G - Bb"
	if result.Display != want {
		t.Errorf("C minor pentatonic = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CMinorBlues(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[5], key)
	want := "C - Eb - F - Gb - G - Bb"
	if result.Display != want {
		t.Errorf("C minor blues = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_GMajorBlues(t *testing.T) {
	result := GenerateScale(G, ScaleTypes[4], KeyFromNote(G))
	want := "G - A - A# - B - D - E"
	if result.Display != want {
		t.Errorf("G major blues = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CDorian(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[6], key)
	want := "C - D - Eb - F - G - A - Bb"
	if result.Display != want {
		t.Errorf("C Dorian = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CPhrygian(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[7], key)
	want := "C - Db - Eb - F - G - Ab - Bb"
	if result.Display != want {
		t.Errorf("C Phrygian = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CLydian(t *testing.T) {
	result := GenerateScale(C, ScaleTypes[8], KeyFromNote(C))
	want := "C - D - E - F# - G - A - B"
	if result.Display != want {
		t.Errorf("C Lydian = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_GMixolydian(t *testing.T) {
	result := GenerateScale(G, ScaleTypes[9], KeyFromNote(G))
	want := "G - A - B - C - D - E - F"
	if result.Display != want {
		t.Errorf("G Mixolydian = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CLocrian(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[10], key)
	want := "C - Db - Eb - F - Gb - Ab - Bb"
	if result.Display != want {
		t.Errorf("C Locrian = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_AHarmonicMinor(t *testing.T) {
	key, _ := KeyFromName("Am")
	result := GenerateScale(A, ScaleTypes[11], key)
	want := "A - B - C - D - E - F - G#"
	if result.Display != want {
		t.Errorf("A Harmonic Minor = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CMelodicMinor(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[12], key)
	want := "C - D - Eb - F - G - A - B"
	if result.Display != want {
		t.Errorf("C Melodic Minor = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CWholeTone(t *testing.T) {
	result := GenerateScale(C, ScaleTypes[13], KeyFromNote(C))
	want := "C - D - E - F# - G# - A#"
	if result.Display != want {
		t.Errorf("C Whole Tone = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CDiminishedHW(t *testing.T) {
	key, _ := KeyFromName("Cm")
	result := GenerateScale(C, ScaleTypes[14], key)
	want := "C - Db - Eb - E - Gb - G - A - Bb"
	if result.Display != want {
		t.Errorf("C Diminished HW = %q, want %q", result.Display, want)
	}
}

func TestGenerateScale_CDiminishedWH(t *testing.T) {
	result := GenerateScale(C, ScaleTypes[15], KeyFromNote(C))
	want := "C - D - D# - F - F# - G# - A - B"
	if result.Display != want {
		t.Errorf("C Diminished WH = %q, want %q", result.Display, want)
	}
}

func TestScaleDegrees(t *testing.T) {
	tests := []struct {
		name    string
		root    Note
		idx     int
		key     Key
		degrees []string
	}{
		{"C Major", C, 0, KeyFromNote(C), []string{"1", "2", "3", "4", "5", "6", "7"}},
		{"C Natural Minor", C, 1, func() Key { k, _ := KeyFromName("Cm"); return k }(), []string{"1", "2", "b3", "4", "5", "b6", "b7"}},
		{"C Dorian", C, 6, func() Key { k, _ := KeyFromName("Cm"); return k }(), []string{"1", "2", "b3", "4", "5", "6", "b7"}},
		{"C Minor Blues", C, 5, func() Key { k, _ := KeyFromName("Cm"); return k }(), []string{"1", "b3", "4", "b5", "5", "b7"}},
		{"C Whole Tone", C, 13, KeyFromNote(C), []string{"1", "2", "3", "b5", "b6", "b7"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateScale(tt.root, ScaleTypes[tt.idx], tt.key)
			if len(result.Degrees) != len(tt.degrees) {
				t.Fatalf("degrees len = %d, want %d", len(result.Degrees), len(tt.degrees))
			}
			for i, want := range tt.degrees {
				if result.Degrees[i] != want {
					t.Errorf("degree[%d] = %q, want %q (full: %v)", i, result.Degrees[i], want, result.Degrees)
				}
			}
		})
	}
}

func TestDegreesLenMatchesNotes(t *testing.T) {
	results := GenerateAllScales(C, KeyFromNote(C))
	for i, r := range results {
		if len(r.Degrees) != len(r.Notes) {
			t.Errorf("scale[%d] %q: degrees len %d != notes len %d", i, r.Name, len(r.Degrees), len(r.Notes))
		}
	}
}

func TestScaleDegrees_Lydian(t *testing.T) {
	result := GenerateScale(C, ScaleTypes[8], KeyFromNote(C))
	want := []string{"1", "2", "3", "b5", "5", "6", "7"}
	if fmt.Sprint(result.Degrees) != fmt.Sprint(want) {
		t.Errorf("C Lydian degrees = %v, want %v", result.Degrees, want)
	}
}

func TestGenerateScalesByType(t *testing.T) {
	key := KeyFromNote(C)

	t.Run("nil filter returns all", func(t *testing.T) {
		results := GenerateScalesByType(C, key, nil)
		if len(results) != len(ScaleTypes) {
			t.Errorf("got %d, want %d", len(results), len(ScaleTypes))
		}
	})

	t.Run("empty filter returns all", func(t *testing.T) {
		results := GenerateScalesByType(C, key, []string{})
		if len(results) != len(ScaleTypes) {
			t.Errorf("got %d, want %d", len(results), len(ScaleTypes))
		}
	})

	t.Run("single filter", func(t *testing.T) {
		results := GenerateScalesByType(C, key, []string{"Major"})
		if len(results) != 1 {
			t.Fatalf("got %d, want 1", len(results))
		}
		if results[0].Name != "C Major" {
			t.Errorf("name = %q, want %q", results[0].Name, "C Major")
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		results := GenerateScalesByType(C, key, []string{"major"})
		if len(results) != 1 {
			t.Fatalf("got %d, want 1", len(results))
		}
	})

	t.Run("multiple filters", func(t *testing.T) {
		results := GenerateScalesByType(C, key, []string{"Major", "Dorian"})
		if len(results) != 2 {
			t.Fatalf("got %d, want 2", len(results))
		}
	})

	t.Run("nonexistent returns empty", func(t *testing.T) {
		results := GenerateScalesByType(C, key, []string{"Nonexistent"})
		if len(results) != 0 {
			t.Errorf("got %d, want 0", len(results))
		}
	})
}

func TestGenerateAllScales(t *testing.T) {
	results := GenerateAllScales(C, KeyFromNote(C))
	if len(results) != 16 {
		t.Fatalf("expected 16 scales, got %d", len(results))
	}
	names := []string{
		"C Major", "C Natural Minor (Aeolian)",
		"C Major Pentatonic", "C Minor Pentatonic",
		"C Major Blues", "C Minor Blues",
		"C Dorian", "C Phrygian", "C Lydian", "C Mixolydian", "C Locrian",
		"C Harmonic Minor", "C Melodic Minor",
		"C Whole Tone",
		"C Diminished (Half-Whole)", "C Diminished (Whole-Half)",
	}
	for i, want := range names {
		if results[i].Name != want {
			t.Errorf("scale[%d].Name = %q, want %q", i, results[i].Name, want)
		}
	}
}
