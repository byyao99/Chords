package music

import "regexp"

// Chord represents a parsed chord symbol.
type Chord struct {
	Root     Note
	Quality  string // e.g., "m7", "maj7", "sus4", "dim", "" (major)
	BassNote *Note  // non-nil for slash chords like C/E
}

// chordRegex matches chord symbols: root note + quality + optional slash bass.
// Examples: C, Am7, F#dim, Bbmaj7, C/E, Am7/G, C#m7/G#
var chordRegex = regexp.MustCompile(`^([A-G][#b]?)(.*?)(?:/([A-G][#b]?))?\s*$`)

// ParseChord parses a chord token like "Am7", "C#m7/G#", "Bbdim".
// Returns the parsed chord and true if successful, or zero value and false if not a valid chord.
func ParseChord(token string) (Chord, bool) {
	if len(token) == 0 {
		return Chord{}, false
	}

	m := chordRegex.FindStringSubmatch(token)
	if m == nil {
		return Chord{}, false
	}

	rootName := m[1]
	quality := m[2]
	bassName := m[3]

	root, _, err := NoteFromName(rootName)
	if err != nil {
		return Chord{}, false
	}

	chord := Chord{
		Root:    root,
		Quality: quality,
	}

	if bassName != "" {
		bass, _, err := NoteFromName(bassName)
		if err != nil {
			return Chord{}, false
		}
		chord.BassNote = &bass
	}

	return chord, true
}

// String returns the chord symbol using sharp names.
func (c Chord) String() string {
	return c.Format(Key{Root: C, UsesFlats: false})
}

// Format returns the chord symbol with note names determined by the given key.
func (c Chord) Format(key Key) string {
	s := key.NoteNameInKey(c.Root) + c.Quality
	if c.BassNote != nil {
		s += "/" + key.NoteNameInKey(*c.BassNote)
	}
	return s
}
