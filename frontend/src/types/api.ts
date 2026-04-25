// matches handler.TransposeRequest
export interface TransposeRequest {
  sheet: string
  mode: 'semitone' | 'key'
  semitones?: number
  from_key?: string
  to_key?: string
}

// matches handler.TransposeResponse
export interface TransposeResponse {
  result: string
  error?: string
}

// matches handler.ScaleCard
export interface ScaleCard {
  name: string
  notes: string[]
  degrees: string[]
  display: string
}

export interface ScalesResponse {
  scales: ScaleCard[]
  error?: string
}

// matches music.ChordVoicing
export interface ChordVoicing {
  name: string
  frets: number[]
  fingers: number[]
  base_fret: number
}

// matches music.DiatonicChord with optional voicing attached by handler
export interface DiatonicChord {
  degree: string
  name: string
  quality: string
  voicing?: ChordVoicing
}

export interface DiatonicResponse {
  key: string
  chords: DiatonicChord[]
  error?: string
}

export interface ChordDiagramResponse {
  voicings: ChordVoicing[]
  notes: string[]
  error?: string
}

export interface Progression {
  pattern_name: string
  mood: string
  style: string
  chords: DiatonicChord[]
}

export interface ProgressionsResponse {
  key: string
  progressions: Progression[]
  error?: string
}
