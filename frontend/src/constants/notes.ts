export interface RootOption {
  value: string
  label: string
}

// Used by Diatonic + Scales pages (mirrors handler.rootOptions).
export const ROOT_OPTIONS: RootOption[] = [
  { value: 'C', label: 'C' },
  { value: 'C#', label: 'C# / Db' },
  { value: 'D', label: 'D' },
  { value: 'Eb', label: 'Eb / D#' },
  { value: 'E', label: 'E' },
  { value: 'F', label: 'F' },
  { value: 'F#', label: 'F# / Gb' },
  { value: 'G', label: 'G' },
  { value: 'Ab', label: 'Ab / G#' },
  { value: 'A', label: 'A' },
  { value: 'Bb', label: 'Bb / A#' },
  { value: 'B', label: 'B' },
]

// Used by Chord lookup page (sharp-only spellings since the API parses these).
export const CHROMATIC_ROOTS: RootOption[] = [
  { value: 'C', label: 'C' },
  { value: 'C#', label: 'C# / Db' },
  { value: 'D', label: 'D' },
  { value: 'D#', label: 'D# / Eb' },
  { value: 'E', label: 'E' },
  { value: 'F', label: 'F' },
  { value: 'F#', label: 'F# / Gb' },
  { value: 'G', label: 'G' },
  { value: 'G#', label: 'G# / Ab' },
  { value: 'A', label: 'A' },
  { value: 'A#', label: 'A# / Bb' },
  { value: 'B', label: 'B' },
]

// Note name → semitone offset from C.
export const NOTE_SEMI: Record<string, number> = {
  C: 0, 'C#': 1, Db: 1, D: 2, 'D#': 3, Eb: 3,
  E: 4, F: 5, 'F#': 6, Gb: 6, G: 7, 'G#': 8,
  Ab: 8, A: 9, 'A#': 10, Bb: 10, B: 11,
}

// Standard guitar tuning open-string frequencies (Hz).
export const OPEN_FREQS: number[] = [82.41, 110.0, 146.83, 196.0, 246.94, 329.63]
