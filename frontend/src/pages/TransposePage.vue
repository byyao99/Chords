<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { postTranspose } from '../composables/useApi'

const SHEET_KEY = 'transpose-sheet'

const SEMITONE_CHOICES = Array.from({ length: 23 }, (_, i) => i - 11)
const KEY_OPTIONS: Array<{ value: string; label: string }> = [
  { value: 'C', label: 'C' }, { value: 'C#', label: 'C# / Db' }, { value: 'D', label: 'D' },
  { value: 'Eb', label: 'Eb / D#' }, { value: 'E', label: 'E' }, { value: 'F', label: 'F' },
  { value: 'F#', label: 'F# / Gb' }, { value: 'G', label: 'G' }, { value: 'Ab', label: 'Ab / G#' },
  { value: 'A', label: 'A' }, { value: 'Bb', label: 'Bb / A#' }, { value: 'B', label: 'B' },
  { value: 'Am', label: 'Am' }, { value: 'Bbm', label: 'Bbm' }, { value: 'Bm', label: 'Bm' },
  { value: 'Cm', label: 'Cm' }, { value: 'C#m', label: 'C#m' }, { value: 'Dm', label: 'Dm' },
  { value: 'Ebm', label: 'Ebm' }, { value: 'Em', label: 'Em' }, { value: 'Fm', label: 'Fm' },
  { value: 'F#m', label: 'F#m' }, { value: 'Gm', label: 'Gm' }, { value: 'G#m', label: 'G#m' },
]

const sheet = ref('')
const mode = ref<'semitone' | 'key'>('semitone')
const semitones = ref(0)
const fromKey = ref('C')
const toKey = ref('G')
const result = ref('')
const error = ref('')
const busy = ref(false)

onMounted(() => {
  const saved = localStorage.getItem(SHEET_KEY)
  if (saved) sheet.value = saved
})

watch(sheet, (v) => {
  localStorage.setItem(SHEET_KEY, v)
})

async function transpose() {
  if (!sheet.value.trim()) {
    error.value = 'Please enter a chord sheet.'
    return
  }
  error.value = ''
  busy.value = true
  try {
    const data = await postTranspose(
      mode.value === 'semitone'
        ? { sheet: sheet.value, mode: 'semitone', semitones: semitones.value }
        : { sheet: sheet.value, mode: 'key', from_key: fromKey.value, to_key: toKey.value },
    )
    if (data.error) {
      error.value = data.error
      result.value = ''
    } else {
      result.value = data.result
    }
  } catch (e) {
    error.value = `Request failed: ${(e as Error).message}`
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="container">
    <h1>Chord Transposer</h1>

    <div class="input-section">
      <label for="sheet">Paste your chord sheet</label>
      <textarea
        id="sheet"
        v-model="sheet"
        rows="12"
        placeholder="  Am        C          G          F&#10;Somebody once told me the world is gonna roll me&#10;&#10;  C         G          Am         F&#10;I ain't the sharpest tool in the shed"
      ></textarea>
    </div>

    <div class="controls">
      <div class="mode-select">
        <label><input type="radio" v-model="mode" value="semitone" /> By semitones</label>
        <label><input type="radio" v-model="mode" value="key" /> By target key</label>
      </div>

      <div v-show="mode === 'semitone'" class="mode-panel">
        <label for="semitones">Semitones</label>
        <select id="semitones" v-model.number="semitones">
          <option v-for="n in SEMITONE_CHOICES" :key="n" :value="n">
            {{ n === 0 ? '0 (no change)' : (n > 0 ? `+${n}` : String(n)) }}
          </option>
        </select>
      </div>

      <div v-show="mode === 'key'" class="mode-panel">
        <label for="from-key">From</label>
        <select id="from-key" v-model="fromKey">
          <option v-for="o in KEY_OPTIONS" :key="`f-${o.value}`" :value="o.value">{{ o.label }}</option>
        </select>

        <label for="to-key">To</label>
        <select id="to-key" v-model="toKey">
          <option v-for="o in KEY_OPTIONS" :key="`t-${o.value}`" :value="o.value">{{ o.label }}</option>
        </select>
      </div>

      <button id="transpose-btn" type="button" :disabled="busy" @click="transpose">
        {{ busy ? 'Transposing...' : 'Transpose' }}
      </button>
    </div>

    <div class="output-section">
      <label>Result</label>
      <pre id="output">{{ result }}</pre>
    </div>

    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>
