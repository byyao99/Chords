<script setup lang="ts">
import { computed } from 'vue'
import { useTheme } from '../composables/useTheme'

const { theme, toggle } = useTheme()
const buttonText = computed(() => (theme.value === 'dark' ? '☀' : '☾'))
const ariaPressed = computed(() => (theme.value === 'dark' ? 'true' : 'false'))
</script>

<template>
  <header class="site-header">
    <div class="site-title">Chord Tools</div>
    <nav class="tab-nav">
      <RouterLink to="/" exact-active-class="active">Transpose</RouterLink>
      <RouterLink to="/chords" exact-active-class="active">Chords</RouterLink>
      <RouterLink to="/diatonic" exact-active-class="active">Diatonic</RouterLink>
      <RouterLink to="/scales" exact-active-class="active">Scales</RouterLink>
      <RouterLink to="/metronome" exact-active-class="active">Metronome</RouterLink>
    </nav>
    <button
      id="theme-toggle"
      class="theme-toggle"
      type="button"
      aria-label="Toggle color theme"
      :aria-pressed="ariaPressed"
      @click="toggle"
    >{{ buttonText }}</button>
  </header>
</template>
