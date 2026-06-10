<script>
  import { query } from '../../stores/scripture.js'

  let inputValue = ''

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
    -webkit-app-region: no-drag;
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
</style>
