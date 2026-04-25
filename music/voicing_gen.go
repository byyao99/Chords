package music

// standardTuning holds the pitch class (0-11) of each open string,
// from low E to high E: E A D G B E.
var standardTuning = [6]Note{E, A, D, G, B, E}

// GenerateVoicing algorithmically finds the best playable guitar voicing for
// the given chord. It returns false if no valid voicing can be found.
func GenerateVoicing(chord Chord) (ChordVoicing, bool) {
	vs := GenerateVoicings(chord, 12, 1)
	if len(vs) == 0 {
		return ChordVoicing{}, false
	}
	return vs[0], true
}

// GenerateVoicings returns up to `limit` distinct voicings for the given chord.
// Only voicings where all pressed frets are <= maxFret are included.
func GenerateVoicings(chord Chord, maxFret int, limit int) []ChordVoicing {
	formula, ok := QualityIntervals(chord.Quality)
	if !ok {
		return nil
	}

	required := formula.Required
	if required == nil {
		required = formula.Intervals
	}
	requiredSet := toMod12Set(required)
	allSet := toMod12Set(formula.Intervals)

	bassNote := chord.Root
	if chord.BassNote != nil {
		bassNote = *chord.BassNote
	}
	bassPC := int(bassNote) % 12

	candidates := collectVoicings(chord, allSet, requiredSet, bassPC, maxFret)
	// Sort by score ascending.
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].score < candidates[i].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	name := chord.Root.SharpName() + chord.Quality
	var result []ChordVoicing
	seen := map[[6]int]bool{}
	for _, c := range candidates {
		if len(result) >= limit {
			break
		}
		if seen[c.frets] {
			continue
		}
		seen[c.frets] = true
		result = append(result, ChordVoicing{
			Name:     name,
			Frets:    c.frets,
			Fingers:  assignFingers(c.frets),
			BaseFret: calcBaseFret(c.frets),
		})
	}
	return result
}

// toMod12Set converts a slice of intervals (possibly > 11) to a set of pitch
// classes mod 12, offset from the root's absolute pitch class is handled at
// call sites. Here we just store the reduced intervals.
func toMod12Set(intervals []int) map[int]bool {
	s := make(map[int]bool, len(intervals))
	for _, iv := range intervals {
		s[iv%12] = true
	}
	return s
}

type candidate struct {
	frets [6]int
	score float64
}

func collectVoicings(chord Chord, allSet, requiredSet map[int]bool, bassPC, maxFret int) []candidate {
	rootPC := int(chord.Root) % 12
	var all []candidate

	// Slide the window from open position (base=0) through maxFret.
	for base := 0; base <= maxFret; base++ {
		// Build per-string candidate frets.
		var perString [6][]int
		for s := 0; s < 6; s++ {
			openPC := int(standardTuning[s]) % 12
			var cands []int
			cands = append(cands, -1) // muted

			// Open string.
			interval := (openPC - rootPC + 12) % 12
			if allSet[interval] {
				cands = append(cands, 0)
			}

			// Frets in window.
			start := base
			if start == 0 {
				start = 1
			}
			for f := start; f <= base+3 && f <= maxFret; f++ {
				pc := (openPC + f) % 12
				iv := (pc - rootPC + 12) % 12
				if allSet[iv] {
					cands = append(cands, f)
				}
			}
			perString[s] = cands
		}

		enumerateAll(perString, [6]int{}, 0, rootPC, bassPC, maxFret, requiredSet, &all)
	}
	return all
}

func enumerateAll(
	perString [6][]int,
	current [6]int,
	s int,
	rootPC, bassPC, maxFret int,
	requiredSet map[int]bool,
	all *[]candidate,
) {
	if s == 6 {
		c := scoreCandidate(current, rootPC, bassPC, requiredSet)
		if c == nil {
			return
		}
		*all = append(*all, *c)
		return
	}
	for _, f := range perString[s] {
		current[s] = f
		enumerateAll(perString, current, s+1, rootPC, bassPC, maxFret, requiredSet, all)
	}
}

