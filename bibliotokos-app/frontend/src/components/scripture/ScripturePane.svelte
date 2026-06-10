<script>
  import { createEventDispatcher, onDestroy } from 'svelte'
  import { query as queryStore } from '../../stores/scripture.js'
  import { QueryPage } from '../../../bindings/bibliotokos/services/bible/bibleservice.js'
  import PaneHeader from './PaneHeader.svelte'
  import VerseList from './VerseList.svelte'

  export let pane
  export let canRemove = false

  const dispatch = createEventDispatcher()

  const CHUNK_SIZE = 50
  const LOOKAHEAD = 1
  const KEEP_RADIUS = 3
  const SETTLE_MS = 120

  let chunks = []
  let loading = false
  let token = 0
  let visibleSet = new Set()
  let lastQ = null
  let lastVersion = null
  let settleTimer

  onDestroy(() => clearTimeout(settleTimer))

  $: maybeLoad($queryStore, pane.versionId)

  function maybeLoad(q, versionId) {
    if (q === lastQ && versionId === lastVersion) return
    lastQ = q
    lastVersion = versionId
    load(q, versionId)
  }

  async function load(q, versionId) {
    token++
    const t = token
    visibleSet = new Set()
    clearTimeout(settleTimer)
    if (!q || !versionId) {
      chunks = []
      loading = false
      dispatch('visiblerange', null)
      return
    }
    loading = true
    dispatch('visiblerange', null)
    try {
      const page = await QueryPage(q, versionId, 0, CHUNK_SIZE)
      if (t !== token) return
      const total = page.total ?? 0
      const count = Math.ceil(total / CHUNK_SIZE)
      chunks = Array.from({ length: count }, (_, i) => ({
        count: i === count - 1 ? total - i * CHUNK_SIZE : CHUNK_SIZE,
        verses: i === 0 ? page.verses ?? [] : null,
        height: null,
        pending: false,
      }))
      loading = false
      ensureWindow()
    } catch {
      if (t !== token) return
      chunks = []
      loading = false
      dispatch('visiblerange', null)
    }
  }

  function handleVisibleChunks(e) {
    visibleSet = e.detail
    clearTimeout(settleTimer)
    settleTimer = setTimeout(ensureWindow, SETTLE_MS)
  }

  function ensureWindow() {
    if (chunks.length === 0) {
      emitVisibleRange()
      return
    }
    const vis = [...visibleSet].filter(i => i >= 0 && i < chunks.length)
    const min = vis.length ? Math.min(...vis) : 0
    const max = vis.length ? Math.max(...vis) : 0
    const loadMin = Math.max(0, min - LOOKAHEAD)
    const loadMax = Math.min(chunks.length - 1, max + LOOKAHEAD)
    for (let i = loadMin; i <= loadMax; i++) {
      if (!chunks[i].verses && !chunks[i].pending) loadChunk(i)
    }
    let changed = false
    for (let i = 0; i < chunks.length; i++) {
      if (chunks[i].verses && (i < min - KEEP_RADIUS || i > max + KEEP_RADIUS)) {
        chunks[i] = { ...chunks[i], verses: null }
        changed = true
      }
    }
    if (changed) chunks = chunks
    emitVisibleRange()
  }

  async function loadChunk(i) {
    const t = token
    chunks[i] = { ...chunks[i], pending: true }
    chunks = chunks
    try {
      const page = await QueryPage(lastQ, lastVersion, i * CHUNK_SIZE, CHUNK_SIZE)
      if (t !== token) return
      chunks[i] = { ...chunks[i], verses: page.verses ?? [], pending: false }
      chunks = chunks
      emitVisibleRange()
    } catch {
      if (t !== token) return
      chunks[i] = { ...chunks[i], pending: false }
      chunks = chunks
    }
  }

  function emitVisibleRange() {
    const vis = [...visibleSet].filter(
      i => i >= 0 && i < chunks.length && chunks[i].verses?.length
    )
    if (vis.length === 0) {
      dispatch('visiblerange', null)
      return
    }
    const min = Math.min(...vis)
    const max = Math.max(...vis)
    const first = chunks[min].verses[0]
    const lastVerses = chunks[max].verses
    dispatch('visiblerange', { first, last: lastVerses[lastVerses.length - 1] })
  }
</script>

<div class="pane" style="flex: {pane.widthFlex}">
  <PaneHeader
    versionId={pane.versionId}
    {canRemove}
    on:versionchange={e => dispatch('versionchange', e.detail)}
    on:addpane={() => dispatch('addpane')}
    on:removepane={() => dispatch('removepane')}
  />
  <VerseList
    {chunks}
    {loading}
    query={$queryStore}
    generation={token}
    on:visiblechunks={handleVisibleChunks}
  />
</div>

<style>
  .pane {
    min-width: 200px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background: var(--pane-bg);
  }
</style>
