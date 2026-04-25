package music

import "testing"

func TestNoteFromName(t *testing.T) {
	tests := []struct {
		input    string
		wantNote Note
		wantLen  int
		wantErr  bool
	}{
		{"C", C, 1, false},
		{"D", D, 1, false},
		{"E", E, 1, false},
		{"F", F, 1, false},
		{"G", G, 1, false},
		{"A", A, 1, false},
		{"B", B, 1, false},
		{"C#", Cs, 2, false},
		{"Db", Cs, 2, false},
		{"D#", Ds, 2, false},
		{"Eb", Ds, 2, false},
		{"F#", Fs, 2, false},
		{"Gb", Fs, 2, false},
		{"G#", Gs, 2, false},
		{"Ab", Gs, 2, false},
		{"A#", As, 2, false},
		{"Bb", As, 2, false},
		{"Cb", B, 2, false},
		{"Fb", E, 2, false},
		{"E#", F, 2, false},
		{"B#", C, 2, false},
		// With trailing characters
		{"Cm7", C, 1, false},
		{"C#m7", Cs, 2, false},
		{"Bbdim", As, 2, false},
		// Errors
		{"", 0, 0, true},
		{"x", 0, 0, true},
		{"1", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			note, n, err := NoteFromName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NoteFromName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if note != tt.wantNote {
				t.Errorf("NoteFromName(%q) note = %d, want %d", tt.input, note, tt.wantNote)
			}
			if n != tt.wantLen {
				t.Errorf("NoteFromName(%q) len = %d, want %d", tt.input, n, tt.wantLen)
			}
		})
	}
}

func TestNoteTranspose(t *testing.T) {
	tests := []struct {
		note      Note
		semitones int
		want      Note
	}{
		{C, 0, C},
		{C, 2, D},
		{C, 12, C},
		{C, -1, B},
		{A, 3, C},
		{B, 1, C},
		{Fs, -6, C},
	}

	for _, tt := range tests {
		got := tt.note.Transpose(tt.semitones)
		if got != tt.want {
			t.Errorf("%v.Transpose(%d) = %v, want %v", tt.note, tt.semitones, got, tt.want)
		}
	}
}

func TestParseNoteName(t *testing.T) {
	tests := []struct {
		input   string
		want    Note
		wantErr bool
	}{
		{"C", C, false},
		{"Bb", As, false},
		{"F#", Fs, false},
		{" G ", G, false},
		{"Cm", 0, true}, // trailing characters
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseNoteName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseNoteName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("ParseNoteName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
