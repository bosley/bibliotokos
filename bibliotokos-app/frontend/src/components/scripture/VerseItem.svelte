<script>
  export let verse
</script>

<div class="verse-item">
  <span class="ref">{verse.book} {verse.chapter}:{verse.verse}</span>
  <span class="text">{verse.text}</span>
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
