<script>
  import { createEventDispatcher } from 'svelte'
  import { versions } from '../../stores/scripture.js'

  export let versionId = ''
  export let canRemove = false

  const dispatch = createEventDispatcher()
</script>

<div class="pane-header">
  <select
    value={versionId}
    on:change={e => dispatch('versionchange', e.currentTarget.value)}
  >
    {#each $versions as v (v.id)}
      <option value={v.id}>{v.name} ({v.id.toUpperCase()})</option>
    {/each}
  </select>

  <div class="actions">
    <button on:click={() => dispatch('addpane')} title="Add translation pane" aria-label="Add pane">
      <svg viewBox="0 0 16 16" fill="currentColor" width="12" height="12">
        <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
      </svg>
    </button>
    {#if canRemove}
      <button on:click={() => dispatch('removepane')} title="Close pane" aria-label="Close pane">
        <svg viewBox="0 0 16 16" fill="currentColor" width="12" height="12">
          <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708z"/>
        </svg>
      </button>
    {/if}
  </div>
</div>

<style>
  .pane-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 16px;
    border-bottom: 1px solid var(--border);
    background: var(--bg-surface);
    flex-shrink: 0;
    gap: 8px;
  }

  select {
    flex: 1;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    padding: 4px 8px;
    font-size: 12px;
    color: var(--text);
    font-family: var(--font-sans);
    cursor: pointer;
    outline: none;
    -webkit-user-select: none;
    user-select: none;
  }

  select:focus {
    border-color: var(--accent);
  }

  .actions {
    display: flex;
    gap: 2px;
    flex-shrink: 0;
  }

  button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: var(--radius-sm);
    color: var(--text-muted);
    transition: color 0.15s, background 0.15s;
  }

  button:hover {
    color: var(--text);
    background: var(--bg-hover);
  }
</style>
