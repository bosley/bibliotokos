<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  export let note = null

  let title = note?.title ?? ''
  let content = note?.content ?? ''
</script>

{#if note}
  <div class="editor">
    <div class="editor-header">
      <input
        class="title-input"
        type="text"
        bind:value={title}
        placeholder="Note title…"
        on:input={() => dispatch('change', { title })}
      />
    </div>
    <textarea
      class="content-area"
      bind:value={content}
      placeholder="Start writing…"
      on:input={() => dispatch('change', { content })}
    ></textarea>
  </div>
{:else}
  <div class="empty-state">
    <p>Select a note or create a new one.</p>
  </div>
{/if}

<style>
  .editor {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .editor-header {
    padding: 20px 28px 16px;
    border-bottom: 1px solid var(--border);
    flex-shrink: 0;
  }

  .title-input {
    width: 100%;
    font-size: 22px;
    font-weight: 600;
    color: var(--text);
    background: transparent;
    outline: none;
    font-family: var(--font-sans);
    padding: 0;
  }

  .title-input::placeholder {
    color: var(--text-muted);
    font-weight: 400;
  }

  .content-area {
    flex: 1;
    padding: 20px 28px;
    font-size: 15px;
    line-height: 1.75;
    font-family: var(--font-reading);
    color: var(--text);
    background: transparent;
    outline: none;
    resize: none;
    width: 100%;
  }

  .content-area::placeholder {
    color: var(--text-muted);
  }

  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    font-size: 14px;
  }
</style>
