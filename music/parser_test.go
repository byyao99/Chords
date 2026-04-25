package music

import (
	"strings"
	"testing"
)

func TestParseLine_ChordLine(t *testing.T) {
	tests := []struct {
		input      string
		isChord    bool
		chordCount int
	}{
		{"Am C G F", true, 4},
		{"  Am    C      G      F", true, 4},
		{"C", true, 1},
		{"C#m7/G# Bbdim Fmaj7sus4", true, 3},
		{"Hello world", false, 0},
		{"Somebody once told me", false, 0},
		{"", false, 0},
		{"   ", false, 0},
		// Edge: mixed tokens, majority are chords
		{"Am C G the", true, 3}, // 3/4 = 75% > 50%
		// Edge: minority chords
		{"the world Am is", false, 0}, // 1/4 = 25% < 50%
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			line := parseLine(tt.input)
			if line.IsChordLine != tt.isChord {
				t.Errorf("parseLine(%q).IsChordLine = %v, want %v", tt.input, line.IsChordLine, tt.isChord)
			}
			if line.IsChordLine && len(line.Chords) != tt.chordCount {
				t.Errorf("parseLine(%q) chord count = %d, want %d", tt.input, len(line.Chords), tt.chordCount)
			}
		})
	}
}

func TestParseSheet(t *testing.T) {
	input := `  Am        C          G          F
Somebody once told me the world is gonna roll me

  C         G          Am         F
I ain't the sharpest tool in the shed`

	sheet := ParseSheet(input)

	if len(sheet.Lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(sheet.Lines))
	}

	// Line 0: chord line
	if !sheet.Lines[0].IsChordLine {
		t.Error("line 0 should be chord line")
	}
	if len(sheet.Lines[0].Chords) != 4 {
		t.Errorf("line 0 chord count = %d, want 4", len(sheet.Lines[0].Chords))
	}

	// Line 1: lyric line
	if sheet.Lines[1].IsChordLine {
		t.Error("line 1 should be lyric line")
	}

	// Line 2: empty
	if sheet.Lines[2].IsChordLine {
		t.Error("line 2 should not be chord line")
	}

	// Line 3: chord line
	if !sheet.Lines[3].IsChordLine {
		t.Error("line 3 should be chord line")
	}

	// Line 4: lyric line
	if sheet.Lines[4].IsChordLine {
		t.Error("line 4 should be lyric line")
	}
}

func TestRenderSheet_Transpose(t *testing.T) {
	input := `  Am  C  G  F
Somebody once told me`

	sheet := ParseSheet(input)
	// Transpose +2 semitones, target key = B (sharps)
	key := KeyFromNote(B)
	result := RenderSheet(sheet, 2, key)

	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	// Chord line should contain transposed chords
	chordLine := lines[0]
	if !strings.Contains(chordLine, "Bm") {
		t.Errorf("chord line should contain Bm, got: %q", chordLine)
	}
	if !strings.Contains(chordLine, "D") {
		t.Errorf("chord line should contain D, got: %q", chordLine)
	}
	if !strings.Contains(chordLine, "A") {
		t.Errorf("chord line should contain A, got: %q", chordLine)
	}
	if !strings.Contains(chordLine, "G") {
		t.Errorf("chord line should contain G, got: %q", chordLine)
	}

	// Lyric line should be unchanged
	if lines[1] != "Somebody once told me" {
		t.Errorf("lyric line changed: %q", lines[1])
	}
}

func TestRenderSheet_FormattingPreservation(t *testing.T) {
	// Test that chord positions are preserved as much as possible
	input := "C       Am      F       G"
	sheet := ParseSheet(input)

	// Transpose +0 should return similar positions
	result := RenderSheet(sheet, 0, KeyFromNote(C))
	if result != input {
		t.Errorf("transpose +0 should preserve:\n got: %q\nwant: %q", result, input)
	}
}

func TestRenderSheet_SlashChords(t *testing.T) {
	input := "C/E  Am7/G  D/F#"
	sheet := ParseSheet(input)

	// Transpose +1 in Eb key
	key := KeyFromNote(Ds) // Eb
	result := RenderSheet(sheet, 1, key)

	if !strings.Contains(result, "Db/F") {
		t.Errorf("expected Db/F in result, got: %q", result)
	}
	if !strings.Contains(result, "Bbm7/Ab") {
		t.Errorf("expected Bbm7/Ab in result, got: %q", result)
	}
	if !strings.Contains(result, "Eb/G") {
		t.Errorf("expected Eb/G in result, got: %q", result)
	}
}

func TestDetectKey(t *testing.T) {
	input := "  Am  C  G  F\nHello world"
	sheet := ParseSheet(input)
	key := DetectKey(sheet)

	if key.Root != A {
		t.Errorf("DetectKey root = %v, want A", key.Root)
	}
}
