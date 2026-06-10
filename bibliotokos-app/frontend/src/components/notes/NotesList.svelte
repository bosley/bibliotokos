<script>
  import NoteItem from './NoteItem.svelte'

  export let noteHeaders = []
  export let activeNoteId = null
</script>

<div class="notes-list">
  {#if noteHeaders.length === 0}
    <div class="empty">No notes yet. Create one to get started.</div>
  {:else}
    {#each noteHeaders as header (header.id)}
      <NoteItem
        note={header}
        active={header.id === activeNoteId}
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
