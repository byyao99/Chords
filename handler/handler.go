package handler

import (
	"chords/music"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TransposeRequest struct {
	Sheet     string `json:"sheet"`
	Mode      string `json:"mode"`      // "semitone" or "key"
	Semitones int    `json:"semitones"` // used when mode = "semitone"
	FromKey   string `json:"from_key"`  // used when mode = "key"
	ToKey     string `json:"to_key"`    // used when mode = "key"
}

type TransposeResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// SetupRouter creates the gin engine. distFS is the embedded Vite build output
// containing index.html and an assets/ subdirectory.
func SetupRouter(distFS fs.FS) (*gin.Engine, error) {
	r := gin.Default()

	indexHTML, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		return nil, err
	}

	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		return nil, err
	}
	r.StaticFS("/assets", http.FS(assetsFS))

	r.POST("/api/transpose", handleTranspose)
	r.GET("/api/scales", handleScalesAPI)
	r.GET("/api/diatonic", handleDiatonicAPI)
	r.GET("/api/chord-diagram", handleChordDiagramAPI)
	r.GET("/api/progressions", handleProgressionsAPI)

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	return r, nil
}

type ScaleCard struct {
	Name    string   `json:"name"`
	Notes   []string `json:"notes"`
	Degrees []string `json:"degrees"`
	Display string   `json:"display"`
}

func handleScalesAPI(c *gin.Context) {
	rootName := c.DefaultQuery("root", "C")
	key, err := music.KeyFromName(rootName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid root: " + rootName})
		return
	}
	var typeFilter []string
	if typeParam := c.Query("type"); typeParam != "" {
		typeFilter = strings.Split(typeParam, ",")
	}
	scales := music.GenerateScalesByType(key.Root, key, typeFilter)

	cards := make([]ScaleCard, len(scales))
	for i, s := range scales {
		cards[i] = ScaleCard{Name: s.Name, Notes: s.Notes, Degrees: s.Degrees, Display: s.Display}
	}

	c.JSON(http.StatusOK, gin.H{"scales": cards})
}

func handleTranspose(c *gin.Context) {
	var req TransposeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TransposeResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	if req.Sheet == "" {
		c.JSON(http.StatusBadRequest, TransposeResponse{Error: "sheet is required"})
		return
	}

	sheet := music.ParseSheet(req.Sheet)

	var semitones int
	var key music.Key

	switch req.Mode {
	case "semitone":
		semitones = req.Semitones
		if semitones == 0 {
			c.JSON(http.StatusOK, TransposeResponse{Result: req.Sheet})
			return
		}
		detected := music.DetectKey(sheet)
		key = music.KeyFromNote(detected.Root.Transpose(semitones))

	case "key":
		fromKey, err := music.KeyFromName(req.FromKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, TransposeResponse{Error: "invalid from_key: " + req.FromKey})
			return
		}
		toKey, err := music.KeyFromName(req.ToKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, TransposeResponse{Error: "invalid to_key: " + req.ToKey})
			return
		}
		semitones = music.SemitoneDiff(fromKey.Root, toKey.Root)
		key = toKey

	default:
		c.JSON(http.StatusBadRequest, TransposeResponse{Error: "mode must be 'semitone' or 'key'"})
		return
	}

	result := music.RenderSheet(sheet, semitones, key)
	c.JSON(http.StatusOK, TransposeResponse{Result: result})
}

// ── Diatonic Chords ──

func handleDiatonicAPI(c *gin.Context) {
	rootName := c.DefaultQuery("root", "C")
	minor := c.DefaultQuery("minor", "false") == "true"

	if _, err := music.ParseNoteName(rootName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid root: " + rootName})
		return
	}

	chords := music.DiatonicChords(rootName, minor)

	type DiatonicWithVoicing struct {
		music.DiatonicChord
		Voicing *music.ChordVoicing `json:"voicing,omitempty"`
	}

	result := make([]DiatonicWithVoicing, len(chords))
	for i, dc := range chords {
		result[i].DiatonicChord = dc
		if v, ok := music.FindChordVoicing(dc.Name); ok {
			result[i].Voicing = &v
		}
	}

	label := rootName + " Major"
	if minor {
		label = rootName + " Minor"
	}

	c.JSON(http.StatusOK, gin.H{
		"key":    label,
		"chords": result,
	})
}

// ── Chord Diagrams ──

func handleChordDiagramAPI(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name parameter is required"})
		return
	}

	maxFret := 12
	if s := c.Query("max_fret"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			maxFret = n
		}
	}
	limit := 3
	if s := c.Query("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			limit = n
		}
	}

	voicings := music.FindChordVoicings(name, maxFret, limit)
	if len(voicings) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "chord not found: " + name})
		return
	}

	notes, _ := music.ChordTones(name)
	c.JSON(http.StatusOK, gin.H{
		"voicings": voicings,
		"notes":    notes,
	})
}

// ── Chord Progression Generator ──

func handleProgressionsAPI(c *gin.Context) {
	rootName := c.DefaultQuery("root", "C")
	minor := c.DefaultQuery("minor", "false") == "true"
	mood := c.DefaultQuery("mood", "")
	style := c.DefaultQuery("style", "")
	length := 0
	if s := c.Query("length"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			length = n
		}
	}

	if _, err := music.ParseNoteName(rootName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid root: " + rootName})
		return
	}

	suggestions, err := music.GenerateProgressions(rootName, minor, mood, style, length)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type ChordWithVoicing struct {
		music.DiatonicChord
		Voicing *music.ChordVoicing `json:"voicing,omitempty"`
	}
	type ProgressionWithVoicing struct {
		PatternName string             `json:"pattern_name"`
		Mood        string             `json:"mood"`
		Style       string             `json:"style"`
		Chords      []ChordWithVoicing `json:"chords"`
	}

	result := make([]ProgressionWithVoicing, len(suggestions))
	for i, sug := range suggestions {
		chords := make([]ChordWithVoicing, len(sug.Chords))
		for j, dc := range sug.Chords {
			chords[j].DiatonicChord = dc
			if v, ok := music.FindChordVoicing(dc.Name); ok {
				chords[j].Voicing = &v
			}
		}
		result[i] = ProgressionWithVoicing{
			PatternName: sug.PatternName,
			Mood:        sug.Mood,
			Style:       sug.Style,
			Chords:      chords,
		}
	}

	keyLabel := rootName + " Major"
	if minor {
		keyLabel = rootName + " Minor"
	}
	c.JSON(http.StatusOK, gin.H{
		"key":          keyLabel,
		"progressions": result,
	})
}
