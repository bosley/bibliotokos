<script>
  import { query, scrollLock } from '../../stores/scripture.js'

  let inputValue = ''

  function toggleLock() {
    scrollLock.update(v => !v)
  }

  function handleKeydown(e) {
    if (e.key === 'Enter' && inputValue.trim()) {
      query.set(inputValue.trim())
    }
    if (e.key === 'Escape') {
      inputValue = ''
      query.set('')
    }
  }
</script>

<div class="search-wrap">
  <svg class="icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5">
    <circle cx="7" cy="7" r="4.5" />
    <line x1="10.5" y1="10.5" x2="14" y2="14" />
  </svg>
  <input
    type="text"
    placeholder="John 3:16 · Gen 1 · Romans 8:28-39…"
    bind:value={inputValue}
    on:keydown={handleKeydown}
    spellcheck="false"
    autocomplete="off"
  />
  <button
    class="lock-btn"
    class:locked={$scrollLock}
    on:click={toggleLock}
    title={$scrollLock ? 'Scroll lock: locked' : 'Scroll lock: unlocked'}
    aria-label={$scrollLock ? 'Unlock scroll sync' : 'Lock scroll sync'}
  >
    {#if $scrollLock}
      <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
        <path d="M8 1a2 2 0 0 1 2 2v4H6V3a2 2 0 0 1 2-2zm3 6V3a3 3 0 0 0-6 0v4a2 2 0 0 0-2 2v5a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2z"/>
      </svg>
    {:else}
      <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
        <path d="M11 1a2 2 0 0 0-2 2v4a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V9a2 2 0 0 1 2-2h5V3a3 3 0 0 1 6 0v4a.5.5 0 0 1-1 0V3a2 2 0 0 0-2-2z"/>
      </svg>
    {/if}
  </button>
</div>

<style>
  .search-wrap {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 0 12px;
    height: 34px;
    max-width: 480px;
    transition: border-color 0.15s;
  }

  .search-wrap:focus-within {
    border-color: var(--accent);
  }

  .icon {
    width: 14px;
    height: 14px;
    color: var(--text-muted);
    flex-shrink: 0;
  }

  input {
    flex: 1;
    font-size: 13px;
    outline: none;
    background: transparent;
    color: var(--text);
  }

  input::placeholder {
    color: var(--text-muted);
  }

  .lock-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    flex-shrink: 0;
    border: none;
    background: transparent;
    border-radius: var(--radius-sm, 4px);
    color: var(--text-muted);
    cursor: pointer;
    transition: color 0.15s, background 0.15s;
  }

  .lock-btn:hover {
    color: var(--text);
    background: var(--bg-hover);
  }

  .lock-btn.locked {
    color: var(--accent);
  }
</style>
