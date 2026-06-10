<script>
  import VerseItem from './VerseItem.svelte'

  export let verses = []
  export let loading = false
  export let query = ''
</script>

<div class="verse-list">
  {#if loading}
    <div class="state-msg">Loading…</div>
  {:else if !query}
    <div class="state-msg empty">
      <p>Enter a scripture reference to begin.</p>
      <p class="hint">Try <em>John 3:16</em>, <em>Gen 1</em>, <em>Romans 8:28-39</em></p>
    </div>
  {:else if verses.length === 0}
    <div class="state-msg">No results for <em>{query}</em></div>
  {:else}
    {#each verses as verse (`${verse.book}-${verse.chapter}-${verse.verse}-${verse.version}`)}
      <VerseItem {verse} />
    {/each}
  {/if}
</div>

<style>
  .verse-list {
    flex: 1;
    overflow-y: auto;
    padding: 20px 24px;
  }

  .state-msg {
    color: var(--text-muted);
    font-size: 14px;
    text-align: center;
    padding: 60px 24px;
    line-height: 1.8;
  }

  .state-msg .hint {
    margin-top: 6px;
    font-size: 13px;
  }

  .state-msg em {
    color: var(--accent);
    font-style: normal;
  }
</style>
