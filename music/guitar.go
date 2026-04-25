package music

// ChordVoicing represents a guitar chord fingering.
type ChordVoicing struct {
	Name     string `json:"name"`
	Frets    [6]int `json:"frets"`   // -1=muted, 0=open, 1+=fret
	Fingers  [6]int `json:"fingers"` // 0=not pressed, 1-4=finger
	BaseFret int    `json:"base_fret"`
}

// chordVoicings stores all known chord voicings, keyed by chord name.
var chordVoicings = map[string]ChordVoicing{
	// ── Open Major ──
	"C": {Name: "C", Frets: [6]int{-1, 3, 2, 0, 1, 0}, Fingers: [6]int{0, 3, 2, 0, 1, 0}, BaseFret: 1},
	"D": {Name: "D", Frets: [6]int{-1, -1, 0, 2, 3, 2}, Fingers: [6]int{0, 0, 0, 1, 3, 2}, BaseFret: 1},
	"E": {Name: "E", Frets: [6]int{0, 2, 2, 1, 0, 0}, Fingers: [6]int{0, 2, 3, 1, 0, 0}, BaseFret: 1},
	"F": {Name: "F", Frets: [6]int{1, 3, 3, 2, 1, 1}, Fingers: [6]int{1, 3, 4, 2, 1, 1}, BaseFret: 1},
	"G": {Name: "G", Frets: [6]int{3, 2, 0, 0, 0, 3}, Fingers: [6]int{2, 1, 0, 0, 0, 3}, BaseFret: 1},
	"A": {Name: "A", Frets: [6]int{-1, 0, 2, 2, 2, 0}, Fingers: [6]int{0, 0, 1, 2, 3, 0}, BaseFret: 1},
	"B": {Name: "B", Frets: [6]int{-1, 2, 4, 4, 4, 2}, Fingers: [6]int{0, 1, 2, 3, 4, 1}, BaseFret: 1},

	// ── Open Minor ──
	"Cm": {Name: "Cm", Frets: [6]int{-1, 3, 5, 5, 4, 3}, Fingers: [6]int{0, 1, 3, 4, 2, 1}, BaseFret: 1},
	"Dm": {Name: "Dm", Frets: [6]int{-1, -1, 0, 2, 3, 1}, Fingers: [6]int{0, 0, 0, 2, 3, 1}, BaseFret: 1},
	"Em": {Name: "Em", Frets: [6]int{0, 2, 2, 0, 0, 0}, Fingers: [6]int{0, 2, 3, 0, 0, 0}, BaseFret: 1},
	"Fm": {Name: "Fm", Frets: [6]int{1, 3, 3, 1, 1, 1}, Fingers: [6]int{1, 3, 4, 1, 1, 1}, BaseFret: 1},
	"Gm": {Name: "Gm", Frets: [6]int{3, 5, 5, 3, 3, 3}, Fingers: [6]int{1, 3, 4, 1, 1, 1}, BaseFret: 1},
	"Am": {Name: "Am", Frets: [6]int{-1, 0, 2, 2, 1, 0}, Fingers: [6]int{0, 0, 2, 3, 1, 0}, BaseFret: 1},
	"Bm": {Name: "Bm", Frets: [6]int{-1, 2, 4, 4, 3, 2}, Fingers: [6]int{0, 1, 3, 4, 2, 1}, BaseFret: 1},

	// ── 7th ──
	"C7": {Name: "C7", Frets: [6]int{-1, 3, 2, 3, 1, 0}, Fingers: [6]int{0, 3, 2, 4, 1, 0}, BaseFret: 1},
	"D7": {Name: "D7", Frets: [6]int{-1, -1, 0, 2, 1, 2}, Fingers: [6]int{0, 0, 0, 2, 1, 3}, BaseFret: 1},
	"E7": {Name: "E7", Frets: [6]int{0, 2, 0, 1, 0, 0}, Fingers: [6]int{0, 2, 0, 1, 0, 0}, BaseFret: 1},
	"F7": {Name: "F7", Frets: [6]int{1, 3, 1, 2, 1, 1}, Fingers: [6]int{1, 3, 1, 2, 1, 1}, BaseFret: 1},
	"G7": {Name: "G7", Frets: [6]int{3, 2, 0, 0, 0, 1}, Fingers: [6]int{3, 2, 0, 0, 0, 1}, BaseFret: 1},
	"A7": {Name: "A7", Frets: [6]int{-1, 0, 2, 0, 2, 0}, Fingers: [6]int{0, 0, 2, 0, 3, 0}, BaseFret: 1},
	"B7": {Name: "B7", Frets: [6]int{-1, 2, 1, 2, 0, 2}, Fingers: [6]int{0, 2, 1, 3, 0, 4}, BaseFret: 1},

	// ── Minor 7th ──
	"Dm7": {Name: "Dm7", Frets: [6]int{-1, -1, 0, 2, 1, 1}, Fingers: [6]int{0, 0, 0, 2, 1, 1}, BaseFret: 1},
	"Em7": {Name: "Em7", Frets: [6]int{0, 2, 0, 0, 0, 0}, Fingers: [6]int{0, 2, 0, 0, 0, 0}, BaseFret: 1},
	"Am7": {Name: "Am7", Frets: [6]int{-1, 0, 2, 0, 1, 0}, Fingers: [6]int{0, 0, 2, 0, 1, 0}, BaseFret: 1},
	"Bm7": {Name: "Bm7", Frets: [6]int{-1, 2, 0, 2, 0, 2}, Fingers: [6]int{0, 1, 0, 2, 0, 3}, BaseFret: 1},

	// ── Major 7th ──
	"Cmaj7": {Name: "Cmaj7", Frets: [6]int{-1, 3, 2, 0, 0, 0}, Fingers: [6]int{0, 3, 2, 0, 0, 0}, BaseFret: 1},
	"Dmaj7": {Name: "Dmaj7", Frets: [6]int{-1, -1, 0, 2, 2, 2}, Fingers: [6]int{0, 0, 0, 1, 2, 3}, BaseFret: 1},
	"Emaj7": {Name: "Emaj7", Frets: [6]int{0, 2, 1, 1, 0, 0}, Fingers: [6]int{0, 3, 1, 2, 0, 0}, BaseFret: 1},
	"Fmaj7": {Name: "Fmaj7", Frets: [6]int{-1, -1, 3, 2, 1, 0}, Fingers: [6]int{0, 0, 3, 2, 1, 0}, BaseFret: 1},
	"Gmaj7": {Name: "Gmaj7", Frets: [6]int{3, 2, 0, 0, 0, 2}, Fingers: [6]int{3, 2, 0, 0, 0, 1}, BaseFret: 1},
	"Amaj7": {Name: "Amaj7", Frets: [6]int{-1, 0, 2, 1, 2, 0}, Fingers: [6]int{0, 0, 2, 1, 3, 0}, BaseFret: 1},

	// ── Sus ──
	"Dsus2": {Name: "Dsus2", Frets: [6]int{-1, -1, 0, 2, 3, 0}, Fingers: [6]int{0, 0, 0, 1, 3, 0}, BaseFret: 1},
	"Dsus4": {Name: "Dsus4", Frets: [6]int{-1, -1, 0, 2, 3, 3}, Fingers: [6]int{0, 0, 0, 1, 2, 3}, BaseFret: 1},
	"Esus4": {Name: "Esus4", Frets: [6]int{0, 2, 2, 2, 0, 0}, Fingers: [6]int{0, 2, 3, 4, 0, 0}, BaseFret: 1},
	"Asus2": {Name: "Asus2", Frets: [6]int{-1, 0, 2, 2, 0, 0}, Fingers: [6]int{0, 0, 2, 3, 0, 0}, BaseFret: 1},
	"Asus4": {Name: "Asus4", Frets: [6]int{-1, 0, 2, 2, 3, 0}, Fingers: [6]int{0, 0, 1, 2, 3, 0}, BaseFret: 1},

	// ── Dim ──
	"Bdim":  {Name: "Bdim", Frets: [6]int{-1, 2, 3, 4, 3, -1}, Fingers: [6]int{0, 1, 2, 4, 3, 0}, BaseFret: 1},
	"C#dim": {Name: "C#dim", Frets: [6]int{-1, -1, 2, 3, 2, 3}, Fingers: [6]int{0, 0, 1, 3, 2, 4}, BaseFret: 1},
	"Ddim":  {Name: "Ddim", Frets: [6]int{-1, -1, 0, 1, 3, 1}, Fingers: [6]int{0, 0, 0, 1, 3, 2}, BaseFret: 1},
	"Edim":  {Name: "Edim", Frets: [6]int{0, 1, 2, 0, -1, -1}, Fingers: [6]int{0, 1, 2, 0, 0, 0}, BaseFret: 1},
	"F#dim": {Name: "F#dim", Frets: [6]int{2, 3, 4, 2, -1, -1}, Fingers: [6]int{1, 2, 3, 1, 0, 0}, BaseFret: 1},
	"G#dim": {Name: "G#dim", Frets: [6]int{4, 5, 6, 4, -1, -1}, Fingers: [6]int{1, 2, 3, 1, 0, 0}, BaseFret: 1},

	// ── Barre Major (E-shape) ──
	"C#": {Name: "C#", Frets: [6]int{-1, 4, 6, 6, 6, 4}, Fingers: [6]int{0, 1, 3, 3, 3, 1}, BaseFret: 1},
	"Eb": {Name: "Eb", Frets: [6]int{-1, -1, 1, 3, 4, 3}, Fingers: [6]int{0, 0, 1, 2, 4, 3}, BaseFret: 1},
	"F#": {Name: "F#", Frets: [6]int{2, 4, 4, 3, 2, 2}, Fingers: [6]int{1, 3, 4, 2, 1, 1}, BaseFret: 1},
	"Ab": {Name: "Ab", Frets: [6]int{4, 6, 6, 5, 4, 4}, Fingers: [6]int{1, 3, 4, 2, 1, 1}, BaseFret: 1},
	"Bb": {Name: "Bb", Frets: [6]int{-1, 1, 3, 3, 3, 1}, Fingers: [6]int{0, 1, 2, 3, 4, 1}, BaseFret: 1},

	// ── Barre Minor ──
	"C#m": {Name: "C#m", Frets: [6]int{-1, 4, 6, 6, 5, 4}, Fingers: [6]int{0, 1, 3, 4, 2, 1}, BaseFret: 1},
	"Ebm": {Name: "Ebm", Frets: [6]int{-1, -1, 1, 3, 4, 2}, Fingers: [6]int{0, 0, 1, 3, 4, 2}, BaseFret: 1},
	"F#m": {Name: "F#m", Frets: [6]int{2, 4, 4, 2, 2, 2}, Fingers: [6]int{1, 3, 4, 1, 1, 1}, BaseFret: 1},
	"G#m": {Name: "G#m", Frets: [6]int{4, 6, 6, 4, 4, 4}, Fingers: [6]int{1, 3, 4, 1, 1, 1}, BaseFret: 1},
	"Bbm": {Name: "Bbm", Frets: [6]int{-1, 1, 3, 3, 2, 1}, Fingers: [6]int{0, 1, 3, 4, 2, 1}, BaseFret: 1},
}

