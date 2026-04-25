<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { ROOT_OPTIONS } from '../constants/notes'
import { octaveForScaleNote, playNoteByName } from '../composables/useAudio'
import { getScales } from '../composables/useApi'
import type { ScaleCard } from '../types/api'

const SCALE_TYPES: string[] = [
  'Major',
  'Natural Minor (Aeolian)',
  'Major Pentatonic',
  'Minor Pentatonic',
  'Major Blues',
  'Minor Blues',
  'Dorian',
  'Phrygian',
  'Lydian',
  'Mixolydian',
  'Locrian',
  'Harmonic Minor',
  'Melodic Minor',
  'Whole Tone',
  'Diminished (Half-Whole)',
  'Diminished (Whole-Half)',
]
const SCALE_TYPE_LABELS: Record<string, string> = {
  'Natural Minor (Aeolian)': 'Natural Minor',
  'Diminished (Half-Whole)': 'Diminished (HW)',
  'Diminished (Whole-Half)': 'Diminished (WH)',
}

const root = ref('C')
const selected = reactive<Record<string, boolean>>(
  Object.fromEntries(SCALE_TYPES.map((t) => [t, false])),
)
const dropdownOpen = ref(false)
const scales = ref<ScaleCard[]>([])
const loading = ref(false)
const error = ref('')

const selectedTypes = computed(() => SCALE_TYPES.filter((t) => selected[t]))
const allSelected = computed(() => selectedTypes.value.length === SCALE_TYPES.length)
const filterLabel = computed(() => {
  const n = selectedTypes.value.length
  if (n === SCALE_TYPES.length) return 'All'
  if (n === 0) return 'None'
  return `${n} selected`
})

function toggleAll(e: Event) {
  const checked = (e.target as HTMLInputElement).checked
  for (const t of SCALE_TYPES) selected[t] = checked
}

async function load() {
  if (selectedTypes.value.length === 0) {
    scales.value = []
    error.value = ''
    return
  }
  loading.value = true
  error.value = ''
  try {
    const types = selectedTypes.value.length === SCALE_TYPES.length ? [] : selectedTypes.value
    const data = await getScales(root.value, types)
    scales.value = data.scales ?? []
  } catch (e) {
    error.value = `Failed to load: ${(e as Error).message}`
    scales.value = []
  } finally {
    loading.value = false
  }
}

watch([root, selectedTypes], load, { deep: false })

function onDocClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.type-filter-wrap')) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onDocClick)
  load()
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocClick)
})

function playPill(note: string, octave: number) {
  playNoteByName(note, octave)
}
</script>

<template>
  <div class="container">
    <h1>Scale Reference</h1>

    <div class="scales-controls">
      <label for="root-select">Root</label>
      <select id="root-select" v-model="root">
        <option v-for="r in ROOT_OPTIONS" :key="r.value" :value="r.value">{{ r.label }}</option>
      </select>

      <label>Type</label>
      <div class="type-filter-wrap">
        <button
          type="button"
          class="type-filter-btn"
          @click.stop="dropdownOpen = !dropdownOpen"
        >{{ filterLabel }}</button>
        <div class="type-filter-dropdown" :class="{ hidden: !dropdownOpen }">
          <label class="type-filter-option type-filter-all">
            <input type="checkbox" :checked="allSelected" @change="toggleAll" /> Select All
          </label>
          <hr class="type-filter-divider" />
          <label v-for="t in SCALE_TYPES" :key="t" class="type-filter-option">
            <input type="checkbox" :value="t" v-model="selected[t]" />
            {{ SCALE_TYPE_LABELS[t] ?? t }}
          </label>
        </div>
      </div>
    </div>

    <div class="scales-grid" :class="{ 'is-loading': loading }">
      <div v-if="loading" class="loading-indicator">Loading…</div>
      <template v-else-if="error">
        <div class="error">{{ error }}</div>
      </template>
      <template v-else-if="selectedTypes.length === 0">
        <div class="empty-hint">Select scale types to display.</div>
      </template>
      <template v-else>
        <div v-for="s in scales" :key="s.name" class="scale-card">
          <div class="scale-name">{{ s.name }}</div>
          <div class="scale-notes">
            <button
              v-for="(note, idx) in s.notes"
              :key="idx"
              type="button"
              class="note-pill playable"
              :class="{ root: idx === 0 }"
              :aria-label="`Play note ${note}`"
              @click="playPill(note, octaveForScaleNote(s.notes[0], note))"
            >{{ note }}</button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
