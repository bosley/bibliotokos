<script>
  import { createEventDispatcher, tick } from 'svelte'
  const dispatch = createEventDispatcher()

  export let note = null
  export let tags = []
  export let isDirty = false

  let title = note?.title ?? ''
  let content = note?.content ?? ''

  let tagMenuOpen = false
  let newTagInput = ''
  let tagMenuEl

  $: existingTagNames = tags.map(t => t.name)
  $: noteTags = note?.tags ?? []
  $: availableTags = existingTagNames.filter(n => !noteTags.includes(n))

  async function openTagMenu() {
    tagMenuOpen = true
    await tick()
    tagMenuEl?.querySelector('input')?.focus()
  }

  function closeTagMenu() {
    tagMenuOpen = false
    newTagInput = ''
  }

  function submitNewTag() {
    const name = newTagInput.trim()
    if (!name) return
    dispatch('addtag', name)
    newTagInput = ''
    closeTagMenu()
  }

  function handleNewTagKeydown(e) {
    if (e.key === 'Enter') submitNewTag()
    if (e.key === 'Escape') closeTagMenu()
  }

  function handleOutsideClick(e) {
    if (tagMenuEl && !tagMenuEl.contains(e.target)) {
      closeTagMenu()
    }
  }

  function handleWindowKeydown(e) {
    if ((e.metaKey || e.ctrlKey) && e.key === 's') {
      e.preventDefault()
      dispatch('save')
    }
  }
</script>

{#if note}
  <div class="editor">
    <div class="editor-header">
      <div class="title-row">
        <input
          class="title-input"
          type="text"
          bind:value={title}
          placeholder="Note title…"
          on:input={() => dispatch('change', { title })}
        />
        <button
          class="save-btn"
          class:dirty={isDirty}
          on:click={() => dispatch('save')}
          title="Save (⌘S)"
          disabled={!isDirty}
        >
          {#if isDirty}
            <span class="dot"></span>Save
          {:else}
            Saved
          {/if}
        </button>
      </div>
      <div class="tag-bar">
        {#each noteTags as tag}
          <span class="tag-chip">
            {tag}
            <button
              class="remove-tag"
              on:click={() => dispatch('removetag', tag)}
              title="Remove tag"
              aria-label="Remove tag {tag}"
            >×</button>
          </span>
        {/each}
        <div class="tag-menu-anchor" bind:this={tagMenuEl}>
          <button class="add-tag-btn" on:click={openTagMenu} title="Add tag">
            <svg viewBox="0 0 16 16" fill="currentColor" width="10" height="10">
              <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
            </svg>
            <span>Tag</span>
          </button>
          {#if tagMenuOpen}
            <div class="tag-dropdown" role="dialog" aria-label="Tag menu">
              <input
                class="tag-input"
                type="text"
                bind:value={newTagInput}
                placeholder="New tag…"
                on:keydown={handleNewTagKeydown}
              />
              {#if availableTags.length > 0}
                <div class="tag-options">
                  {#each availableTags as name}
                    <button class="tag-option" on:click={() => { dispatch('addtag', name); closeTagMenu() }}>
                      {name}
                    </button>
                  {/each}
                </div>
              {/if}
              {#if newTagInput.trim() && !existingTagNames.includes(newTagInput.trim())}
                <button class="tag-option create" on:click={submitNewTag}>
                  Create "{newTagInput.trim()}"
                </button>
              {/if}
            </div>
          {/if}
        </div>
      </div>
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

<svelte:window on:click={handleOutsideClick} on:keydown={handleWindowKeydown} />

<style>
  .editor {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .editor-header {
    padding: 20px 28px 14px;
    border-bottom: 1px solid var(--border);
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .title-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .title-input {
    flex: 1;
    font-size: 22px;
    font-weight: 600;
    color: var(--text);
    background: transparent;
    outline: none;
    font-family: var(--font-sans);
    padding: 0;
    min-width: 0;
  }

  .title-input::placeholder {
    color: var(--text-muted);
    font-weight: 400;
  }

  .save-btn {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 4px 10px;
    border-radius: var(--radius-sm, 5px);
    font-size: 11px;
    font-weight: 500;
    color: var(--text-muted);
    background: transparent;
    transition: color 0.15s, background 0.15s;
    white-space: nowrap;
    flex-shrink: 0;
    cursor: default;
  }

  .save-btn.dirty {
    color: var(--accent);
    background: color-mix(in srgb, var(--accent) 12%, transparent);
    cursor: pointer;
  }

  .save-btn.dirty:hover {
    background: color-mix(in srgb, var(--accent) 20%, transparent);
  }

  .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--accent);
    flex-shrink: 0;
  }

  .tag-bar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 5px;
  }

  .tag-chip {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    padding: 2px 6px 2px 8px;
    border-radius: 999px;
    font-size: 11px;
    font-weight: 500;
    color: var(--accent);
    background: color-mix(in srgb, var(--accent) 12%, transparent);
  }

  .remove-tag {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    font-size: 13px;
    line-height: 1;
    color: var(--accent);
    opacity: 0.7;
    transition: opacity 0.1s, background 0.1s;
    cursor: pointer;
  }

  .remove-tag:hover {
    opacity: 1;
    background: color-mix(in srgb, var(--accent) 20%, transparent);
  }

  .tag-menu-anchor {
    position: relative;
  }

  .add-tag-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 2px 8px;
    border-radius: 999px;
    font-size: 11px;
    font-weight: 500;
    color: var(--text-muted);
    background: var(--bg-hover);
    transition: color 0.15s, background 0.15s;
    cursor: pointer;
  }

  .add-tag-btn:hover {
    color: var(--text);
  }

  .tag-dropdown {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    z-index: 100;
    min-width: 180px;
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm, 6px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
    overflow: hidden;
  }

  .tag-input {
    display: block;
    width: 100%;
    padding: 8px 10px;
    font-size: 12px;
    color: var(--text);
    background: transparent;
    outline: none;
    border-bottom: 1px solid var(--border);
    font-family: var(--font-sans);
  }

  .tag-input::placeholder {
    color: var(--text-muted);
  }

  .tag-options {
    max-height: 160px;
    overflow-y: auto;
  }

  .tag-option {
    display: block;
    width: 100%;
    padding: 7px 10px;
    font-size: 12px;
    text-align: left;
    color: var(--text);
    transition: background 0.1s;
    cursor: pointer;
  }

  .tag-option:hover {
    background: var(--bg-hover);
  }

  .tag-option.create {
    color: var(--accent);
    font-style: italic;
    border-top: 1px solid var(--border);
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
