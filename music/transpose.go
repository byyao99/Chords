package music

// TransposeChord transposes a chord by the given number of semitones.
// The key determines the enharmonic spelling of the result.
func TransposeChord(c Chord, semitones int, key Key) Chord {
	newRoot := c.Root.Transpose(semitones)
	result := Chord{
		Root:    newRoot,
		Quality: c.Quality,
	}
	if c.BassNote != nil {
		newBass := c.BassNote.Transpose(semitones)
		result.BassNote = &newBass
	}
	return result
}

// TransposeChordToKey transposes a chord from one key to another.
func TransposeChordToKey(c Chord, from, to Key) Chord {
	semitones := SemitoneDiff(from.Root, to.Root)
	return TransposeChord(c, semitones, to)
}
