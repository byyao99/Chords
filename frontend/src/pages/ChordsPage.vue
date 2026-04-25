<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { CHROMATIC_ROOTS } from '../constants/notes'
import { getChordDiagram } from '../composables/useApi'
import { playChordVoicing } from '../composables/useAudio'
import type { ChordVoicing } from '../types/api'
import ChordDiagram from '../components/ChordDiagram.vue'

interface QualityOption { value: string; label: string }
interface QualityGroup { label: string; options: QualityOption[] }

const QUALITY_GROUPS: QualityGroup[] = [
  { label: 'Triad', options: [
    { value: '', label: 'Major' },
    { value: 'm', label: 'm (minor)' },
    { value: 'aug', label: 'aug' },
    { value: 'dim', label: 'dim' },
    { value: '5', label: '5 (power)' },
  ]},
  { label: '7th', options: [
    { value: '7', label: '7' }, { value: 'm7', label: 'm7' }, { value: 'maj7', label: 'maj7' },
    { value: 'dim7', label: 'dim7' }, { value: 'm7b5', label: 'm7b5 (half-dim)' },
    { value: 'aug7', label: 'aug7' }, { value: 'mmaj7', label: 'mmaj7' },
  ]},
  { label: '6th', options: [{ value: '6', label: '6' }, { value: 'm6', label: 'm6' }] },
  { label: 'Sus', options: [
    { value: 'sus2', label: 'sus2' }, { value: 'sus4', label: 'sus4' }, { value: '7sus4', label: '7sus4' },
  ]},
  { label: 'Add / Extended', options: [
    { value: 'add9', label: 'add9' }, { value: '9', label: '9' },
    { value: 'm9', label: 'm9' }, { value: 'maj9', label: 'maj9' },
  ]},
  { label: 'Altered', options: [
    { value: '7b9', label: '7b9' }, { value: '7#9', label: '7#9' },
  ]},
]
const MAX_FRET_OPTIONS = [5, 7, 9, 12]

const root = ref('C')
const quality = ref('')
const maxFret = ref(5)
const loading = ref(false)
const error = ref('')
const voicings = ref<ChordVoicing[]>([])
const notes = ref<string[]>([])

const chordName = computed(() => root.value + quality.value)
const fretsLabel = (frets: number[]) => frets.map((f) => (f === -1 ? 'x' : String(f))).join('  ')

async function search() {
  if (!chordName.value) return
  loading.value = true
  error.value = ''
  try {
    const data = await getChordDiagram(chordName.value, maxFret.value)
    voicings.value = data.voicings ?? []
    notes.value = data.notes ?? []
  } catch (e) {
    error.value = (e as Error).message
    voicings.value = []
    notes.value = []
  } finally {
    loading.value = false
  }
}

watch([root, quality, maxFret], search)
onMounted(search)

function play(frets: number[]) {
  playChordVoicing(frets)
}

function onKey(e: KeyboardEvent, frets: number[]) {
  if (e.key === 'Enter' || e.key === ' ') {
    e.preventDefault()
    play(frets)
  }
}
</script>

<template>
  <div class="container">
    <h1>Chord Lookup</h1>

    <div class="scales-controls">
      <label for="root-select">Root</label>
      <select id="root-select" v-model="root">
        <option v-for="r in CHROMATIC_ROOTS" :key="r.value" :value="r.value">{{ r.label }}</option>
      </select>

      <label for="quality-select">Quality</label>
      <select id="quality-select" v-model="quality">
        <optgroup v-for="g in QUALITY_GROUPS" :key="g.label" :label="g.label">
          <option v-for="o in g.options" :key="`${g.label}-${o.value}`" :value="o.value">{{ o.label }}</option>
        </optgroup>
      </select>

      <label for="max-fret">Max Fret</label>
      <select id="max-fret" v-model.number="maxFret">
        <option v-for="n in MAX_FRET_OPTIONS" :key="n" :value="n">{{ n }}</option>
      </select>
    </div>

    <div class="chord-result" :class="{ 'is-loading': loading }">
      <div v-if="loading" class="loading-indicator">Loading…</div>
      <template v-else-if="!error && voicings.length">
        <div class="chord-display-name">{{ chordName }}</div>
        <div v-if="notes.length" class="chord-notes-label chord-notes-header">{{ notes.join('  ') }}</div>
        <div class="chord-cards-row">
          <div
            v-for="(v, i) in voicings"
            :key="i"
            class="chord-display-card playable-chord"
            role="button"
            tabindex="0"
            title="Click to play"
            :aria-label="`Play ${chordName}`"
            @click="play(v.frets)"
            @keydown="onKey($event, v.frets)"
          >
            <div class="chord-display-diagram">
              <ChordDiagram :voicing="v" :name="chordName" :notes="notes" :width="140" :height="170" />
            </div>
            <div class="chord-frets-label">{{ fretsLabel(v.frets) }}</div>
          </div>
        </div>
      </template>
    </div>

    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>
