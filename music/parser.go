package music

import "strings"

// ChordOccurrence represents a chord found at a specific position in a line.
type ChordOccurrence struct {
	Chord    Chord
	StartCol int // 0-based start column
	EndCol   int // exclusive end column
	Original string
}

// Line represents a parsed line from a chord sheet.
type Line struct {
	IsChordLine bool
	Raw         string
	Chords      []ChordOccurrence // populated only for chord lines
}

// Sheet represents a parsed chord sheet.
type Sheet struct {
	Lines []Line
}

// ParseSheet parses a chord sheet text into structured lines.
func ParseSheet(text string) Sheet {
	rawLines := strings.Split(text, "\n")
	sheet := Sheet{Lines: make([]Line, len(rawLines))}

	for i, raw := range rawLines {
		sheet.Lines[i] = parseLine(raw)
	}

	return sheet
}

func parseLine(raw string) Line {
	line := Line{Raw: raw}

	if strings.TrimSpace(raw) == "" {
		return line
	}

	// Split into tokens preserving positions
	type token struct {
		text     string
		startCol int
	}

	var tokens []token
	inToken := false
	start := 0

	for i, ch := range raw {
		if ch == ' ' || ch == '\t' {
			if inToken {
				tokens = append(tokens, token{raw[start:i], start})
				inToken = false
			}
		} else {
			if !inToken {
				start = i
				inToken = true
			}
		}
	}
	if inToken {
		tokens = append(tokens, token{raw[start:], start})
	}

	if len(tokens) == 0 {
		return line
	}

	// Try to parse each token as a chord
	chordCount := 0
	var occurrences []ChordOccurrence

	for _, tok := range tokens {
		chord, ok := ParseChord(tok.text)
		if ok {
			chordCount++
			occurrences = append(occurrences, ChordOccurrence{
				Chord:    chord,
				StartCol: tok.startCol,
				EndCol:   tok.startCol + len(tok.text),
				Original: tok.text,
			})
		}
	}

	// A line is a chord line if at least 50% of tokens are chords and there's at least one chord
	ratio := float64(chordCount) / float64(len(tokens))
	if chordCount > 0 && ratio >= 0.5 {
		line.IsChordLine = true
		line.Chords = occurrences
	}

	return line
}

// RenderSheet renders a transposed sheet as a string.
// semitones is the transposition offset.
// key determines enharmonic spelling.
func RenderSheet(sheet Sheet, semitones int, key Key) string {
	var lines []string

	for _, line := range sheet.Lines {
		if !line.IsChordLine {
			lines = append(lines, line.Raw)
			continue
		}
		lines = append(lines, renderChordLine(line, semitones, key))
	}

	return strings.Join(lines, "\n")
}

func renderChordLine(line Line, semitones int, key Key) string {
	if len(line.Chords) == 0 {
		return line.Raw
	}

	// Build the output line by placing transposed chords at their original positions
	var buf []byte
	cursor := 0

	for _, occ := range line.Chords {
		transposed := TransposeChord(occ.Chord, semitones, key)
		newName := transposed.Format(key)

		// Pad with spaces to reach the chord's start column
		for cursor < occ.StartCol {
			buf = append(buf, ' ')
			cursor++
		}

		// If cursor is past the start column (previous chord was longer), add at least one space
		if cursor > occ.StartCol {
			buf = append(buf, ' ')
			cursor++
		}

		buf = append(buf, []byte(newName)...)
		cursor += len(newName)
	}

	return string(buf)
}

// DetectKey attempts to detect the key of a chord sheet from its first chord.
func DetectKey(sheet Sheet) Key {
	for _, line := range sheet.Lines {
		if line.IsChordLine && len(line.Chords) > 0 {
			root := line.Chords[0].Chord.Root
			return KeyFromNote(root)
		}
	}
	return KeyFromNote(C) // default
}
