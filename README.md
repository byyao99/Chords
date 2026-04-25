# Chords

A guitar and music theory toolkit — chord transposition, scale explorer, diatonic chord lookup, chord diagrams, progression generator, and metronome — served as a single self-contained binary.

## Features

| Tool | Description |
|------|-------------|
| **Transpose** | Paste a chord sheet and shift by semitones or between keys |
| **Chord Lookup** | Find guitar voicings by name (open chords, barre shapes, alternatives) |
| **Diatonic Chords** | All 7 diatonic chords for any major or minor key, with guitar voicings |
| **Chord Progressions** | Common progressions filtered by mood and style (pop, rock, jazz, blues, folk) |
| **Scales** | Scale notes and degrees for any root, filterable by scale type |
| **Metronome** | BPM slider, tap tempo, adjustable time signature, per-beat accent toggle |

## Tech Stack

- **Backend**: Go + [Gin](https://github.com/gin-gonic/gin); pure music-theory library in `music/`
- **Frontend**: Vue 3 + TypeScript + Vite; single-page app with Vue Router
- **Deployment**: Single binary via `//go:embed all:frontend/dist` — no runtime dependencies

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+

### Development

```bash
# Terminal A — Go backend on :8080
go run .

# Terminal B — Vite dev server on :5173 (proxies /api to :8080)
cd frontend && npm install && npm run dev
```

Open http://localhost:5173.

### Production Build

```bash
cd frontend && npm run build   # outputs frontend/dist/
cd .. && go build -o chords .  # embeds dist/ into the binary
./chords                        # serves on :8080 (set PORT to override)
```

## Project Structure

```
chords/
├── main.go                 Entry point; embeds frontend/dist
├── handler/
│   ├── handler.go          Gin router — JSON APIs + SPA fallback
│   └── handler_test.go
├── music/                  Pure music-theory library (no web deps)
│   ├── note.go, key.go, intervals.go, scale.go
│   ├── chord.go, parser.go, transpose.go, diatonic.go
│   ├── progression.go, guitar.go, voicing_gen.go
│   └── *_test.go
└── frontend/
    └── src/
        ├── pages/          One .vue file per tool
        ├── components/     ChordDiagram.vue, ChordCard.vue, SiteHeader.vue
        ├── composables/    useAudio.ts, useMetronome.ts, useTheme.ts, useApi.ts
        ├── types/api.ts    TypeScript interfaces mirroring Go response structs
        └── constants/      ROOT_OPTIONS, NOTE_SEMI, OPEN_FREQS
```

## API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/transpose` | Transpose a chord sheet |
| `GET` | `/api/scales` | Scale cards for a root note |
| `GET` | `/api/diatonic` | Diatonic chords for a key |
| `GET` | `/api/chord-diagram` | Guitar voicings for a chord name |
| `GET` | `/api/progressions` | Chord progressions with optional mood/style filter |

## Running Tests

```bash
go test ./...          # all Go tests
go test ./music        # music theory library only
go test ./handler      # HTTP handler tests
```
