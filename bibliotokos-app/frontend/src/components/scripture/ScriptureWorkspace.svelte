<script>
  import { onMount, onDestroy } from 'svelte'
  import { get } from 'svelte/store'
  import { books, versions } from '../../stores/scripture.js'
  import { GetLinkedNotes } from '../../../bindings/bibliotokos/services/notes/notesservice.js'
  import ScripturePane from './ScripturePane.svelte'
  import PaneDivider from './PaneDivider.svelte'
  import LinkedNotesBar from './LinkedNotesBar.svelte'

  let panes = [
    { id: crypto.randomUUID(), versionId: '', widthFlex: 1 }
  ]

  let containerEl
  let linkedNotes = []
  let visibleRanges = new Map()
  let debounceTimer

  onMount(() => {
    const unsubVersions = versions.subscribe(vs => {
      if (vs.length === 0) return
      let changed = false
      panes = panes.map(p => {
        if (!p.versionId) {
          changed = true
          return { ...p, versionId: vs[0].id }
        }
        return p
      })
      if (changed) panes = [...panes]
    })

    return () => {
      unsubVersions()
    }
  })

  onDestroy(() => {
    clearTimeout(debounceTimer)
  })

  function bookOrd(abbr) {
    const b = get(books).find(bk => bk.abbr === abbr)
    return b ? b.ord : 0
  }

  function versePos(v) {
    return bookOrd(v.book) * 1_000_000 + v.chapter * 1_000 + v.verse
  }

  function handleVisibleRange(paneId, detail) {
    if (detail) {
      visibleRanges.set(paneId, detail)
    } else {
      visibleRanges.delete(paneId)
    }
    scheduleLinkedNotesRefresh()
  }

  function scheduleLinkedNotesRefresh() {
    clearTimeout(debounceTimer)
    debounceTimer = setTimeout(refreshLinkedNotes, 300)
  }

  async function refreshLinkedNotes() {
    const ranges = [...visibleRanges.values()]
    if (ranges.length === 0) {
      linkedNotes = []
      return
    }
    let first = ranges[0].first
    let last = ranges[0].last
    for (const r of ranges.slice(1)) {
      if (versePos(r.first) < versePos(first)) first = r.first
      if (versePos(r.last) > versePos(last)) last = r.last
    }
    const ref = `${first.book} ${first.chapter}:${first.verse} - ${last.book} ${last.chapter}:${last.verse}`
    linkedNotes = await GetLinkedNotes(ref).catch(() => [])
  }

  function handleWindowFocus() {
    refreshLinkedNotes()
  }

  function addPane(afterIndex) {
    const defaultVersionId = $versions[0]?.id ?? ''
    const newPane = {
      id: crypto.randomUUID(),
      versionId: defaultVersionId,
      widthFlex: 1,
    }
    panes = [
      ...panes.slice(0, afterIndex + 1),
      newPane,
      ...panes.slice(afterIndex + 1),
    ]
  }

  function removePane(index) {
    if (panes.length <= 1) return
    const removed = panes[index]
    panes = panes.filter((_, i) => i !== index)
    visibleRanges.delete(removed.id)
    scheduleLinkedNotesRefresh()
  }

  function changeVersion(index, versionId) {
    panes[index] = { ...panes[index], versionId }
    panes = [...panes]
  }

  function startResize(leftIndex, startX) {
    const rightIndex = leftIndex + 1
    const startLeftFlex = panes[leftIndex].widthFlex
    const startRightFlex = panes[rightIndex].widthFlex
    const totalFlex = panes.reduce((sum, p) => sum + p.widthFlex, 0)
    const cw = containerEl?.clientWidth ?? 800

    function onMove(e) {
      const dx = e.clientX - startX
      const deltaFlex = (dx / cw) * totalFlex
      panes[leftIndex] = { ...panes[leftIndex], widthFlex: Math.max(0.15, startLeftFlex + deltaFlex) }
      panes[rightIndex] = { ...panes[rightIndex], widthFlex: Math.max(0.15, startRightFlex - deltaFlex) }
      panes = [...panes]
    }

    function onUp() {
      document.removeEventListener('mousemove', onMove)
      document.removeEventListener('mouseup', onUp)
    }

    document.addEventListener('mousemove', onMove)
    document.addEventListener('mouseup', onUp)
  }
</script>

<svelte:window on:focus={handleWindowFocus} />

<div class="workspace">
  <div class="panes" bind:this={containerEl}>
    {#each panes as pane, i (pane.id)}
      <ScripturePane
        {pane}
        canRemove={panes.length > 1}
        on:addpane={() => addPane(i)}
        on:removepane={() => removePane(i)}
        on:versionchange={e => changeVersion(i, e.detail)}
        on:visiblerange={e => handleVisibleRange(pane.id, e.detail)}
      />
      {#if i < panes.length - 1}
        <PaneDivider on:dragstart={e => startResize(i, e.detail)} />
      {/if}
    {/each}
  </div>
  <LinkedNotesBar {linkedNotes} />
</div>

<style>
  .workspace {
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 100%;
    overflow: hidden;
  }

  .panes {
    flex: 1;
    display: flex;
    flex-direction: row;
    min-height: 0;
    overflow: hidden;
  }
</style>
