package music

import "testing"

func notePtr(n Note) *Note { return &n }

func TestParseChord(t *testing.T) {
	tests := []struct {
		input    string
		wantOK   bool
		wantRoot Note
		wantQual string
		wantBass *Note
	}{
		// Basic chords
		{"C", true, C, "", nil},
		{"Am", true, A, "m", nil},
		{"G7", true, G, "7", nil},
		{"Fmaj7", true, F, "maj7", nil},
		{"Dm7", true, D, "m7", nil},

		// Advanced chords
		{"Gsus4", true, G, "sus4", nil},
		{"Gsus2", true, G, "sus2", nil},
		{"Bdim", true, B, "dim", nil},
		{"Caug", true, C, "aug", nil},
		{"Eadd9", true, E, "add9", nil},
		{"A6", true, A, "6", nil},
		{"D9", true, D, "9", nil},

		// Complex chords
		{"Fmaj7sus4", true, F, "maj7sus4", nil},
		{"Am7b5", true, A, "m7b5", nil},
		{"C7#9", true, C, "7#9", nil},

		// With accidentals
		{"C#m7", true, Cs, "m7", nil},
		{"Bbdim", true, As, "dim", nil},
		{"F#11", true, Fs, "11", nil},
		{"Ebmaj7", true, Ds, "maj7", nil},

		// Slash chords
		{"C/E", true, C, "", notePtr(E)},
		{"D/F#", true, D, "", notePtr(Fs)},
		{"Am7/G", true, A, "m7", notePtr(G)},
		{"C#m7/G#", true, Cs, "m7", notePtr(Gs)},

		// Not chords
		{"the", false, 0, "", nil},
		{"Hello", false, 0, "", nil},
		{"123", false, 0, "", nil},
		{"", false, 0, "", nil},
		{"/G", false, 0, "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			chord, ok := ParseChord(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ParseChord(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if chord.Root != tt.wantRoot {
				t.Errorf("ParseChord(%q).Root = %v, want %v", tt.input, chord.Root, tt.wantRoot)
			}
			if chord.Quality != tt.wantQual {
				t.Errorf("ParseChord(%q).Quality = %q, want %q", tt.input, chord.Quality, tt.wantQual)
			}
			if tt.wantBass == nil && chord.BassNote != nil {
				t.Errorf("ParseChord(%q).BassNote = %v, want nil", tt.input, *chord.BassNote)
			}
			if tt.wantBass != nil {
				if chord.BassNote == nil {
					t.Errorf("ParseChord(%q).BassNote = nil, want %v", tt.input, *tt.wantBass)
				} else if *chord.BassNote != *tt.wantBass {
					t.Errorf("ParseChord(%q).BassNote = %v, want %v", tt.input, *chord.BassNote, *tt.wantBass)
				}
			}
		})
	}
}

func TestChordFormat(t *testing.T) {
	tests := []struct {
		chord Chord
		key   Key
		want  string
	}{
		{Chord{C, "", nil}, KeyFromNote(C), "C"},
		{Chord{As, "m7", nil}, KeyFromNote(F), "Bbm7"},
		{Chord{Cs, "m7", notePtr(Gs)}, KeyFromNote(A), "C#m7/G#"},
		{Chord{Cs, "m7", notePtr(Gs)}, Key{Cs, true}, "Dbm7/Ab"},
	}

	for _, tt := range tests {
		got := tt.chord.Format(tt.key)
		if got != tt.want {
			t.Errorf("Chord(%v).Format(%v) = %q, want %q", tt.chord, tt.key, got, tt.want)
		}
	}
}
