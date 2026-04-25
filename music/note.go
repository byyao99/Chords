package music

import (
	"fmt"
	"strings"
)

// Note represents a musical note as a semitone value 0-11 (C=0, C#/Db=1, ... B=11).
type Note int

const (
	C  Note = 0
	Cs Note = 1
	D  Note = 2
	Ds Note = 3
	E  Note = 4
	F  Note = 5
	Fs Note = 6
	G  Note = 7
	Gs Note = 8
	A  Note = 9
	As Note = 10
	B  Note = 11
)

var noteNames = map[string]Note{
	"C":  C,
	"C#": Cs, "Db": Cs,
	"D":  D,
	"D#": Ds, "Eb": Ds,
	"E": E, "Fb": E,
	"E#": F,
	"F":  F,
	"F#": Fs, "Gb": Fs,
	"G":  G,
	"G#": Gs, "Ab": Gs,
	"A":  A,
	"A#": As, "Bb": As,
	"B": B, "Cb": B,
	"B#": C,
}

var sharpNames = [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
var flatNames = [12]string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B"}

// NoteFromName parses a note name from the beginning of s.
// Returns the note, the number of characters consumed, and an error if invalid.
func NoteFromName(s string) (Note, int, error) {
	if len(s) == 0 {
		return 0, 0, fmt.Errorf("empty string")
	}

	first := s[0]
	if first < 'A' || first > 'G' {
		return 0, 0, fmt.Errorf("invalid note: %q", s)
	}

	// Check for accidental (# or b)
	if len(s) >= 2 && (s[1] == '#' || s[1] == 'b') {
		name := s[:2]
		if note, ok := noteNames[name]; ok {
			return note, 2, nil
		}
	}

	name := s[:1]
	if note, ok := noteNames[name]; ok {
		return note, 1, nil
	}

	return 0, 0, fmt.Errorf("invalid note: %q", s)
}

// SharpName returns the note name using sharps (e.g., C#).
func (n Note) SharpName() string {
	return sharpNames[((int(n)%12)+12)%12]
}

// FlatName returns the note name using flats (e.g., Db).
func (n Note) FlatName() string {
	return flatNames[((int(n)%12)+12)%12]
}

// Transpose returns a new note transposed by the given number of semitones.
func (n Note) Transpose(semitones int) Note {
	return Note(((int(n)+semitones)%12 + 12) % 12)
}

// String returns the sharp name by default.
func (n Note) String() string {
	return n.SharpName()
}

// ParseNoteName is a convenience that parses a full string as a note name.
func ParseNoteName(s string) (Note, error) {
	s = strings.TrimSpace(s)
	note, consumed, err := NoteFromName(s)
	if err != nil {
		return 0, err
	}
	if consumed != len(s) {
		return 0, fmt.Errorf("unexpected trailing characters: %q", s[consumed:])
	}
	return note, nil
}
