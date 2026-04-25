package music

// Key represents a musical key with its root note and accidental preference.
type Key struct {
	Root      Note
	UsesFlats bool
}

// usesFlats maps root notes to whether the key conventionally uses flats.
// Sharp keys: G, D, A, E, B, F#
// Flat keys: F, Bb, Eb, Ab, Db, Gb
// C: no accidentals, defaults to sharps
var keyUsesFlats = map[Note]bool{
	C:  false,
	Cs: false, // C# / Db — treat as Db (flats) only when user writes Db
	D:  false,
	Ds: true, // Eb
	E:  false,
	F:  true,
	Fs: false, // F#
	G:  false,
	Gs: true, // Ab
	A:  false,
	As: true, // Bb
	B:  false,
}

// KeyFromNote creates a Key from a note, using the conventional accidental preference.
func KeyFromNote(root Note) Key {
	return Key{Root: root, UsesFlats: keyUsesFlats[root]}
}

// KeyFromName parses a key name like "C", "Bb", "F#", "Am", "Ebm".
// For minor keys, the relative major is used for accidental preference.
func KeyFromName(s string) (Key, error) {
	note, consumed, err := NoteFromName(s)
	if err != nil {
		return Key{}, err
	}

	rest := s[consumed:]

	// Determine if the user wrote it as a flat (e.g., "Db" vs "C#")
	usesFlats := keyUsesFlats[note]
	if consumed == 2 && s[1] == 'b' {
		usesFlats = true
	} else if consumed == 2 && s[1] == '#' {
		usesFlats = false
	}

	// Handle minor key suffix — use relative major for accidental preference
	if rest == "m" || rest == "min" {
		relativeMajor := note.Transpose(3)
		usesFlats = keyUsesFlats[relativeMajor]
		// Override based on the accidental written
		if consumed == 2 && s[1] == 'b' {
			usesFlats = true
		} else if consumed == 2 && s[1] == '#' {
			usesFlats = false
		}
	}

	return Key{Root: note, UsesFlats: usesFlats}, nil
}

// NoteNameInKey returns the display name for a note in this key's context.
func (k Key) NoteNameInKey(n Note) string {
	normalized := Note(((int(n) % 12) + 12) % 12)
	if k.UsesFlats {
		return flatNames[normalized]
	}
	return sharpNames[normalized]
}

// SemitoneDiff returns the number of semitones from one note to another (ascending).
func SemitoneDiff(from, to Note) int {
	diff := int(to) - int(from)
	return ((diff % 12) + 12) % 12
}
