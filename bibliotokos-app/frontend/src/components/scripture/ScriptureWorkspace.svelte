<script>
  import { query, versions } from '../../stores/scripture.js'
  import { Query } from '../../../bindings/bibliotokos/services/bible/bibleservice.js'
  import ScripturePane from './ScripturePane.svelte'
  import PaneDivider from './PaneDivider.svelte'

  let panes = [
    { id: crypto.randomUUID(), versionId: '', verses: [], widthFlex: 1, loading: false }
  ]

  let containerEl

  $: if ($versions.length > 0 && panes.some(p => !p.versionId)) {
    panes = panes.map(p => ({ ...p, versionId: p.versionId || $versions[0].id }))
    if ($query) queryAllPanes()
  }

  $: if ($query) {
    queryAllPanes()
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

<div class="workspace" bind:this={containerEl}>
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

<style>
  .workspace {
    display: flex;
    flex-direction: row;
    height: 100%;
    width: 100%;
    overflow: hidden;
  }
</style>
