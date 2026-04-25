import { onBeforeUnmount, reactive, ref } from 'vue'
import { getAudioContext } from './useAudio'

type BeatCallback = (beatIndex: number, isAccent: boolean) => void

const running = ref(false)
const bpm = ref(120)
const beats = ref(4)
const accents = reactive<Record<number, boolean>>({ 0: true })

let currentBeat = 0
let nextNoteTime = 0
const scheduleAhead = 0.1 // seconds
const lookAhead = 25 // ms
let timer: ReturnType<typeof setInterval> | null = null
let onBeatCb: BeatCallback | null = null

// Tap-tempo state.
let taps: number[] = []

function scheduleBeat(beat: number, time: number) {
  const ac = getAudioContext()
  const isAccent = !!accents[beat]
  const freq = isAccent ? 1000 : 800
  const dur = isAccent ? 0.1 : 0.06

  const gain = ac.createGain()
  gain.gain.setValueAtTime(isAccent ? 0.5 : 0.35, time)
  gain.gain.exponentialRampToValueAtTime(0.001, time + dur)
  gain.connect(ac.destination)

  const osc = ac.createOscillator()
  osc.type = 'square'
  osc.frequency.value = freq
  osc.connect(gain)
  osc.start(time)
  osc.stop(time + dur)

  if (onBeatCb) {
    let delay = (time - ac.currentTime) * 1000
    if (delay < 0) delay = 0
    const cb = onBeatCb
    setTimeout(() => cb(beat, isAccent), delay)
  }
}

function scheduler() {
  const ac = getAudioContext()
  while (nextNoteTime < ac.currentTime + scheduleAhead) {
    scheduleBeat(currentBeat, nextNoteTime)
    const secondsPerBeat = 60.0 / bpm.value
    nextNoteTime += secondsPerBeat
    currentBeat = (currentBeat + 1) % beats.value
  }
}

export function useMetronome() {
  function start() {
    if (running.value) return
    running.value = true
    currentBeat = 0
    nextNoteTime = getAudioContext().currentTime
    timer = setInterval(scheduler, lookAhead)
  }

  function stop() {
    running.value = false
    if (timer !== null) {
      clearInterval(timer)
      timer = null
    }
  }

  function setBpm(n: number) {
    bpm.value = Math.max(30, Math.min(300, n))
    return bpm.value
  }

  function setBeats(n: number) {
    beats.value = n
    currentBeat = 0
    // Drop accents outside the new beat range.
    for (const k of Object.keys(accents)) {
      if (Number(k) >= n) delete accents[Number(k)]
    }
  }

  function toggleAccent(beat: number) {
    if (accents[beat]) delete accents[beat]
    else accents[beat] = true
  }

  function isAccent(beat: number): boolean {
    return !!accents[beat]
  }

  function onBeat(fn: BeatCallback | null) {
    onBeatCb = fn
  }

  function tapTempo(): number | null {
    const now = Date.now()
    taps.push(now)
    taps = taps.filter((t) => now - t < 2000)
    if (taps.length < 2) return null
    let sum = 0
    for (let i = 1; i < taps.length; i++) sum += taps[i] - taps[i - 1]
    const avg = sum / (taps.length - 1)
    return Math.round(60000 / avg)
  }

  return {
    running,
    bpm,
    beats,
    accents,
    start,
    stop,
    setBpm,
    setBeats,
    toggleAccent,
    isAccent,
    onBeat,
    tapTempo,
  }
}

// Convenience: page can call this from setup() so that leaving the route stops
// the metronome and detaches the beat callback automatically.
export function useMetronomeAutoStop() {
  const m = useMetronome()
  onBeforeUnmount(() => {
    m.stop()
    m.onBeat(null)
  })
  return m
}
