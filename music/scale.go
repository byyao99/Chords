package music

import "strings"

// ScaleType defines a named scale with its interval pattern.
type ScaleType struct {
	Name      string // e.g., "大調音階"
	Intervals []int  // semitone steps between consecutive notes
}

var ScaleTypes = []ScaleType{
	{"Major", []int{2, 2, 1, 2, 2, 2, 1}},
	{"Natural Minor (Aeolian)", []int{2, 1, 2, 2, 1, 2, 2}},
	{"Major Pentatonic", []int{2, 2, 3, 2, 3}},
	{"Minor Pentatonic", []int{3, 2, 2, 3, 2}},
	{"Major Blues", []int{2, 1, 1, 3, 2, 3}},
	{"Minor Blues", []int{3, 2, 1, 1, 3, 2}},
	{"Dorian", []int{2, 1, 2, 2, 2, 1, 2}},
	{"Phrygian", []int{1, 2, 2, 2, 1, 2, 2}},
	{"Lydian", []int{2, 2, 2, 1, 2, 2, 1}},
	{"Mixolydian", []int{2, 2, 1, 2, 2, 1, 2}},
	{"Locrian", []int{1, 2, 2, 1, 2, 2, 2}},
	{"Harmonic Minor", []int{2, 1, 2, 2, 1, 3, 1}},
	{"Melodic Minor", []int{2, 1, 2, 2, 2, 2, 1}},
	{"Whole Tone", []int{2, 2, 2, 2, 2, 2}},
	{"Diminished (Half-Whole)", []int{1, 2, 1, 2, 1, 2, 1, 2}},
	{"Diminished (Whole-Half)", []int{2, 1, 2, 1, 2, 1, 2, 1}},
}

// ScaleResult holds the generated scale info for display.
type ScaleResult struct {
	Name    string   // e.g., "G Major"
	Notes   []string // e.g., ["G", "A", "B", "C", "D", "E", "F#"]
	Degrees []string // e.g., ["1", "2", "3", "4", "5", "6", "7"]
	Display string   // e.g., "G - A - B - C - D - E - F#"
}

var semitoneToDegree = map[int]string{
	0: "1", 1: "b2", 2: "2", 3: "b3", 4: "3", 5: "4",
	6: "b5", 7: "5", 8: "b6", 9: "6", 10: "b7", 11: "7",
}

// GenerateScale produces the notes of a scale given a root and scale type.
func GenerateScale(root Note, st ScaleType, key Key) ScaleResult {
	notes := make([]string, 0, len(st.Intervals))
	degrees := make([]string, 0, len(st.Intervals))
	current := root
	notes = append(notes, key.NoteNameInKey(current))
	degrees = append(degrees, "1")

	for i := 0; i < len(st.Intervals)-1; i++ {
		current = current.Transpose(st.Intervals[i])
		notes = append(notes, key.NoteNameInKey(current))
		semis := SemitoneDiff(root, current)
		degrees = append(degrees, semitoneToDegree[semis])
	}

	rootName := key.NoteNameInKey(root)
	return ScaleResult{
		Name:    rootName + " " + st.Name,
		Notes:   notes,
		Degrees: degrees,
		Display: strings.Join(notes, " - "),
	}
}

// GenerateAllScales produces all scale types for a given root note.
func GenerateAllScales(root Note, key Key) []ScaleResult {
	return GenerateScalesByType(root, key, nil)
}

// GenerateScalesByType produces scales filtered by type name (case-insensitive).
// If typeFilter is nil or empty, all scale types are returned.
func GenerateScalesByType(root Note, key Key, typeFilter []string) []ScaleResult {
	types := ScaleTypes
	if len(typeFilter) > 0 {
		filterSet := make(map[string]bool, len(typeFilter))
		for _, t := range typeFilter {
			filterSet[strings.ToLower(strings.TrimSpace(t))] = true
		}
		types = nil
		for _, st := range ScaleTypes {
			if filterSet[strings.ToLower(st.Name)] {
				types = append(types, st)
			}
		}
	}
	results := make([]ScaleResult, len(types))
	for i, st := range types {
		results[i] = GenerateScale(root, st, key)
	}
	return results
}
