package music

import "testing"

func TestTransposeChord(t *testing.T) {
	tests := []struct {
		name      string
		chord     Chord
		semitones int
		key       Key
		want      string
	}{
		{"C+2=D", Chord{C, "", nil}, 2, KeyFromNote(D), "D"},
		{"Am+3=Cm", Chord{A, "m", nil}, 3, KeyFromNote(C), "Cm"},
		{"C#m7/G#+1=Dm7/A", Chord{Cs, "m7", notePtr(Gs)}, 1, KeyFromNote(D), "Dm7/A"},
		{"Bb+2=C", Chord{As, "", nil}, 2, KeyFromNote(C), "C"},
		{"C+12=C (full circle)", Chord{C, "", nil}, 12, KeyFromNote(C), "C"},
		{"C-1=B", Chord{C, "", nil}, -1, KeyFromNote(B), "B"},
		{"D/F#+1=Eb/G in Eb", Chord{D, "", notePtr(Fs)}, 1, KeyFromNote(Ds), "Eb/G"},
		{"Bbdim+2=Cdim", Chord{As, "dim", nil}, 2, KeyFromNote(C), "Cdim"},
		{"Fmaj7+2=Gmaj7", Chord{F, "maj7", nil}, 2, KeyFromNote(G), "Gmaj7"},
		{"Am7b5+5=Dm7b5 in F", Chord{A, "m7b5", nil}, 5, KeyFromNote(F), "Dm7b5"},

		// Flat key rendering
		{"G7+1=Ab7 in Eb", Chord{G, "7", nil}, 1, KeyFromNote(Ds), "Ab7"},
		{"A+3=C in C", Chord{A, "", nil}, 3, KeyFromNote(C), "C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransposeChord(tt.chord, tt.semitones, tt.key)
			got := result.Format(tt.key)
			if got != tt.want {
				t.Errorf("TransposeChord(%v, %d).Format(%v) = %q, want %q",
					tt.chord, tt.semitones, tt.key, got, tt.want)
			}
		})
	}
}

func TestTransposeChordToKey(t *testing.T) {
	tests := []struct {
		name  string
		chord Chord
		from  Key
		to    Key
		want  string
	}{
		{"C->G: C=G", Chord{C, "", nil}, KeyFromNote(C), KeyFromNote(G), "G"},
		{"C->G: Am=Em", Chord{A, "m", nil}, KeyFromNote(C), KeyFromNote(G), "Em"},
		{"C->G: F=C", Chord{F, "", nil}, KeyFromNote(C), KeyFromNote(G), "C"},
		{"C->G: G=D", Chord{G, "", nil}, KeyFromNote(C), KeyFromNote(G), "D"},
		{"G->C: G=C", Chord{G, "", nil}, KeyFromNote(G), KeyFromNote(C), "C"},
		{"Em->Gm: Em=Gm", Chord{E, "m", nil}, KeyFromNote(E), KeyFromNote(G), "Gm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransposeChordToKey(tt.chord, tt.from, tt.to)
			got := result.Format(tt.to)
			if got != tt.want {
				t.Errorf("TransposeChordToKey = %q, want %q", got, tt.want)
			}
		})
	}
}
