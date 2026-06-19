<script>
  import { createEventDispatcher, tick } from 'svelte'
  import { get } from 'svelte/store'
  import Editor from '@toast-ui/editor'
  import '@toast-ui/editor/dist/toastui-editor.css'
  import '@toast-ui/editor/dist/theme/toastui-editor-dark.css'
  import { theme } from '../../stores/theme.js'

  const dispatch = createEventDispatcher()

  export let note = null
  export let tags = []
  export let isDirty = false
  export let passageError = ''

  let title = note?.title ?? ''

  let editor = null
  let editorRoot = null

  function interceptSave(e) {
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 's') {
      e.preventDefault()
      e.stopPropagation()
      dispatch('save')
    }
  }

  function initEditor(el) {
    el.addEventListener('keydown', interceptSave, { capture: true })
    editor = new Editor({
      el,
      initialValue: note?.content ?? '',
      initialEditType: 'wysiwyg',
      previewStyle: 'tab',
      height: '100%',
      usageStatistics: false,
      hideModeSwitch: true,
      autofocus: false,
      placeholder: 'Start writing…',
      theme: get(theme) === 'dark' ? 'dark' : 'default',
      toolbarItems: [
        ['heading', 'bold', 'italic', 'strike'],
        ['hr', 'quote'],
        ['ul', 'ol', 'task'],
        ['table', 'link'],
        ['code', 'codeblock'],
      ],
      events: {
        change: () => dispatch('change', { content: editor.getMarkdown() }),
      },
    })
    editorRoot = el
    return {
      destroy() {
        el.removeEventListener('keydown', interceptSave, { capture: true })
        editor?.destroy()
        editor = null
        editorRoot = null
      },
    }
  }

  $: if (editorRoot) {
    editorRoot
      .querySelector('.toastui-editor-defaultUI')
      ?.classList.toggle('toastui-editor-dark', $theme === 'dark')
  }

  let tagMenuOpen = false
  let newTagInput = ''
  let tagMenuEl

  let passageMenuOpen = false
  let passageInput = ''
  let passageMenuEl

  $: existingTagNames = tags.map(t => t.name)
  $: noteTags = note?.tags ?? []
  $: availableTags = existingTagNames.filter(n => !noteTags.includes(n))
  $: notePassages = note?.passages ?? []

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

  async function openPassageMenu() {
    passageMenuOpen = true
    dispatch('passageerrorclear')
    await tick()
    passageMenuEl?.querySelector('input')?.focus()
  }

  function closePassageMenu() {
    passageMenuOpen = false
    passageInput = ''
    dispatch('passageerrorclear')
  }

  function submitPassage() {
    const ref = passageInput.trim()
    if (!ref) return
    dispatch('addpassage', ref)
    passageInput = ''
  }

  function handlePassageKeydown(e) {
    if (e.key === 'Enter') submitPassage()
    if (e.key === 'Escape') closePassageMenu()
  }

  function handleOutsideClick(e) {
    if (tagMenuEl && !tagMenuEl.contains(e.target)) {
      closeTagMenu()
    }
    if (passageMenuEl && !passageMenuEl.contains(e.target)) {
      closePassageMenu()
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
      <div class="passage-bar">
        {#each notePassages as passage (passage.id)}
          <span class="passage-chip">
            {passage.reference}
            <button
              class="remove-passage"
              on:click={() => dispatch('removepassage', passage.id)}
              title="Unlink passage"
              aria-label="Unlink passage {passage.reference}"
            >×</button>
          </span>
        {/each}
        <div class="passage-menu-anchor" bind:this={passageMenuEl}>
          <button class="add-passage-btn" on:click={openPassageMenu} title="Link a Bible passage">
            <svg viewBox="0 0 16 16" fill="currentColor" width="10" height="10">
              <path d="M4.715 6.542 3.343 7.914a3 3 0 1 0 4.243 4.243l1.828-1.829A3 3 0 0 0 8.586 5.5L8 6.086a1 1 0 0 0-.154.199 2 2 0 0 1 .861 3.337L6.88 11.45a2 2 0 1 1-2.83-2.83l.793-.792a4 4 0 0 1-.128-1.287z"/>
              <path d="M6.586 4.672A3 3 0 0 0 7.414 9.5l.775-.776a2 2 0 0 1-.896-3.346L9.12 3.55a2 2 0 1 1 2.83 2.83l-.793.792c.112.42.155.855.128 1.287l1.372-1.372a3 3 0 1 0-4.243-4.243L6.586 4.672z"/>
            </svg>
            <span>Link passage</span>
          </button>
          {#if passageMenuOpen}
            <div class="passage-dropdown" role="dialog" aria-label="Link passage">
              <input
                class="passage-input"
                type="text"
                bind:value={passageInput}
                placeholder="e.g. John 1:1-3, Gen 1, Rom 8"
                on:keydown={handlePassageKeydown}
                on:input={() => dispatch('passageerrorclear')}
              />
              {#if passageError}
                <div class="passage-error">{passageError}</div>
              {:else}
                <div class="passage-hint">Press Enter to link</div>
              {/if}
            </div>
          {/if}
        </div>
      </div>
    </div>
    <div class="content-area" use:initEditor></div>
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

  .passage-bar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 5px;
  }

  .passage-chip {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    padding: 2px 6px 2px 8px;
    border-radius: var(--radius-sm, 5px);
    font-size: 11px;
    font-weight: 500;
    font-variant-numeric: tabular-nums;
    color: var(--text);
    background: var(--bg-hover);
    border: 1px solid var(--border);
  }

  .remove-passage {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    font-size: 13px;
    line-height: 1;
    color: var(--text-muted);
    opacity: 0.7;
    transition: opacity 0.1s, background 0.1s;
    cursor: pointer;
  }

  .remove-passage:hover {
    opacity: 1;
    background: color-mix(in srgb, currentColor 15%, transparent);
  }

  .passage-menu-anchor {
    position: relative;
  }

  .add-passage-btn {
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

  .add-passage-btn:hover {
    color: var(--text);
  }

  .passage-dropdown {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    z-index: 100;
    min-width: 220px;
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm, 6px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
    overflow: hidden;
  }

  .passage-input {
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

  .passage-input::placeholder {
    color: var(--text-muted);
  }

  .passage-hint {
    padding: 6px 10px;
    font-size: 11px;
    color: var(--text-muted);
  }

  .passage-error {
    padding: 6px 10px;
    font-size: 11px;
    color: #e5534b;
  }

  .content-area {
    flex: 1;
    overflow: hidden;
    min-height: 0;
  }

  .content-area :global(.toastui-editor-defaultUI) {
    border: none;
    border-radius: 0;
  }

  .content-area :global(.toastui-editor-defaultUI-toolbar) {
    background: var(--bg-surface);
    border-bottom: 1px solid var(--border);
    border-radius: 0;
  }

  .content-area :global(.toastui-editor-main),
  .content-area :global(.toastui-editor-md-container),
  .content-area :global(.toastui-editor-ww-container),
  .content-area :global(.toastui-editor-md-preview) {
    background: transparent;
  }

  .content-area :global(.toastui-editor-mode-switch) {
    display: none;
  }

  .content-area :global(.toastui-editor-md-tab-container) {
    display: none;
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