// Enharmonic aliases for chord lookup
var chordAliases = map[string]string{
	"Db":    "C#",
	"D#":    "Eb",
	"Gb":    "F#",
	"G#":    "Ab",
	"A#":    "Bb",
	"Dbm":   "C#m",
	"D#m":   "Ebm",
	"Gbm":   "F#m",
	"Abm":   "G#m",
	"A#m":   "Bbm",
	"Dbdim": "C#dim",
	"Abdim": "G#dim",
	"Gbdim": "F#dim",
}

// FindChordVoicing looks up a chord voicing by name.
// Handles enharmonic aliases (e.g., "Db" -> "C#").
func FindChordVoicing(name string) (ChordVoicing, bool) {
	if v, ok := chordVoicings[name]; ok {
		return v, true
	}
	if alias, ok := chordAliases[name]; ok {
		if v, ok := chordVoicings[alias]; ok {
			v.Name = name // display with the requested name
			return v, true
		}
	}
	// Fall back to algorithmic generation.
	chord, ok := ParseChord(name)
	if !ok {
		return ChordVoicing{}, false
	}
	v, ok := GenerateVoicing(chord)
	if ok {
		v.Name = name
	}
	return v, ok
}

// FindChordVoicings returns up to `limit` voicings for the given chord name,
// only including voicings where all pressed frets are <= maxFret.
// The hardcoded voicing (if any) is always first.
func FindChordVoicings(name string, maxFret int, limit int) []ChordVoicing {
	var result []ChordVoicing
	seen := map[[6]int]bool{}

	// Try hardcoded lookup first.
	var hardcoded *ChordVoicing
	if v, ok := chordVoicings[name]; ok {
		hardcoded = &v
	} else if alias, ok := chordAliases[name]; ok {
		if v, ok := chordVoicings[alias]; ok {
			v.Name = name
			hardcoded = &v
		}
	}
	if hardcoded != nil && fitsMaxFret(hardcoded.Frets, maxFret) {
		result = append(result, *hardcoded)
		seen[hardcoded.Frets] = true
	}

	// Fill remaining slots with algorithmically generated voicings.
	if len(result) < limit {
		chord, ok := ParseChord(name)
		if ok {
			generated := GenerateVoicings(chord, maxFret, limit*2)
			for _, v := range generated {
				if len(result) >= limit {
					break
				}
				if seen[v.Frets] {
					continue
				}
				seen[v.Frets] = true
				v.Name = name
				result = append(result, v)
			}
		}
	}
	return result
}

// fitsMaxFret reports whether all pressed frets in v are <= maxFret.
func fitsMaxFret(frets [6]int, maxFret int) bool {
	for _, f := range frets {
		if f > maxFret {
			return false
		}
	}
	return true
}

// AllChordNames returns all available chord voicing names (sorted would be nice, but map order is fine for now).
func AllChordNames() []string {
	names := make([]string, 0, len(chordVoicings))
	for name := range chordVoicings {
		names = append(names, name)
	}
	return names
}
