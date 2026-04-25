<script setup lang="ts">
import type { DiatonicChord } from '../types/api'
import { playChordVoicing } from '../composables/useAudio'
import ChordDiagram from './ChordDiagram.vue'

const props = defineProps<{ chord: DiatonicChord }>()

function play() {
  if (props.chord.voicing) playChordVoicing(props.chord.voicing.frets)
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Enter' || e.key === ' ') {
    e.preventDefault()
    play()
  }
}
</script>

<template>
  <div
    v-if="chord.voicing"
    class="diatonic-card playable-chord"
    role="button"
    tabindex="0"
    :title="`Click to play`"
    :aria-label="`Play ${chord.name}`"
    @click="play"
    @keydown="onKey"
  >
    <div class="diatonic-degree">{{ chord.degree }}</div>
    <div class="diatonic-name">{{ chord.name }}</div>
    <div class="diatonic-diagram">
      <ChordDiagram :voicing="chord.voicing" :name="chord.name" />
    </div>
    <div class="diatonic-quality">{{ chord.quality }}</div>
  </div>
  <div v-else class="diatonic-card">
    <div class="diatonic-degree">{{ chord.degree }}</div>
    <div class="diatonic-name">{{ chord.name }}</div>
    <div class="diatonic-quality">{{ chord.quality }}</div>
  </div>
</template>
