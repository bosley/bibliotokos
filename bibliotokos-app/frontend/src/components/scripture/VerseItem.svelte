<script>
  export let verse

  let copied = null

  function stripHtml(html) {
    if (!html) return ''
    const doc = new DOMParser().parseFromString(html, 'text/html')
    return doc.body.textContent || ''
  }

  $: plainText = stripHtml(verse.text)

  function copy(type) {
    const value = type === 'ref'
      ? `${verse.book} ${verse.chapter}:${verse.verse}`
      : plainText
    navigator.clipboard.writeText(value)
    copied = type
    setTimeout(() => copied = null, 800)
  }
</script>

<div class="verse-item">
  <div class="copy-pills">
    <button 
      class="pill" 
      class:copied={copied === 'ref'}
      on:click={() => copy('ref')}
    >ref</button>
    <button 
      class="pill" 
      class:copied={copied === 'txt'}
      on:click={() => copy('txt')}
    >txt</button>
  </div>
  <span class="ref">{verse.book} {verse.chapter}:{verse.verse}</span>
  <span class="text">{@html verse.text}</span>
</div>

<style>
  .verse-item {
    display: flex;
    gap: 16px;
    padding: 5px 0;
    line-height: 1.65;
    border-bottom: 1px solid transparent;
  }

  .verse-item:hover {
    border-bottom-color: var(--border);
  }

  .copy-pills {
    display: flex;
    gap: 4px;
    opacity: 0;
    transition: opacity 0.15s;
    flex-shrink: 0;
    padding-top: 2px;
  }

  .verse-item:hover .copy-pills {
    opacity: 1;
  }

  .pill {
    font-size: 9px;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.02em;
    padding: 2px 5px;
    border-radius: 3px;
    border: none;
    cursor: pointer;
    background: var(--bg-alt);
    color: var(--text-muted);
    transition: background 0.15s, color 0.15s;
  }

  .pill:hover {
    background: var(--accent);
    color: white;
  }

  .pill.copied {
    background: var(--accent);
    color: white;
  }

  .ref {
    color: var(--text-muted);
    font-size: 11px;
    font-variant-numeric: tabular-nums;
    flex-shrink: 0;
    padding-top: 3px;
    min-width: 88px;
    text-align: right;
  }

  .text {
    font-size: 15px;
    font-family: var(--font-reading);
    color: var(--text);
  }
</style>
