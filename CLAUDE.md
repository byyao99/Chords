# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- Frontend dev: `cd frontend && npm install` (first time), then `npm run dev` (Vite on `:5173`, proxies `/api` to `:8080`). Run `go run .` in another shell so the proxy has a target.
- Frontend build: `cd frontend && npm run build` — outputs `frontend/dist/`. Required before `go build`, since `main.go` embeds `frontend/dist` via `//go:embed`.
- Production build: `cd frontend && npm run build && cd .. && go build -o chords .` — produces a single self-contained binary.
- Run server: `go run .` (listens on `$PORT`, default `8080`; serves the SPA + JSON APIs).
- Test all (Go): `go test ./...`
- Test one package: `go test ./music` or `go test ./handler`
- Test single function: `go test ./music -run TestParseChord`
- Verbose: `go test -v ./music -run TestParseChord`
- Format / vet: `go fmt ./...` and `go vet ./...`

The Vite build output (`frontend/dist`) is embedded via `//go:embed` in `main.go`, so the built binary is self-contained — no working-directory assumptions at runtime. `frontend/dist/` is gitignored; rebuilds must precede `go build`.

## Architecture

This is a Gin-based web app for guitar / music-theory tools (transpose, scales, diatonic chords, chord diagrams, progression generator, metronome). It is split into three layers:

### 1. `music/` — pure music-theory library

All domain logic lives here as a standalone package with no web dependencies. Everything is keyed off `Note` (an `int` 0–11 representing a pitch class) so transposition and enharmonic respelling are cheap. Key files and their roles:

- `note.go`, `intervals.go` — pitch-class primitives and interval formulas
- `key.go`, `scale.go` — `Key` struct carries an accidental preference (sharp/flat) that drives enharmonic spelling everywhere downstream; `GenerateAllScales` produces the scale cards shown on `/scales`
- `chord.go`, `parser.go` — `ParseChord` and `ParseSheet` read chord tokens (including slash chords like `C/E`) and full chord-sheet text; `ChordTones` returns the notes of a chord
- `diatonic.go` — `DiatonicChords(root, minor)` returns the 7 diatonic chords with roman-numeral degrees
- `progression.go` — `patternDatabase` of named progressions (tagged by mood/style) + `GenerateProgressions` which realizes them against a key
- `transpose.go` — `DetectKey`, `TransposeChord`, and `RenderSheet`. The transpose endpoint supports two modes: by semitone count (auto-detects source key) or by explicit from/to key
- `guitar.go` — hardcoded `chordVoicings` map of common open/barre shapes; `FindChordVoicing` looks up by chord name
- `voicing_gen.go` — algorithmically generates voicings on a standard-tuned fretboard (`standardTuning = E A D G B E`) and is used by `FindChordVoicings` when the static map lacks a match or more alternatives are requested

Every `.go` file in `music/` has a paired `_test.go` — run these when changing theory logic.

### 2. `handler/` — thin Gin layer

`handler.go` owns `SetupRouter(distFS fs.FS)`. The embedded `frontend/dist` filesystem is mounted at `/assets` (Vite's hashed bundles), and any non-API path is served the SPA's `index.html` via `r.NoRoute` so client-side routing works on refresh / direct navigation. Handlers are deliberately thin: parse query/JSON → call into `music` → shape into a JSON response. The one notable pattern is that diatonic and progression responses wrap each `DiatonicChord` with an optional `Voicing` field via an inline struct — keep this shape stable; `frontend/src/components/ChordDiagram.vue` depends on it.

Routes:
- `POST /api/transpose` — transpose a chord sheet
- `GET /api/scales`, `/api/diatonic`, `/api/chord-diagram`, `/api/progressions` — JSON APIs
- `/assets/*` — Vite hashed bundles
- All other paths fall through to `index.html` (SPA). `/api/*` paths that miss return 404 JSON, not the SPA.

### 3. `frontend/` — Vue 3 + TypeScript + Vite SPA

Single-page app with Vue Router. Layout: `App.vue` mounts `<SiteHeader>` + `<RouterView>`; pages live in `src/pages/`; reusable bits in `src/components/`; cross-cutting state in `src/composables/`. The previous vanilla-JS modules are now:

- `static/theme.js` → `src/composables/useTheme.ts` (also bootstrapped inline in `index.html` to avoid theme-flash before mount)
- `static/audio.js` → `src/composables/useAudio.ts` — module-scope `AudioContext` singleton (shared with `useMetronome`); never `close()` it across SPA navigation
- `static/metronome.js` → `src/composables/useMetronome.ts` — same look-ahead scheduler; pages use `useMetronomeAutoStop()` so navigating away calls `stop()` in `onBeforeUnmount`
- `static/diagram.js` → `src/components/ChordDiagram.vue` — declarative SVG renderer
- `static/app.js` and the inline scripts in each `templates/*.html` → `src/pages/*.vue`

Type definitions for the JSON API live in `src/types/api.ts` and mirror the Go structs in `handler/handler.go` — change them together. `src/constants/notes.ts` holds `ROOT_OPTIONS` (the 12 root-note options previously injected via `gin.H`), `NOTE_SEMI`, and `OPEN_FREQS`. `src/style.css` is the original monolithic stylesheet, imported once from `main.ts`; component-scoped styles can be migrated opportunistically.

### Spec document

`chord_service_spec_go_gin.md` describes an earlier, broader API design (progression generator, voice leading, analysis under `/api/v1/...`). The current code only partially implements this — the live routes are the ones in `handler.go`, not the spec. Treat the spec as aspirational context, not a source of truth.
