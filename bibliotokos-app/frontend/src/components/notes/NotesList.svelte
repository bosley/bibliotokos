<script>
  import NoteItem from './NoteItem.svelte'

  export let notes = []
  export let activeNoteId = null
</script>

<div class="notes-list">
  {#if notes.length === 0}
    <div class="empty">No notes yet. Create one to get started.</div>
  {:else}
    {#each notes as note (note.id)}
      <NoteItem
        {note}
        active={note.id === activeNoteId}
        on:select
        on:delete
      />
    {/each}
  {/if}
</div>

<style>
  .notes-list {
    flex: 1;
    overflow-y: auto;
  }

  .empty {
    padding: 32px 16px;
    color: var(--text-muted);
    font-size: 13px;
    text-align: center;
    line-height: 1.6;
  }
</style>
