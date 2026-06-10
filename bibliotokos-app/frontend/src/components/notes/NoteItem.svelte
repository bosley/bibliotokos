<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  export let note
  export let active = false

  function formatDate(iso) {
    try {
      return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
    } catch {
      return ''
    }
  }
</script>

<div
  class="note-item"
  class:active
  role="button"
  tabindex="0"
  on:click={() => dispatch('select', note.id)}
  on:keydown={e => e.key === 'Enter' && dispatch('select', note.id)}
>
  <div class="note-title">{note.title || 'Untitled Note'}</div>
  <div class="note-meta">
    <span class="date">{formatDate(note.updatedAt)}</span>
    <button
      class="delete"
      on:click|stopPropagation={() => dispatch('delete', note.id)}
      title="Delete note"
      aria-label="Delete note"
    >
      <svg viewBox="0 0 16 16" fill="currentColor" width="11" height="11">
        <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/>
        <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/>
      </svg>
    </button>
  </div>
</div>

<style>
  .note-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 10px 16px;
    cursor: pointer;
    border-bottom: 1px solid var(--border);
    transition: background 0.1s;
    border-left: 2px solid transparent;
  }

  .note-item:hover {
    background: var(--bg-hover);
  }

  .note-item.active {
    background: var(--bg-hover);
    border-left-color: var(--accent);
  }

  .note-title {
    font-size: 13px;
    font-weight: 500;
    color: var(--text);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .note-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .date {
    font-size: 11px;
    color: var(--text-muted);
  }

  .delete {
    opacity: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: var(--radius-sm);
    color: var(--text-muted);
    transition: opacity 0.1s, color 0.1s, background 0.1s;
  }

  .note-item:hover .delete {
    opacity: 1;
  }

  .delete:hover {
    color: var(--text);
    background: var(--bg-hover);
  }
</style>
