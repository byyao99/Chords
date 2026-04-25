import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import TransposePage from '../pages/TransposePage.vue'
import ChordsPage from '../pages/ChordsPage.vue'
import DiatonicPage from '../pages/DiatonicPage.vue'
import ScalesPage from '../pages/ScalesPage.vue'
import MetronomePage from '../pages/MetronomePage.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'transpose', component: TransposePage, meta: { title: 'Transpose' } },
  { path: '/chords', name: 'chords', component: ChordsPage, meta: { title: 'Chord Lookup' } },
  { path: '/diatonic', name: 'diatonic', component: DiatonicPage, meta: { title: 'Diatonic Chords' } },
  { path: '/scales', name: 'scales', component: ScalesPage, meta: { title: 'Scales' } },
  { path: '/metronome', name: 'metronome', component: MetronomePage, meta: { title: 'Metronome' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.afterEach((to) => {
  const title = (to.meta?.title as string | undefined) ?? 'Chord Tools'
  document.title = `${title} - Chord Tools`
})

export default router
