<script>
  import { Events } from '@wailsio/runtime'

  export let linkedNotes = []

  let expanded = false

  $: if (linkedNotes.length === 0) expanded = false

  function openNote(noteId) {
    Events.Emit('open-notes', noteId)
  }
</script>

{#if linkedNotes.length > 0}
  <div class="linked-notes-bar">
    <button class="toggle" on:click={() => { expanded = !expanded }} aria-expanded={expanded}>
      <svg viewBox="0 0 16 16" fill="currentColor" width="11" height="11">
        <path d="M2.5 1A1.5 1.5 0 0 0 1 2.5v11A1.5 1.5 0 0 0 2.5 15h6.086a1.5 1.5 0 0 0 1.06-.44l4.915-4.914A1.5 1.5 0 0 0 15 8.586V2.5A1.5 1.5 0 0 0 13.5 1h-11zm6 8.5a1 1 0 0 1 1-1h4.396l-5.396 5.396V9.5z"/>
      </svg>
      <span>{linkedNotes.length} linked {linkedNotes.length === 1 ? 'note' : 'notes'}</span>
      <svg class="chevron" class:expanded viewBox="0 0 16 16" fill="currentColor" width="10" height="10">
        <path d="M7.646 4.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1-.708.708L8 5.707l-5.646 5.647a.5.5 0 0 1-.708-.708l6-6z"/>
      </svg>
    </button>
    {#if expanded}
      <div class="note-list">
        {#each linkedNotes as ln (ln.noteId)}
          <button class="linked-note" on:click={() => openNote(ln.noteId)} title="Open note">
            <span class="note-title">{ln.title || 'Untitled Note'}</span>
            <span class="refs">
              {#each ln.references as ref}
                <span class="ref-chip">{ref}</span>
              {/each}
            </span>
          </button>
        {/each}
      </div>
    {/if}
  </div>
{/if}

<style>
  .linked-notes-bar {
    flex-shrink: 0;
    border-top: 1px solid var(--border);
    background: var(--bg-surface);
  }

  .toggle {
    display: flex;
    align-items: center;
    gap: 7px;
    width: 100%;
    padding: 6px 16px;
    font-size: 11px;
    font-weight: 500;
    color: var(--text-muted);
    background: transparent;
    cursor: pointer;
    transition: color 0.15s, background 0.15s;
  }

  .toggle:hover {
    color: var(--text);
    background: var(--bg-hover);
  }

  .chevron {
    transition: transform 0.15s;
  }

  .chevron.expanded {
    transform: rotate(180deg);
  }

  .note-list {
    max-height: 200px;
    overflow-y: auto;
    border-top: 1px solid var(--border);
  }

  .linked-note {
    display: flex;
    align-items: baseline;
    gap: 10px;
    width: 100%;
    padding: 7px 16px;
    text-align: left;
    background: transparent;
    cursor: pointer;
    transition: background 0.1s;
  }

  .linked-note:hover {
    background: var(--bg-hover);
  }

  .note-title {
    font-size: 12px;
    font-weight: 500;
    color: var(--text);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex-shrink: 0;
    max-width: 50%;
  }

  .refs {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    overflow: hidden;
  }

  .ref-chip {
    padding: 1px 6px;
    border-radius: 999px;
    font-size: 10px;
    font-weight: 500;
    font-variant-numeric: tabular-nums;
    color: var(--accent);
    background: color-mix(in srgb, var(--accent) 12%, transparent);
    white-space: nowrap;
  }
</style>
