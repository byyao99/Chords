# Chord Service Spec (v1)

## Overview

This service provides three core features:

1. Chord Progression Generator
2. Voice Leading Suggestion
3. Chord Progression Analysis

Base URL:

```
/api/v1
```

---

# 1. Chord Progression Generator

## Description
Generate chord progressions based on key, scale, mood, and style.

## Endpoint

```
POST /progressions/generate
```

## Request

```json
{
  "key": "C",
  "scale": "major",
  "mood": "happy",
  "style": "pop",
  "complexity": "simple",
  "length": 4
}
```

## Parameters

| Field       | Type   | Required | Description |
|------------|--------|----------|------------|
| key        | string | yes      | Musical key (C, G, F#...) |
| scale      | string | yes      | major / minor |
| mood       | string | no       | happy / sad / emotional / dark |
| style      | string | no       | pop / rock / jazz / neo-soul |
| complexity | string | no       | simple / medium / complex |
| length     | int    | no       | number of chords |

## Response

```json
{
  "key": "C",
  "scale": "major",
  "progressions": [
    {
      "name": "Pop Classic",
      "roman": ["I", "V", "vi", "IV"],
      "chords": ["C", "G", "Am", "F"]
    }
  ]
}
```

## Implementation Notes

- Maintain a progression pattern database

```go
type ProgressionPattern struct {
    Name       string
    Roman      []string
    Tags       []string
    Complexity string
}
```

- Map roman numerals to chords based on key

---

# 2. Voice Leading Suggestion

## Description
Optimize chord transitions by minimizing pitch movement.

## Endpoint

```
POST /voice-leading
```

## Request

```json
{
  "chords": ["C", "G", "Am", "F"],
  "instrument": "piano"
}
```

## Response

```json
{
  "input": ["C", "G", "Am", "F"],
  "suggestions": [
    {
      "chords": ["C", "G/B", "Am", "F"],
      "voicing": [
        ["C", "E", "G"],
        ["B", "D", "G"],
        ["A", "C", "E"],
        ["A", "C", "F"]
      ],
      "movement_score": 2.1
    }
  ]
}
```

## Implementation Steps

1. Convert chords to notes
2. Generate inversions
3. Compute movement cost
4. Select optimal transitions

```go
func voiceLeadingCost(a, b []Note) float64 {
    // calculate minimal pitch distance
}
```

---

# 3. Chord Progression Analysis

## Description
Analyze chord input to determine key, roman numerals, and style.

## Endpoint

```
POST /analysis
```

## Request

```json
{
  "input": "C G Am F"
}
```

## Response

```json
{
  "key": "C",
  "scale": "major",
  "roman": ["I", "V", "vi", "IV"],
  "style": ["pop"],
  "borrowed_chords": [],
  "secondary_dominants": [],
  "confidence": 0.92
}
```

## Implementation Steps

1. Parse chord string into array
2. Try all keys (major/minor)
3. Score matches
4. Map to roman numerals
5. Detect patterns and styles
6. Identify special chords

---

# Shared Models

```go
type Chord struct {
    Root string
    Type string
}

type AnalysisResult struct {
    Key        string
    Scale      string
    Roman      []string
    Style      []string
    Confidence float64
}
```

---

# Gin Router Setup

```go
r := gin.Default()

api := r.Group("/api/v1")
{
    api.POST("/progressions/generate", GenerateProgression)
    api.POST("/voice-leading", VoiceLeading)
    api.POST("/analysis", AnalyzeChords)
}
```

---

# Suggested Enhancements

- Add caching (Redis)
- Keep core logic as pure functions
- Design for future AI prompt-based input

Example:

```json
{
  "prompt": "sad japanese chord progression"
}
```

