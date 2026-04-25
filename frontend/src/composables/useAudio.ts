import { NOTE_SEMI, OPEN_FREQS } from '../constants/notes'

// Module-scope singleton. Shared with useMetronome via getAudioContext().
let ctx: AudioContext | null = null

export function getAudioContext(): AudioContext {
  if (!ctx) {
    const Ctor = window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext
    ctx = new Ctor()
  }
  if (ctx.state === 'suspended') void ctx.resume()
  return ctx
}

function playFreq(freq: number, startDelay: number) {
  const ac = getAudioContext()
  const t = ac.currentTime + startDelay

  const master = ac.createGain()
  master.gain.setValueAtTime(0, t)
  master.gain.linearRampToValueAtTime(0.38, t + 0.005)
  master.gain.exponentialRampToValueAtTime(0.001, t + 1.2)

  const filter = ac.createBiquadFilter()
  filter.type = 'lowpass'
  filter.frequency.setValueAtTime(4000, t)
  filter.frequency.exponentialRampToValueAtTime(600, t + 0.6)
  filter.Q.value = 1

  filter.connect(master)
  master.connect(ac.destination)

  const osc1 = ac.createOscillator()
  const g1 = ac.createGain()
  osc1.type = 'triangle'
  osc1.frequency.setValueAtTime(freq, t)
  g1.gain.value = 0.55
  osc1.connect(g1)
  g1.connect(filter)
  osc1.start(t)
  osc1.stop(t + 1.2)

  const osc2 = ac.createOscillator()
  const g2 = ac.createGain()
  osc2.type = 'sine'
  osc2.frequency.setValueAtTime(freq * 2, t)
  g2.gain.value = 0.25
  osc2.connect(g2)
  g2.connect(filter)
  osc2.start(t)
  osc2.stop(t + 1.2)

  const osc3 = ac.createOscillator()
  const g3 = ac.createGain()
  osc3.type = 'sine'
  osc3.frequency.setValueAtTime(freq * 3, t)
  g3.gain.setValueAtTime(0.15, t)
  g3.gain.exponentialRampToValueAtTime(0.001, t + 0.4)
  osc3.connect(g3)
  g3.connect(filter)
  osc3.start(t)
  osc3.stop(t + 1.2)
}

function noteFreq(name: string, octave: number): number | null {
  const semi = NOTE_SEMI[name]
  if (semi === undefined) return null
  const midi = 12 * (octave + 1) + semi
  return 440 * Math.pow(2, (midi - 69) / 12)
}

function fretFreq(stringIndex: number, fret: number): number {
  return OPEN_FREQS[stringIndex] * Math.pow(2, fret / 12)
}

export function playNoteByName(name: string, octave = 4) {
  const freq = noteFreq(name, octave)
  if (freq) playFreq(freq, 0)
}

export function playChordVoicing(frets: number[]) {
  let delay = 0
  for (let s = 0; s < 6; s++) {
    if (frets[s] === -1) continue
    playFreq(fretFreq(s, frets[s]), delay)
    delay += 0.05
  }
}

// Returns the octave to use when playing scale notes — bumps to 5 when the
// note semitone is below the root's, so the scale plays ascending.
// Mirrors templates/scales.html:146.
export function octaveForScaleNote(rootName: string, noteName: string): number {
  const rootSemi = NOTE_SEMI[rootName]
  const semi = NOTE_SEMI[noteName]
  if (semi !== undefined && rootSemi !== undefined && semi < rootSemi) return 5
  return 4
}
