import { ref } from 'vue'

const KEY = 'chord-tools-theme'

type Theme = 'light' | 'dark'

function readInitial(): Theme {
  const stored = localStorage.getItem(KEY) as Theme | null
  if (stored === 'light' || stored === 'dark') return stored
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

const theme = ref<Theme>(readInitial())

function apply(next: Theme) {
  theme.value = next
  document.documentElement.setAttribute('data-theme', next)
  localStorage.setItem(KEY, next)
}

// Apply on first import to keep <html data-theme> in sync with the ref.
apply(theme.value)

export function useTheme() {
  return {
    theme,
    toggle() {
      apply(theme.value === 'light' ? 'dark' : 'light')
    },
  }
}
