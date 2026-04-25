<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import { useMetronomeAutoStop } from '../composables/useMetronome'

const m = useMetronomeAutoStop()
const playingBeat = ref(-1)
const playingAccent = ref(false)

const beatDots = computed(() => Array.from({ length: m.beats.value }, (_, i) => i))

function toggleStart() {
  if (m.running.value) {
    m.stop()
    playingBeat.value = -1
  } else {
    m.start()
  }
}

function tap() {
  const newBpm = m.tapTempo()
  if (newBpm && newBpm >= 30 && newBpm <= 300) m.setBpm(newBpm)
}

function onSpace(e: KeyboardEvent) {
  if (e.code === 'Space' && e.target === document.body) {
    e.preventDefault()
    toggleStart()
  }
}

onMounted(() => {
  m.onBeat((beat, isAccent) => {
    playingBeat.value = beat
    playingAccent.value = isAccent
  })
  document.addEventListener('keydown', onSpace)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onSpace)
})
</script>

<template>
  <div class="container">
    <h1>Metronome</h1>

    <div class="metronome-wrapper">
      <div class="bpm-section">
        <div class="bpm-display">{{ m.bpm.value }}</div>
        <input
          type="range" class="bpm-slider"
          :min="30" :max="300"
          :value="m.bpm.value"
          @input="m.setBpm(parseInt(($event.target as HTMLInputElement).value, 10))"
        />
      </div>

      <div class="metronome-controls">
        <div>
          <label for="time-sig">Beats</label>
          <select
            id="time-sig"
            :value="m.beats.value"
            @change="m.setBeats(parseInt(($event.target as HTMLSelectElement).value, 10))"
          >
            <option v-for="n in [2, 3, 4, 5, 6, 7]" :key="n" :value="n">{{ n }}</option>
          </select>
        </div>
      </div>

      <div class="beat-display">
        <button
          v-for="i in beatDots"
          :key="i"
          type="button"
          class="beat-dot"
          :class="{
            'accent-mark': m.isAccent(i),
            active: m.running.value && playingBeat === i && !playingAccent,
            'playing-accent': m.running.value && playingBeat === i && playingAccent,
          }"
          :aria-label="`Beat ${i + 1} accent`"
          :aria-pressed="m.isAccent(i) ? 'true' : 'false'"
          @click="m.toggleAccent(i)"
        ></button>
      </div>

      <div class="metronome-buttons">
        <button
          class="start-btn"
          type="button"
          :class="{ running: m.running.value }"
          @click="toggleStart"
        >{{ m.running.value ? 'Stop' : 'Start' }}</button>
        <button class="tap-btn" type="button" @click="tap">Tap</button>
      </div>
    </div>
  </div>
</template>
