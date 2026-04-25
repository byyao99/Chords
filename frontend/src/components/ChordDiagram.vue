<script setup lang="ts">
import { computed } from 'vue'
import type { ChordVoicing } from '../types/api'

const props = withDefaults(
  defineProps<{
    voicing: ChordVoicing
    name?: string
    notes?: string[]
    width?: number
    height?: number
  }>(),
  {
    name: '',
    notes: () => [],
    width: 80,
    height: 100,
  },
)

const PAD_TOP = 28
const PAD_LEFT = 16
const PAD_RIGHT = 8
const PAD_BOTTOM = 8
const STRINGS = 6
const FRET_COUNT = 4

interface Layout {
  fretW: number
  fretH: number
  dotR: number
  nutY: number
  displayBase: number
  baseFretLabel: { x: number; y: number; text: string } | null
  fretLines: { x1: number; y1: number; x2: number; y2: number }[]
  stringLines: { x1: number; y1: number; x2: number; y2: number }[]
  nutLine: { x1: number; y1: number; x2: number; y2: number; strokeW: number }
  markers: Array<
    | { kind: 'mute'; x: number; y: number }
    | { kind: 'open'; x: number; y: number; r: number }
    | { kind: 'dot'; x: number; y: number; r: number }
  >
}

const layout = computed<Layout>(() => {
  const w = props.width
  const h = props.height
  const fretW = (w - PAD_LEFT - PAD_RIGHT) / (STRINGS - 1)
  const fretH = (h - PAD_TOP - PAD_BOTTOM) / FRET_COUNT
  const dotR = fretW * 0.3
  const nutY = PAD_TOP

  const frets = props.voicing.frets ?? [0, 0, 0, 0, 0, 0]

  let minFret = 99
  let maxFret = 0
  for (let i = 0; i < 6; i++) {
    if (frets[i] > 0) {
      if (frets[i] < minFret) minFret = frets[i]
      if (frets[i] > maxFret) maxFret = frets[i]
    }
  }
  let displayBase = 1
  if (minFret > 4) displayBase = minFret
  else if (maxFret > 4) displayBase = minFret

  const nutLine = {
    x1: PAD_LEFT,
    y1: nutY,
    x2: PAD_LEFT + fretW * (STRINGS - 1),
    y2: nutY,
    strokeW: displayBase === 1 ? 3 : 1,
  }
  const baseFretLabel =
    displayBase === 1
      ? null
      : { x: PAD_LEFT - 10, y: nutY + fretH * 0.65, text: `${displayBase}fr` }

  const fretLines = [] as Layout['fretLines']
  for (let f = 1; f <= FRET_COUNT; f++) {
    const y = nutY + f * fretH
    fretLines.push({ x1: PAD_LEFT, y1: y, x2: PAD_LEFT + fretW * (STRINGS - 1), y2: y })
  }

  const stringLines = [] as Layout['stringLines']
  for (let s = 0; s < STRINGS; s++) {
    const x = PAD_LEFT + s * fretW
    stringLines.push({ x1: x, y1: nutY, x2: x, y2: nutY + FRET_COUNT * fretH })
  }

  const markers: Layout['markers'] = []
  for (let s = 0; s < STRINGS; s++) {
    const x = PAD_LEFT + s * fretW
    const fretVal = frets[s]
    if (fretVal === -1) {
      markers.push({ kind: 'mute', x, y: nutY - 6 })
    } else if (fretVal === 0) {
      markers.push({ kind: 'open', x, y: nutY - 8, r: dotR * 0.8 })
    } else {
      const relFret = fretVal - displayBase + 1
      if (relFret >= 1 && relFret <= FRET_COUNT) {
        const cy = nutY + (relFret - 0.5) * fretH
        markers.push({ kind: 'dot', x, y: cy, r: dotR })
      }
    }
  }

  return { fretW, fretH, dotR, nutY, displayBase, baseFretLabel, fretLines, stringLines, nutLine, markers }
})

function buildAriaLabel(name: string, frets: number[], notes: string[]): string {
  const parts: string[] = []
  if (name) parts.push(`${name} chord`)
  if (frets && frets.length) {
    const fretStr = frets.map((f) => (f === -1 ? 'x' : String(f))).join(' ')
    parts.push(`Frets: ${fretStr}`)
  }
  if (notes && notes.length) parts.push(`Notes: ${notes.join(' ')}`)
  return parts.length ? parts.join('. ') + '.' : 'Chord diagram.'
}

const ariaLabel = computed(() =>
  buildAriaLabel(props.name || props.voicing.name, props.voicing.frets, props.notes),
)
</script>

<template>
  <svg
    class="chord-svg"
    role="img"
    :width="width"
    :height="height"
    :viewBox="`0 0 ${width} ${height}`"
    :aria-label="ariaLabel"
  >
    <title>{{ ariaLabel }}</title>

    <line
      :x1="layout.nutLine.x1" :y1="layout.nutLine.y1"
      :x2="layout.nutLine.x2" :y2="layout.nutLine.y2"
      stroke="currentColor" :stroke-width="layout.nutLine.strokeW"
    />
    <text
      v-if="layout.baseFretLabel"
      :x="layout.baseFretLabel.x" :y="layout.baseFretLabel.y"
      text-anchor="middle" font-size="8" font-family="-apple-system, sans-serif"
      fill="currentColor"
    >{{ layout.baseFretLabel.text }}</text>

    <line
      v-for="(l, i) in layout.fretLines" :key="`f${i}`"
      :x1="l.x1" :y1="l.y1" :x2="l.x2" :y2="l.y2"
      stroke="currentColor" stroke-width="1"
    />
    <line
      v-for="(l, i) in layout.stringLines" :key="`s${i}`"
      :x1="l.x1" :y1="l.y1" :x2="l.x2" :y2="l.y2"
      stroke="currentColor" stroke-width="1"
    />

    <template v-for="(m, i) in layout.markers" :key="`m${i}`">
      <text
        v-if="m.kind === 'mute'"
        :x="m.x" :y="m.y" text-anchor="middle" font-size="11"
        font-family="-apple-system, sans-serif" fill="currentColor"
      >&#215;</text>
      <circle
        v-else-if="m.kind === 'open'"
        :cx="m.x" :cy="m.y" :r="m.r"
        fill="none" stroke="currentColor" stroke-width="1.5"
      />
      <circle
        v-else
        :cx="m.x" :cy="m.y" :r="m.r" fill="currentColor"
      />
    </template>
  </svg>
</template>
