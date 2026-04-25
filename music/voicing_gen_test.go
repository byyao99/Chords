package music

import "testing"

// validateVoicing checks that a generated voicing is musically correct.
func validateVoicing(t *testing.T, chord Chord, v ChordVoicing) {
	t.Helper()

	formula, _ := QualityIntervals(chord.Quality)
	required := formula.Required
	if required == nil {
		required = formula.Intervals
	}
	requiredSet := toMod12Set(required)
	rootPC := int(chord.Root) % 12

	// Collect sounding strings.
	var sounding []int // pitch classes
	for s, f := range v.Frets {
		if f == -1 {
			continue
		}
		pc := (int(standardTuning[s]) + f) % 12
		sounding = append(sounding, pc)
	}

	if len(sounding) < 3 {
		t.Errorf("%s: too few sounding strings (%d)", v.Name, len(sounding))
	}

	// All sounding strings must produce chord tones.
	allSet := toMod12Set(formula.Intervals)
	for _, pc := range sounding {
		iv := (pc - rootPC + 12) % 12
		if !allSet[iv] {
			t.Errorf("%s: sounding pitch class %d (interval %d) not in chord", v.Name, pc, iv)
		}
	}

	// All required intervals must be present.
	present := toMod12Set(sounding)
	for iv := range requiredSet {
		needed := (rootPC + iv) % 12
		if !present[needed] {
			t.Errorf("%s: required interval %d (pc %d) not present", v.Name, iv, needed)
		}
	}

	// Fret span must be <= 4.
	minF, maxF := 999, 0
	for _, f := range v.Frets {
		if f > 0 {
			if f < minF {
				minF = f
			}
			if f > maxF {
				maxF = f
			}
		}
	}
	if minF != 999 && maxF-minF > 4 {
		t.Errorf("%s: fret span %d exceeds 4", v.Name, maxF-minF)
	}
}

func TestGenerateVoicing_CommonChords(t *testing.T) {
	cases := []string{
		"C", "Dm", "Em", "F", "G", "Am", "Bm",
		"C7", "G7", "Am7", "Dm7",
		"Cmaj7", "Gmaj7",
	}
	for _, name := range cases {
		chord, ok := ParseChord(name)
		if !ok {
			t.Fatalf("ParseChord(%q) failed", name)
		}
		v, ok := GenerateVoicing(chord)
		if !ok {
			t.Errorf("GenerateVoicing(%q) returned false", name)
			continue
		}
		validateVoicing(t, chord, v)
	}
}

func TestGenerateVoicing_NewChords(t *testing.T) {
	// Chords that are NOT in the hardcoded database.
	cases := []string{
		"Caug", "Faug", "Gaug",
		"Cm7", "Fm7", "Gm7",
		"Dm7b5", "Am7b5",
		"Gadd9", "Cadd9",
		"E9", "A9",
		"Csus2", "Gsus4",
		"Cdim7", "Fdim7",
	}
	for _, name := range cases {
		chord, ok := ParseChord(name)
		if !ok {
			t.Fatalf("ParseChord(%q) failed", name)
		}
		v, ok := GenerateVoicing(chord)
		if !ok {
			t.Errorf("GenerateVoicing(%q) returned false", name)
			continue
		}
		validateVoicing(t, chord, v)
	}
}

func TestGenerateVoicing_AllRootsAllQualities(t *testing.T) {
	roots := []Note{C, Cs, D, Ds, E, F, Fs, G, Gs, A, As, B}
	for _, root := range roots {
		for quality := range qualityToFormula {
			chord := Chord{Root: root, Quality: quality}
			v, ok := GenerateVoicing(chord)
			if !ok {
				t.Errorf("GenerateVoicing(%s%s) returned false", root.SharpName(), quality)
				continue
			}
			validateVoicing(t, chord, v)
		}
	}
}

func TestGenerateVoicing_UnknownQuality(t *testing.T) {
	chord := Chord{Root: C, Quality: "xyz"}
	_, ok := GenerateVoicing(chord)
	if ok {
		t.Error("expected unknown quality to return false")
	}
}

func TestFindChordVoicing_Fallback(t *testing.T) {
	// These chords are not hardcoded; the fallback generator must handle them.
	cases := []string{"Caug", "Fm7", "Dm7b5", "Gadd9"}
	for _, name := range cases {
		v, ok := FindChordVoicing(name)
		if !ok {
			t.Errorf("FindChordVoicing(%q) returned false", name)
			continue
		}
		if v.Name != name {
			t.Errorf("FindChordVoicing(%q): Name = %q, want %q", name, v.Name, name)
		}
	}
}

func TestFindChordVoicing_HardcodedPriority(t *testing.T) {
	// Hardcoded voicings must still be returned for known chords.
	hardcoded := []string{"Am", "G7", "Cmaj7", "Dm"}
	for _, name := range hardcoded {
		v, ok := FindChordVoicing(name)
		if !ok {
			t.Errorf("FindChordVoicing(%q) returned false", name)
			continue
		}
		expected := chordVoicings[name]
		if v.Frets != expected.Frets {
			t.Errorf("FindChordVoicing(%q): expected hardcoded frets %v, got %v", name, expected.Frets, v.Frets)
		}
	}
}