func scoreCandidate(frets [6]int, rootPC, bassPC int, requiredSet map[int]bool) *candidate {
	// Collect sounding strings and their pitch classes.
	type stringInfo struct {
		str int
		pc  int
	}
	var sounding []stringInfo
	for s := 0; s < 6; s++ {
		if frets[s] == -1 {
			continue
		}
		pc := (int(standardTuning[s]) + frets[s]) % 12
		sounding = append(sounding, stringInfo{s, pc})
	}

	// Must have at least 3 sounding strings.
	if len(sounding) < 3 {
		return nil
	}

	// Lowest sounding string must produce bassPC.
	if sounding[0].pc != bassPC {
		return nil
	}

	// Check all required intervals are present.
	present := make(map[int]bool, len(sounding))
	for _, si := range sounding {
		iv := (si.pc - rootPC + 12) % 12
		present[iv] = true
	}
	for iv := range requiredSet {
		if !present[iv] {
			return nil
		}
	}

	// Fret span check.
	minF, maxF := 999, 0
	for _, f := range frets {
		if f > 0 {
			if f < minF {
				minF = f
			}
			if f > maxF {
				maxF = f
			}
		}
	}
	span := 0
	if minF != 999 {
		span = maxF - minF
	}
	if span > 4 {
		return nil
	}

	// Scoring (lower is better).
	score := float64(span) * 10

	openCount := 0
	muteCount := 0
	for s, f := range frets {
		if f == 0 {
			openCount++
		} else if f == -1 {
			muteCount++
			// Interior muted string: heavier penalty.
			if s > 0 && s < 5 && frets[s-1] != -1 && frets[s+1] != -1 {
				score += 8
			}
		}
	}
	score -= float64(openCount) * 3
	score += float64(muteCount) * 5

	// Reward root in bass.
	if sounding[0].pc == rootPC {
		score -= 15
	}

	// Prefer low positions.
	if minF != 999 {
		score += float64(minF) * 0.5
	}

	// Detect barre (same fret on consecutive strings).
	if minF != 999 {
		count := 0
		for _, f := range frets {
			if f == minF {
				count++
			}
		}
		if count >= 2 {
			score -= 5
		}
	}

	return &candidate{frets: frets, score: score}
}

// assignFingers heuristically assigns finger numbers (1-4) to pressed frets.
func assignFingers(frets [6]int) [6]int {
	var fingers [6]int

	// Find the minimum pressed fret — that's the barre fret (finger 1).
	minF := 999
	for _, f := range frets {
		if f > 0 && f < minF {
			minF = f
		}
	}
	if minF == 999 {
		return fingers // all open/muted
	}

	// Assign finger 1 to all strings at minF (barre or single).
	for s, f := range frets {
		if f == minF {
			fingers[s] = 1
		}
	}

	// Collect remaining fret positions above minF, sorted low to high.
	type fretPos struct{ fret, str int }
	var remaining []fretPos
	for s, f := range frets {
		if f > minF {
			remaining = append(remaining, fretPos{f, s})
		}
	}
	// Sort by fret then string.
	for i := 0; i < len(remaining); i++ {
		for j := i + 1; j < len(remaining); j++ {
			if remaining[j].fret < remaining[i].fret ||
				(remaining[j].fret == remaining[i].fret && remaining[j].str < remaining[i].str) {
				remaining[i], remaining[j] = remaining[j], remaining[i]
			}
		}
	}

	// Assign fingers 2, 3, 4.
	finger := 2
	seen := map[int]int{} // fret -> finger already assigned
	for _, fp := range remaining {
		if f, ok := seen[fp.fret]; ok {
			fingers[fp.str] = f
		} else {
			if finger <= 4 {
				fingers[fp.str] = finger
				seen[fp.fret] = finger
				finger++
			}
		}
	}

	return fingers
}

// calcBaseFret returns the base fret for diagram rendering.
// Follows the same convention as the hardcoded voicings: always 1,
// since diagram.js handles the visual offset from the fret numbers.
func calcBaseFret(frets [6]int) int {
	return 1
}
