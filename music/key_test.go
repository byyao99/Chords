package music

import "testing"

func TestNoteNameInKey(t *testing.T) {
	tests := []struct {
		key  Key
		note Note
		want string
	}{
		{KeyFromNote(F), As, "Bb"},  // F key uses Bb, not A#
		{KeyFromNote(G), Fs, "F#"},  // G key uses F#, not Gb
		{KeyFromNote(Ds), Gs, "Ab"}, // Eb key uses Ab
		{KeyFromNote(Ds), As, "Bb"}, // Eb key uses Bb
		{KeyFromNote(C), C, "C"},
		{KeyFromNote(A), Cs, "C#"},  // A key uses sharps
		{KeyFromNote(As), Cs, "Db"}, // Bb key uses flats
	}

	for _, tt := range tests {
		got := tt.key.NoteNameInKey(tt.note)
		if got != tt.want {
			t.Errorf("Key(%v, flats=%v).NoteNameInKey(%v) = %q, want %q",
				tt.key.Root, tt.key.UsesFlats, tt.note, got, tt.want)
		}
	}
}

func TestKeyFromName(t *testing.T) {
	tests := []struct {
		input     string
		wantRoot  Note
		wantFlats bool
		wantErr   bool
	}{
		{"C", C, false, false},
		{"G", G, false, false},
		{"D", D, false, false},
		{"A", A, false, false},
		{"E", E, false, false},
		{"B", B, false, false},
		{"F#", Fs, false, false},
		{"F", F, true, false},
		{"Bb", As, true, false},
		{"Eb", Ds, true, false},
		{"Ab", Gs, true, false},
		{"Db", Cs, true, false},
		{"Gb", Fs, true, false},
		// Minor keys
		{"Am", A, false, false},  // relative major = C, no flats
		{"Dm", D, true, false},   // relative major = F, uses flats
		{"Em", E, false, false},  // relative major = G, sharps
		{"Bbm", As, true, false}, // relative major = Db, flats
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, err := KeyFromName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("KeyFromName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if key.Root != tt.wantRoot {
				t.Errorf("KeyFromName(%q).Root = %v, want %v", tt.input, key.Root, tt.wantRoot)
			}
			if key.UsesFlats != tt.wantFlats {
				t.Errorf("KeyFromName(%q).UsesFlats = %v, want %v", tt.input, key.UsesFlats, tt.wantFlats)
			}
		})
	}
}

func TestSemitoneDiff(t *testing.T) {
	tests := []struct {
		from, to Note
		want     int
	}{
		{C, G, 7},
		{C, C, 0},
		{G, C, 5},
		{A, C, 3},
		{E, F, 1},
		{B, C, 1},
	}

	for _, tt := range tests {
		got := SemitoneDiff(tt.from, tt.to)
		if got != tt.want {
			t.Errorf("SemitoneDiff(%v, %v) = %d, want %d", tt.from, tt.to, got, tt.want)
		}
	}
}
