import type {
  ChordDiagramResponse,
  DiatonicResponse,
  ProgressionsResponse,
  ScalesResponse,
  TransposeRequest,
  TransposeResponse,
} from '../types/api'

async function getJson<T>(url: string): Promise<T> {
  const res = await fetch(url)
  const data = (await res.json()) as T & { error?: string }
  if (!res.ok) {
    throw new Error(data.error || `Request failed (${res.status})`)
  }
  return data
}

export async function postTranspose(req: TransposeRequest): Promise<TransposeResponse> {
  const res = await fetch('/api/transpose', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  return (await res.json()) as TransposeResponse
}

export function getScales(root: string, types: string[]): Promise<ScalesResponse> {
  let url = `/api/scales?root=${encodeURIComponent(root)}`
  if (types.length > 0) {
    url += `&type=${encodeURIComponent(types.join(','))}`
  }
  return getJson<ScalesResponse>(url)
}

export function getDiatonic(root: string, minor: boolean): Promise<DiatonicResponse> {
  return getJson<DiatonicResponse>(
    `/api/diatonic?root=${encodeURIComponent(root)}&minor=${minor}`,
  )
}

export function getChordDiagram(name: string, maxFret: number): Promise<ChordDiagramResponse> {
  return getJson<ChordDiagramResponse>(
    `/api/chord-diagram?name=${encodeURIComponent(name)}&max_fret=${maxFret}`,
  )
}

export function getProgressions(
  root: string,
  minor: boolean,
  mood: string,
  style: string,
): Promise<ProgressionsResponse> {
  const params = new URLSearchParams({
    root,
    minor: String(minor),
    mood,
    style,
  })
  return getJson<ProgressionsResponse>(`/api/progressions?${params.toString()}`)
}
