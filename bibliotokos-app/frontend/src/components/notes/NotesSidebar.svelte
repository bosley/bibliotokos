<script>
  import { createEventDispatcher } from 'svelte'
  import NotesList from './NotesList.svelte'

  export let notes = []
  export let activeNoteId = null

  const dispatch = createEventDispatcher()
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <span class="label">Notes</span>
    <button on:click={() => dispatch('create')} title="New note" aria-label="New note">
      <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
        <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
      </svg>
    </button>
  </div>
  <NotesList {notes} {activeNoteId} on:select on:delete />
</aside>

<style>
  .sidebar {
    width: var(--sidebar-width);
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
