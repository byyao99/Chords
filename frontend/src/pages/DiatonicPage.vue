<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { ROOT_OPTIONS } from '../constants/notes'
import { getDiatonic, getProgressions } from '../composables/useApi'
import type { DiatonicChord, Progression } from '../types/api'
import ChordCard from '../components/ChordCard.vue'

const MOOD_OPTIONS = [
  { value: '', label: 'Any' },
  { value: 'happy', label: 'Happy' },
  { value: 'sad', label: 'Sad' },
  { value: 'energetic', label: 'Energetic' },
  { value: 'dark', label: 'Dark' },
  { value: 'chill', label: 'Chill' },
]
const STYLE_OPTIONS = [
  { value: '', label: 'Any' },
  { value: 'pop', label: 'Pop' },
  { value: 'rock', label: 'Rock' },
  { value: 'jazz', label: 'Jazz' },
  { value: 'blues', label: 'Blues' },
  { value: 'folk', label: 'Folk' },
]

const root = ref('C')
const isMinor = ref(false)
const mood = ref('')
const style = ref('')

const keyLabel = ref('')
const chords = ref<DiatonicChord[]>([])
const diatonicLoading = ref(false)
const diatonicError = ref('')

const progressions = ref<Progression[]>([])
const progLoading = ref(false)
const progError = ref('')

async function loadDiatonic() {
  diatonicLoading.value = true
  diatonicError.value = ''
  try {
    const data = await getDiatonic(root.value, isMinor.value)
    keyLabel.value = data.key
    chords.value = data.chords ?? []
  } catch (e) {
    diatonicError.value = `Failed to load: ${(e as Error).message}`
    chords.value = []
  } finally {
    diatonicLoading.value = false
  }
}

async function loadProgressions() {
  progLoading.value = true
  progError.value = ''
  try {
    const data = await getProgressions(root.value, isMinor.value, mood.value, style.value)
    progressions.value = data.progressions ?? []
  } catch (e) {
    progError.value = `Failed to load: ${(e as Error).message}`
    progressions.value = []
  } finally {
    progLoading.value = false
  }
}

function loadAll() {
  loadDiatonic()
  loadProgressions()
}

watch([root, isMinor], loadAll)
watch([mood, style], loadProgressions)
onMounted(loadAll)
</script>

<template>
  <div class="container">
    <h1>Diatonic Chords</h1>

    <div class="scales-controls">
      <label for="root-select">Root</label>
      <select id="root-select" v-model="root">
        <option v-for="r in ROOT_OPTIONS" :key="r.value" :value="r.value">{{ r.label }}</option>
      </select>

      <div class="mode-toggle">
        <button
          type="button"
          class="toggle-btn"
          :class="{ active: !isMinor }"
          @click="isMinor = false"
        >Major</button>
        <button
          type="button"
          class="toggle-btn"
          :class="{ active: isMinor }"
          @click="isMinor = true"
        >Minor</button>
      </div>
    </div>

    <div class="key-label">{{ keyLabel }}</div>
    <div class="diatonic-grid" :class="{ 'is-loading': diatonicLoading }">
      <div v-if="diatonicLoading" class="loading-indicator">Loading…</div>
      <div v-else-if="diatonicError" class="error">{{ diatonicError }}</div>
      <template v-else>
        <ChordCard v-for="ch in chords" :key="ch.name + ch.degree" :chord="ch" />
      </template>
    </div>

    <h1>Chord Progressions</h1>
    <div class="scales-controls">
      <label for="mood-select">Mood</label>
      <select id="mood-select" v-model="mood">
        <option v-for="m in MOOD_OPTIONS" :key="m.value" :value="m.value">{{ m.label }}</option>
      </select>

      <label for="style-select">Style</label>
      <select id="style-select" v-model="style">
        <option v-for="s in STYLE_OPTIONS" :key="s.value" :value="s.value">{{ s.label }}</option>
      </select>
    </div>

    <div :class="{ 'is-loading': progLoading }">
      <div v-if="progLoading" class="loading-indicator">Loading…</div>
      <div v-else-if="progError" class="error">{{ progError }}</div>
      <div v-else-if="progressions.length === 0" class="error">No progressions found.</div>
      <template v-else>
        <div v-for="(p, pi) in progressions" :key="pi" class="progression-block">
          <div class="progression-header">
            <span class="progression-name">{{ p.pattern_name }}</span>
            <span class="progression-tags">{{ p.style }} · {{ p.mood }}</span>
          </div>
          <div class="diatonic-grid">
            <ChordCard v-for="(ch, ci) in p.chords" :key="ci" :chord="ch" />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
