<script>
  import { createEventDispatcher } from 'svelte'
  import { query as queryStore } from '../../stores/scripture.js'
  import PaneHeader from './PaneHeader.svelte'
  import VerseList from './VerseList.svelte'

  export let pane
  export let canRemove = false

  const dispatch = createEventDispatcher()
</script>

<div class="pane" style="flex: {pane.widthFlex}">
  <PaneHeader
    versionId={pane.versionId}
    {canRemove}
    on:versionchange={e => dispatch('versionchange', e.detail)}
    on:addpane={() => dispatch('addpane')}
    on:removepane={() => dispatch('removepane')}
  />
  <VerseList verses={pane.verses} loading={pane.loading} query={$queryStore} />
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
