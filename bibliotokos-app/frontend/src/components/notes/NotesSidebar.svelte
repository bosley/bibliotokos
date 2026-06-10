<script>
  import { createEventDispatcher } from 'svelte'
  import NotesList from './NotesList.svelte'

  export let noteHeaders = []
  export let tags = []
  export let activeNoteId = null
  export let selectedTagId = null

  const dispatch = createEventDispatcher()

  $: filteredHeaders = selectedTagId === '__untagged__'
    ? noteHeaders.filter(h => h.tags.length === 0)
    : selectedTagId
      ? noteHeaders.filter(h => h.tags.includes(selectedTagId))
      : noteHeaders

  function selectTag(id) {
    dispatch('selecttag', selectedTagId === id ? null : id)
  }
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <span class="label">Notes</span>
    <div class="header-actions">
      <button on:click={() => dispatch('create')} title="New note" aria-label="New note" class="icon-btn">
        <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
        </svg>
      </button>
      <button on:click={() => dispatch('collapse')} title="Collapse sidebar" aria-label="Collapse sidebar" class="icon-btn">
        <svg viewBox="0 0 16 16" fill="currentColor" width="12" height="12">
          <path d="M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z"/>
        </svg>
      </button>
    </div>
  </div>

  {#if tags.length > 0}
    <div class="tag-filter">
      <button
        class="tag-pill"
        class:active={selectedTagId === null}
        on:click={() => dispatch('selecttag', null)}
      >All</button>
      {#each tags as tag (tag.id)}
        <span class="tag-pill-wrap" class:active={selectedTagId === tag.name}>
          <button
            class="tag-pill"
            on:click={() => selectTag(tag.name)}
          >{tag.name}</button>
          <button
            class="tag-delete"
            on:click|stopPropagation={() => dispatch('deletetag', tag.id)}
            title="Delete tag"
            aria-label="Delete tag {tag.name}"
          >×</button>
        </span>
      {/each}
      <button
        class="tag-pill"
        class:active={selectedTagId === '__untagged__'}
        on:click={() => selectTag('__untagged__')}
      >Untagged</button>
    </div>
  {/if}

  <NotesList noteHeaders={filteredHeaders} {activeNoteId} on:select on:delete />
</aside>

<style>
  .sidebar {
    width: var(--sidebar-w, 260px);
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    border-right: 1px solid var(--border);
    background: var(--bg-surface);
    overflow: hidden;
  }

  .sidebar-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border);
    flex-shrink: 0;
  }

  .label {
    font-size: 11px;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: var(--radius-sm);
    color: var(--text-muted);
    transition: color 0.15s, background 0.15s;
  }

  .icon-btn:hover {
    color: var(--text);
    background: var(--bg-hover);
  }

  .tag-filter {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    padding: 8px 12px;
    border-bottom: 1px solid var(--border);
    flex-shrink: 0;
  }

  .tag-pill-wrap {
    display: inline-flex;
    align-items: center;
    border-radius: 999px;
    background: var(--bg-hover);
    transition: background 0.1s;
    overflow: hidden;
  }

  .tag-pill-wrap.active {
    background: color-mix(in srgb, var(--accent) 15%, transparent);
  }

  .tag-pill-wrap.active .tag-pill {
    color: var(--accent);
  }

  .tag-pill {
    padding: 2px 6px 2px 8px;
    font-size: 11px;
    font-weight: 500;
    color: var(--text-muted);
    transition: color 0.1s;
    cursor: pointer;
    white-space: nowrap;
    background: transparent;
  }

  .tag-pill:hover {
    color: var(--text);
  }

  .tag-delete {
    display: none;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 16px;
    margin-right: 3px;
    border-radius: 50%;
    font-size: 13px;
    line-height: 1;
    color: var(--text-muted);
    transition: color 0.1s, background 0.1s;
    cursor: pointer;
    background: transparent;
  }

  .tag-pill-wrap:hover .tag-delete {
    display: flex;
  }

  .tag-delete:hover {
    color: var(--text);
    background: color-mix(in srgb, currentColor 15%, transparent);
  }
</style>
