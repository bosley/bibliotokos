<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
</script>

<div
  class="divider"
  role="presentation"
  on:mousedown|preventDefault={e => dispatch('dragstart', e.clientX)}
/>

<style>
  .divider {
    width: 4px;
    background: var(--divider);
    cursor: col-resize;
    flex-shrink: 0;
    transition: background 0.15s;
  }

  .divider:hover {
    background: var(--accent);
  }
</style>
