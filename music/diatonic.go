package music

// DiatonicChord represents one chord in a diatonic chord set.
type DiatonicChord struct {
	Degree  string `json:"degree"`  // Roman numeral: "I", "ii", "iii", etc.
	Name    string `json:"name"`    // Chord name: "C", "Dm", "Em", etc.
	Quality string `json:"quality"` // "Major", "minor", "dim"
}

// Major key: scale degree semitone offsets from root
var majorIntervals = [7]int{0, 2, 4, 5, 7, 9, 11}

// Major key chord qualities and Roman numerals
var majorDegrees = [7]struct {
	degree  string
	quality string
	suffix  string
}{
	{"I", "Major", ""},
	{"ii", "minor", "m"},
	{"iii", "minor", "m"},
	{"IV", "Major", ""},
	{"V", "Major", ""},
	{"vi", "minor", "m"},
	{"vii°", "dim", "dim"},
}

// Natural minor: scale degree semitone offsets from root
var minorIntervals = [7]int{0, 2, 3, 5, 7, 8, 10}

// Minor key chord qualities and Roman numerals
var minorDegrees = [7]struct {
	degree  string
	quality string
	suffix  string
}{
	{"i", "minor", "m"},
	{"ii°", "dim", "dim"},
	{"III", "Major", ""},
	{"iv", "minor", "m"},
	{"v", "minor", "m"},
	{"VI", "Major", ""},
	{"VII", "Major", ""},
}

// naturalSemitones maps letter index (0=C,1=D,2=E,3=F,4=G,5=A,6=B) to semitone value.
var naturalSemitones = [7]int{0, 2, 4, 5, 7, 9, 11}

// diatonicLetters lists the note letters in ascending order starting from C.
var diatonicLetters = [7]string{"C", "D", "E", "F", "G", "A", "B"}

// letterIdx maps a letter byte to its index in diatonicLetters.
var letterIdx = map[byte]int{'C': 0, 'D': 1, 'E': 2, 'F': 3, 'G': 4, 'A': 5, 'B': 6}

// DiatonicNoteName returns the properly spelled note name for a given scale degree.
// rootLetterIdx is the letter index of the scale root (e.g. 0 for C, 2 for E).
// degree is the scale step (0–6). note is the actual semitone value.
func DiatonicNoteName(rootLetterIdx, degree int, note Note) string {
	idx := (rootLetterIdx + degree) % 7
	letter := diatonicLetters[idx]
	natural := naturalSemitones[idx]
	diff := (int(note) - natural + 12) % 12
	// Normalise to range -2..2 (prefer nearby accidentals over distant ones).
	if diff > 6 {
		diff -= 12
	}
	switch diff {
	case 0:
		return letter
	case 1:
		return letter + "#"
	case -1:
		return letter + "b"
	case 2:
		return letter + "##"
	case -2:
		return letter + "bb"
	default:
		return letter // fallback; shouldn't occur in standard keys
	}
}

// DiatonicChords returns the 7 diatonic chords for a given root name and mode.
// rootName must be the spelled root (e.g. "C#", "Eb") so that diatonic letter
// names are derived correctly (e.g. E# instead of F in C# major).
func DiatonicChords(rootName string, minor bool) []DiatonicChord {
	root, _, _ := NoteFromName(rootName)
	rootLetterIdx := letterIdx[rootName[0]]

	intervals := majorIntervals
	degrees := majorDegrees
	if minor {
		intervals = minorIntervals
		degrees = minorDegrees
	}

	chords := make([]DiatonicChord, 7)
	for i := 0; i < 7; i++ {
		note := root.Transpose(intervals[i])
		noteName := DiatonicNoteName(rootLetterIdx, i, note)
		chords[i] = DiatonicChord{
			Degree:  degrees[i].degree,
			Name:    noteName + degrees[i].suffix,
			Quality: degrees[i].quality,
		}
	}
	return chords
}
