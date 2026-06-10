<script>
  import { onMount } from 'svelte'
  import { get } from 'svelte/store'
  import { query, versions } from '../../stores/scripture.js'
  import { Query } from '../../../bindings/bibliotokos/services/bible/bibleservice.js'
  import { GetLinkedNotes } from '../../../bindings/bibliotokos/services/notes/notesservice.js'
  import ScripturePane from './ScripturePane.svelte'
  import PaneDivider from './PaneDivider.svelte'
  import LinkedNotesBar from './LinkedNotesBar.svelte'

  let panes = [
    { id: crypto.randomUUID(), versionId: '', verses: [], widthFlex: 1, loading: false }
  ]

  let containerEl
  let linkedNotes = []

  async function refreshLinkedNotes(q) {
    if (!q) {
      linkedNotes = []
      return
    }
    linkedNotes = await GetLinkedNotes(q).catch(() => [])
  }

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
      if (changed && get(query)) queryAllPanes()
    })

    const unsubQuery = query.subscribe(q => {
      if (q) {
        queryAllPanes()
      } else {
        panes = panes.map(p => ({ ...p, verses: [], loading: false }))
      }
      refreshLinkedNotes(q)
    })

    return () => {
      unsubVersions()
      unsubQuery()
    }
  })

  function handleWindowFocus() {
    refreshLinkedNotes(get(query))
  }

  async function queryPane(index) {
    const pane = panes[index]
    if (!pane || !pane.versionId || !$query) return
    panes[index] = { ...pane, loading: true }
    panes = [...panes]
    try {
      const verses = await Query($query, [], [pane.versionId])
      panes[index] = { ...panes[index], verses: verses ?? [], loading: false }
    } catch {
      panes[index] = { ...panes[index], verses: [], loading: false }
    }
    panes = [...panes]
  }

  function queryAllPanes() {
    panes.forEach((_, i) => queryPane(i))
  }

  function addPane(afterIndex) {
    const defaultVersionId = $versions[0]?.id ?? ''
    const newPane = {
      id: crypto.randomUUID(),
      versionId: defaultVersionId,
      verses: [],
      widthFlex: 1,
      loading: false,
    }
    panes = [
      ...panes.slice(0, afterIndex + 1),
      newPane,
      ...panes.slice(afterIndex + 1),
    ]
    if ($query) {
      const newIndex = afterIndex + 1
      queryPane(newIndex)
    }
  }

  function removePane(index) {
    if (panes.length <= 1) return
    panes = panes.filter((_, i) => i !== index)
  }

  function changeVersion(index, versionId) {
    panes[index] = { ...panes[index], versionId }
    panes = [...panes]
    if ($query) queryPane(index)
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
