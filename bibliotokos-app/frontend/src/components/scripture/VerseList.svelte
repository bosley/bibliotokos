<script>
  import { afterUpdate, beforeUpdate, createEventDispatcher, onDestroy } from 'svelte'
  import VerseItem from './VerseItem.svelte'

  export let chunks = []
  export let loading = false
  export let query = ''
  export let generation = 0

  const dispatch = createEventDispatcher()
  const FALLBACK_VERSE_HEIGHT = 32

  let listEl
  let chunkEls = []
  let observer
  let visible = new Set()
  let lastGeneration = -1
  let pendingScrollReset = false
  let anchorIndex = -1
  let anchorOffset = 0
  let suppress = 0
  let suppressAt = 0

  $: if (generation !== lastGeneration) {
    lastGeneration = generation
    pendingScrollReset = true
    for (const i of [...visible]) {
      if (i >= chunks.length) visible.delete(i)
    }
  }

  $: avgVerseHeight = computeAvg(chunks)

  function computeAvg(cs) {
    let h = 0
    let n = 0
    for (const c of cs) {
      if (c.verses && c.height != null) {
        h += c.height
        n += c.count
      }
    }
    return n > 0 ? h / n : FALLBACK_VERSE_HEIGHT
  }

  function ensureObserver() {
    if (observer || !listEl) return
    observer = new IntersectionObserver(
      entries => {
        for (const entry of entries) {
          const idx = Number(entry.target.dataset.chunkIndex)
          if (Number.isNaN(idx)) continue
          if (entry.isIntersecting) visible.add(idx)
          else visible.delete(idx)
        }
        emitVisible()
      },
      { root: listEl, rootMargin: '50% 0px' }
    )
  }

  function emitVisible() {
    dispatch('visiblechunks', new Set([...visible].filter(i => i < chunks.length)))
  }

  function observe(el) {
    ensureObserver()
    observer.observe(el)
    return {
      destroy() {
        observer?.unobserve(el)
      },
    }
  }

  function elTop(el) {
    return (
      el.getBoundingClientRect().top -
      listEl.getBoundingClientRect().top +
      listEl.scrollTop
    )
  }

  function setScrollTop(target) {
    const max = Math.max(0, listEl.scrollHeight - listEl.clientHeight)
    const t = Math.max(0, Math.min(target, max))
    if (Math.abs(listEl.scrollTop - t) <= 1) return
    suppress++
    suppressAt = performance.now()
    listEl.scrollTop = t
  }

  function handleScroll() {
    if (suppress > 0) {
      if (performance.now() - suppressAt < 300) {
        suppress--
        return
      }
      suppress = 0
    }
    dispatch('scrollsync', getAnchor())
  }

  export function getAnchor() {
    if (!listEl || chunks.length === 0) return null
    const st = listEl.scrollTop
    let base = 0
    for (let i = 0; i < chunks.length; i++) {
      const count = chunks[i].count
      const el = chunkEls[i]
      if (!el || !el.isConnected) {
        base += count
        continue
      }
      const top = elTop(el)
      const h = el.offsetHeight
      if (top + h <= st) {
        base += count
        continue
      }
      if (chunks[i].verses) {
        const kids = el.children
        for (let j = 0; j < kids.length; j++) {
          const kt = elTop(kids[j])
          const kh = kids[j].offsetHeight
          if (kt + kh > st) {
            return {
              index: base + j,
              fraction: kh > 0 ? Math.max(0, (st - kt) / kh) : 0,
            }
          }
        }
        return { index: base + count - 1, fraction: 0 }
      }
      const within = h > 0 ? ((st - top) / h) * count : 0
      const j = Math.max(0, Math.min(count - 1, Math.floor(within)))
      return { index: base + j, fraction: Math.max(0, within - j) }
    }
    return null
  }

  export function scrollToAnchor(anchor) {
    if (!listEl || !anchor || chunks.length === 0) return
    let base = 0
    for (let i = 0; i < chunks.length; i++) {
      const count = chunks[i].count
      if (anchor.index >= base + count) {
        base += count
        continue
      }
      const el = chunkEls[i]
      if (!el || !el.isConnected) return
      const j = anchor.index - base
      if (chunks[i].verses && el.children.length > j) {
        const kid = el.children[j]
        setScrollTop(elTop(kid) + anchor.fraction * kid.offsetHeight)
      } else {
        const per = count > 0 ? el.offsetHeight / count : 0
        setScrollTop(elTop(el) + (j + anchor.fraction) * per)
      }
      return
    }
  }

  beforeUpdate(() => {
    anchorIndex = -1
    if (!listEl || pendingScrollReset) return
    const st = listEl.scrollTop
    for (let i = 0; i < chunkEls.length; i++) {
      const el = chunkEls[i]
      if (!el || !el.isConnected) continue
      if (elTop(el) + el.offsetHeight > st) {
        anchorIndex = i
        anchorOffset = st - elTop(el)
        break
      }
    }
  })

  afterUpdate(() => {
    if (!listEl) return
    for (let i = 0; i < chunks.length; i++) {
      const el = chunkEls[i]
      if (el && el.isConnected && chunks[i].verses) {
        chunks[i].height = el.offsetHeight
      }
    }
    if (pendingScrollReset) {
      pendingScrollReset = false
      setScrollTop(0)
      return
    }
    if (anchorIndex < 0) return
    const el = chunkEls[anchorIndex]
    if (!el || !el.isConnected) return
    const maxOffset = Math.max(0, el.offsetHeight - 1)
    setScrollTop(elTop(el) + Math.min(anchorOffset, maxOffset))
  })

  onDestroy(() => {
    observer?.disconnect()
  })
</script>

<div class="verse-list" bind:this={listEl} on:scroll={handleScroll}>
  {#if loading}
    <div class="state-msg">Loading…</div>
  {:else if !query}
    <div class="state-msg empty">
      <p>Enter a scripture reference to begin.</p>
      <p class="hint">Try <em>John 3:16</em>, <em>Gen 1</em>, <em>Romans 8:28-39</em></p>
    </div>
  {:else if chunks.length === 0}
    <div class="state-msg">No results for <em>{query}</em></div>
  {:else}
    {#each chunks as chunk, i}
      <div
        class="chunk"
        data-chunk-index={i}
        use:observe
        bind:this={chunkEls[i]}
      >
        {#if chunk.verses}
          {#each chunk.verses as verse (`${verse.book}-${verse.chapter}-${verse.verse}`)}
            <VerseItem {verse} />
          {/each}
        {:else}
          <div
            class="placeholder"
            style="height: {chunk.height ?? Math.max(1, chunk.count * avgVerseHeight)}px"
          ></div>
        {/if}
      </div>
    {/each}
  {/if}
</div>

<style>
  .verse-list {
    flex: 1;
    overflow-y: auto;
    padding: 20px 24px;
  }

  .state-msg {
    color: var(--text-muted);
    font-size: 14px;
    text-align: center;
    padding: 60px 24px;
    line-height: 1.8;
  }

  .state-msg .hint {
    margin-top: 6px;
    font-size: 13px;
  }

  .state-msg em {
    color: var(--accent);
    font-style: normal;
  }
</style>
