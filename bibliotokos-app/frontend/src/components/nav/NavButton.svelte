<script>
  export let title = ''
  export let active = false
</script>

<button class:active {title} aria-label={title} on:click>
  <slot />
</button>

<style>
  button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border-radius: var(--radius-sm);
    color: var(--text-muted);
    transition: color 0.15s, background 0.15s;
    flex-shrink: 0;
    -webkit-app-region: no-drag;
  }

  button:hover {
    color: var(--text);
    background: var(--bg-hover);
  }

  button.active {
    color: var(--accent);
  }
</style>
