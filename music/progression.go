package music

import "strings"

// ProgressionPattern is a named chord-progression template.
type ProgressionPattern struct {
	Name    string
	Degrees []string // e.g. ["I", "V", "vi", "IV"]
	Mood    string   // happy / sad / energetic / dark / chill
	Style   string   // pop / rock / jazz / blues / folk
}

// ProgressionSuggestion is a generated progression with concrete chord cards.
type ProgressionSuggestion struct {
	PatternName string          `json:"pattern_name"`
	Mood        string          `json:"mood"`
	Style       string          `json:"style"`
	Chords      []DiatonicChord `json:"chords"`
}

// patternDatabase is the built-in library of progression patterns.
var patternDatabase = []ProgressionPattern{
	// Pop
	{"Pop Classic", []string{"I", "V", "vi", "IV"}, "happy", "pop"},
	{"Pop Ballad", []string{"vi", "IV", "I", "V"}, "sad", "pop"},
	{"Pop Alt", []string{"I", "IV", "vi", "V"}, "happy", "pop"},
	{"Pop Descending", []string{"I", "V", "IV", "I"}, "happy", "pop"},
	// Rock
	{"Rock Standard", []string{"I", "IV", "V"}, "energetic", "rock"},
	{"Rock Power", []string{"I", "bVII", "IV"}, "energetic", "rock"},
	{"Rock Alt", []string{"I", "V", "bVII", "IV"}, "energetic", "rock"},
	{"Rock Anthem", []string{"I", "IV", "I", "V"}, "energetic", "rock"},
	// Jazz
	{"ii-V-I", []string{"ii", "V", "I"}, "chill", "jazz"},
	{"Rhythm Changes", []string{"I", "vi", "ii", "V"}, "chill", "jazz"},
	{"Jazz Turnaround", []string{"iii", "vi", "ii", "V"}, "chill", "jazz"},
	{"Minor ii-V-i", []string{"ii°", "V", "i"}, "dark", "jazz"},
	// Blues
	{"12-Bar Blues", []string{"I", "I", "I", "I", "IV", "IV", "I", "I", "V", "IV", "I", "V"}, "chill", "blues"},
	{"8-Bar Blues", []string{"I", "IV", "I", "V", "IV", "I", "V", "I"}, "chill", "blues"},
	// Folk
	{"Folk Classic", []string{"I", "IV", "I", "V"}, "happy", "folk"},
	{"Folk Ballad", []string{"I", "V", "IV", "I"}, "chill", "folk"},
	// Sad / Minor
	{"Minor Sad", []string{"i", "VI", "III", "VII"}, "sad", "pop"},
	{"Minor Dark", []string{"i", "iv", "v", "i"}, "dark", "rock"},
	{"Andalusian", []string{"i", "VII", "VI", "V"}, "dark", "folk"},
	// Neo-soul / R&B
	{"Neo-Soul", []string{"I", "iii", "IV", "V"}, "chill", "jazz"},
	{"R&B Smooth", []string{"I", "vi", "IV", "V"}, "chill", "jazz"},
}

// nonDiatonicSemitones maps special degree tokens to their semitone offset from root.
var nonDiatonicSemitones = map[string]int{
	"bII": 1, "bIII": 3, "bV": 6, "bVI": 8, "bVII": 10,
	"#I": 1, "#II": 3, "#IV": 6, "#V": 8,
}

// GenerateProgressions returns progression suggestions for the given key and filters.
// Pass empty strings for mood/style to skip those filters. Pass 0 for length to skip.
func GenerateProgressions(rootName string, minor bool, mood, style string, length int) ([]ProgressionSuggestion, error) {
	// Validate root.
	if _, err := ParseNoteName(rootName); err != nil {
		return nil, err
	}

	// Build degree→chord map from diatonic chords.
	diatonic := DiatonicChords(rootName, minor)
	degMap := make(map[string]string, 7)
	for _, dc := range diatonic {
		degMap[dc.Degree] = dc.Name
	}

	var results []ProgressionSuggestion
	for _, pat := range patternDatabase {
		if mood != "" && !strings.EqualFold(pat.Mood, mood) {
			continue
		}
		if style != "" && !strings.EqualFold(pat.Style, style) {
			continue
		}

		degrees := pat.Degrees
		if length > 0 && length != len(degrees) {
			// Repeat or trim to match requested length.
			degrees = adjustLength(degrees, length)
		}

		chords := degreesToChords(degrees, rootName, minor, degMap)
		results = append(results, ProgressionSuggestion{
			PatternName: pat.Name,
			Mood:        pat.Mood,
			Style:       pat.Style,
			Chords:      chords,
		})
		if len(results) >= 6 {
			break
		}
	}
	return results, nil
}

// degreesToChords converts a degree slice to DiatonicChord cards in the given key.
func degreesToChords(degrees []string, rootName string, minor bool, degMap map[string]string) []DiatonicChord {
	chords := make([]DiatonicChord, len(degrees))
	for i, deg := range degrees {
		name := degreeToChordName(deg, rootName, minor, degMap)
		chords[i] = DiatonicChord{
			Degree:  deg,
			Name:    name,
			Quality: degreeQuality(deg),
		}
	}
	return chords
}

// degreeQuality returns the chord quality for a roman-numeral degree.
func degreeQuality(degree string) string {
	switch degree {
	case "ii", "iii", "vi", "i", "iv", "v":
		return "minor"
	case "vii°", "ii°":
		return "dim"
	default:
		return "Major"
	}
}

// degreeToChordName converts a single roman-numeral degree to a chord name.
func degreeToChordName(degree, rootName string, minor bool, degMap map[string]string) string {
	// Direct lookup in diatonic map.
	if name, ok := degMap[degree]; ok {
		return name
	}

	// Handle non-diatonic degrees like bVII, bIII.
	offset, ok := nonDiatonicSemitones[degree]
	if !ok {
		// Strip lower/upper to try a case-insensitive lookup or return as-is.
		return degree
	}

	root, _, _ := NoteFromName(rootName)
	targetNote := root.Transpose(offset)
	rootLetterIdx := letterIdx[rootName[0]]

	// Determine which diatonic degree position this corresponds to.
	// bVII, bIII etc. are always major chords.
	// We find the letter by the degree name's roman numeral ordinal.
	degreeOrdinal := parseDegreeOrdinal(degree)
	if degreeOrdinal < 0 {
		return targetNote.SharpName()
	}

	letterIdx2 := (rootLetterIdx + degreeOrdinal) % 7
	noteName := DiatonicNoteName(rootLetterIdx, degreeOrdinal, targetNote)
	_ = letterIdx2

	// Non-diatonic bII, bVII etc. are major chords.
	return noteName
}

// parseDegreeOrdinal returns the 0-based ordinal of a roman numeral (I=0, II=1, ..., VII=6).
// Returns -1 if unparseable.
func parseDegreeOrdinal(degree string) int {
	// Strip accidentals and case.
	d := strings.TrimLeft(degree, "b#")
	d = strings.ToUpper(d)
	switch d {
	case "I":
		return 0
	case "II":
		return 1
	case "III":
		return 2
	case "IV":
		return 3
	case "V":
		return 4
	case "VI":
		return 5
	case "VII":
		return 6
	}
	return -1
}

// adjustLength repeats or trims a degree slice to the target length.
func adjustLength(degrees []string, length int) []string {
	if len(degrees) == 0 {
		return degrees
	}
	result := make([]string, length)
	for i := range result {
		result[i] = degrees[i%len(degrees)]
	}
	return result
}
