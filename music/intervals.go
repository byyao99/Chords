package music

// ChordFormula defines the intervals and required tones for a chord quality.
type ChordFormula struct {
	Intervals []int // all chord tones as semitones from root (may exceed 11 for extensions)
	Required  []int // subset that must appear; nil means all are required
}

// qualityToFormula maps chord quality strings to their interval formulas.
var qualityToFormula = map[string]ChordFormula{
	// Triads
	"":    {Intervals: []int{0, 4, 7}},
	"m":   {Intervals: []int{0, 3, 7}},
	"aug": {Intervals: []int{0, 4, 8}},
	"dim": {Intervals: []int{0, 3, 6}},
	"5":   {Intervals: []int{0, 7}},

	// 7th chords
	"7":     {Intervals: []int{0, 4, 7, 10}},
	"m7":    {Intervals: []int{0, 3, 7, 10}},
	"maj7":  {Intervals: []int{0, 4, 7, 11}},
	"dim7":  {Intervals: []int{0, 3, 6, 9}},
	"m7b5":  {Intervals: []int{0, 3, 6, 10}}, // half-diminished
	"aug7":  {Intervals: []int{0, 4, 8, 10}},
	"mmaj7": {Intervals: []int{0, 3, 7, 11}}, // minor-major 7

	// 6th chords
	"6":  {Intervals: []int{0, 4, 7, 9}},
	"m6": {Intervals: []int{0, 3, 7, 9}},

	// Sus chords
	"sus2":  {Intervals: []int{0, 2, 7}},
	"sus4":  {Intervals: []int{0, 5, 7}},
	"7sus4": {Intervals: []int{0, 5, 7, 10}},

	// Add chords
	"add9": {Intervals: []int{0, 4, 7, 14}},

	// Extended chords (5th is optional)
	"9":    {Intervals: []int{0, 4, 7, 10, 14}, Required: []int{0, 4, 10, 14}},
	"m9":   {Intervals: []int{0, 3, 7, 10, 14}, Required: []int{0, 3, 10, 14}},
	"maj9": {Intervals: []int{0, 4, 7, 11, 14}, Required: []int{0, 4, 11, 14}},

	// Altered dominants
	"7b9": {Intervals: []int{0, 4, 7, 10, 13}, Required: []int{0, 4, 10, 13}},
	"7#9": {Intervals: []int{0, 4, 7, 10, 15}, Required: []int{0, 4, 10, 15}},
}

// QualityIntervals returns the ChordFormula for the given quality string.
func QualityIntervals(quality string) (ChordFormula, bool) {
	f, ok := qualityToFormula[quality]
	return f, ok
}

// ChordTones returns the note names that make up the given chord (e.g. "Am7" → ["A","C","E","G"]).
// The accidental spelling (sharp vs flat) follows the root note's conventional key.
func ChordTones(name string) ([]string, bool) {
	chord, ok := ParseChord(name)
	if !ok {
		return nil, false
	}
	formula, ok := QualityIntervals(chord.Quality)
	if !ok {
		return nil, false
	}

	// Determine sharp vs flat based on how the root was written.
	// If the original name contains a flat accidental on the root, use flats.
	usesFlats := keyUsesFlats[chord.Root]
	if len(name) >= 2 && name[1] == 'b' {
		usesFlats = true
	} else if len(name) >= 2 && name[1] == '#' {
		usesFlats = false
	}

	seen := make(map[Note]bool)
	var notes []string
	for _, iv := range formula.Intervals {
		n := chord.Root.Transpose(iv % 12)
		if seen[n] {
			continue
		}
		seen[n] = true
		if usesFlats {
			notes = append(notes, n.FlatName())
		} else {
			notes = append(notes, n.SharpName())
		}
	}
	return notes, true
}
