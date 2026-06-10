<script>
  import { onMount, onDestroy } from 'svelte'
  import {
    ListNotes,
    GetNote,
    GetTags,
    SaveNote,
    DeleteNote,
    CreateTag,
    DeleteTag,
    LinkPassage,
    UnlinkPassage,
  } from '../../bindings/bibliotokos/services/notes/notesservice.js'
  import NotesSidebar from '../components/notes/NotesSidebar.svelte'
  import NoteEditor from '../components/notes/NoteEditor.svelte'
  import PaneDivider from '../components/scripture/PaneDivider.svelte'

  let noteHeaders = []
  let tags = []
  let activeNote = null
  let activeNoteId = null
  let selectedTagId = null
  let isDirty = false
  let saveTimer = null
  let passageError = ''

  let sidebarWidth = 260
  let collapsed = false

  onMount(async () => {
    const [headers, tagList] = await Promise.all([
      ListNotes().catch(() => []),
      GetTags().catch(() => []),
    ])
    noteHeaders = headers ?? []
    tags = tagList ?? []
    const requestedId = new URLSearchParams(window.location.search).get('note')
    if (requestedId && noteHeaders.some(h => h.id === requestedId)) {
      await selectNote(requestedId)
    } else if (noteHeaders.length > 0) {
      await selectNote(noteHeaders[0].id)
    }
  })

  onDestroy(() => {
    clearTimeout(saveTimer)
  })

  async function flushSave() {
    clearTimeout(saveTimer)
    if (!activeNote || !isDirty) return
    await SaveNote(activeNote).catch(() => {})
    noteHeaders = await ListNotes().catch(() => noteHeaders)
    isDirty = false
  }

  function scheduleSave() {
    isDirty = true
    clearTimeout(saveTimer)
    saveTimer = setTimeout(flushSave, 3000)
  }

  async function selectNote(id) {
    await flushSave()
    const fetched = await GetNote(id).catch(() => null)
    activeNote = fetched
    activeNoteId = id
    isDirty = false
    passageError = ''
  }

  async function createNote() {
    await flushSave()
    const note = {
      id: crypto.randomUUID(),
      title: 'Untitled Note',
      content: '',
      updatedAt: '',
      tags: [],
      passages: [],
    }
    await SaveNote(note)
    noteHeaders = await ListNotes().catch(() => noteHeaders)
    const fetched = await GetNote(note.id).catch(() => null)
    activeNote = fetched
    activeNoteId = note.id
    isDirty = false
  }

  function handleEditorChange(id, changes) {
    if (!activeNote || activeNote.id !== id) return
    activeNote = { ...activeNote, ...changes }
    scheduleSave()
  }

  async function deleteNote(id) {
    if (activeNoteId === id) {
      clearTimeout(saveTimer)
      isDirty = false
    }
    await DeleteNote(id)
    noteHeaders = await ListNotes().catch(() => noteHeaders)
    if (activeNoteId === id) {
      const remaining = noteHeaders
      if (remaining.length > 0) {
        await selectNote(remaining[0].id)
      } else {
        activeNote = null
        activeNoteId = null
      }
    }
  }

  async function addTagToNote(name) {
    if (!activeNote) return
    let tag = tags.find(t => t.name === name)
    if (!tag) {
      tag = await CreateTag(name).catch(() => null)
      if (tag) tags = [...tags, tag]
    }
    if (tag && !activeNote.tags.includes(name)) {
      activeNote = { ...activeNote, tags: [...activeNote.tags, name] }
      await SaveNote(activeNote).catch(() => {})
      noteHeaders = await ListNotes().catch(() => noteHeaders)
    }
  }

  async function deleteTag(id) {
    await DeleteTag(id).catch(() => {})
    tags = await GetTags().catch(() => tags)
    noteHeaders = await ListNotes().catch(() => noteHeaders)
    if (selectedTagId && !tags.find(t => t.name === selectedTagId)) {
      selectedTagId = null
    }
    if (activeNote) {
      const refreshed = await GetNote(activeNote.id).catch(() => activeNote)
      activeNote = refreshed
    }
  }

  async function removeTagFromNote(name) {
    if (!activeNote) return
    activeNote = { ...activeNote, tags: activeNote.tags.filter(t => t !== name) }
    await SaveNote(activeNote).catch(() => {})
    noteHeaders = await ListNotes().catch(() => noteHeaders)
  }

  async function addPassageToNote(ref) {
    if (!activeNote) return
    try {
      const passage = await LinkPassage(activeNote.id, ref)
      if (!activeNote.passages.some(p => p.id === passage.id)) {
        activeNote = { ...activeNote, passages: [...activeNote.passages, passage] }
      }
      passageError = ''
    } catch (err) {
      passageError = 'Unknown passage'
    }
  }

  async function removePassageFromNote(passageId) {
    if (!activeNote) return
    await UnlinkPassage(passageId).catch(() => {})
    activeNote = { ...activeNote, passages: activeNote.passages.filter(p => p.id !== passageId) }
  }

  function startResize(startX) {
    const startWidth = sidebarWidth
    function onMove(e) {
      sidebarWidth = Math.max(160, Math.min(520, startWidth + (e.clientX - startX)))
    }
    function onUp() {
      document.removeEventListener('mousemove', onMove)
      document.removeEventListener('mouseup', onUp)
    }
    document.addEventListener('mousemove', onMove)
    document.addEventListener('mouseup', onUp)
  }

  function handleBeforeUnload() {
    if (activeNote && isDirty) {
      SaveNote(activeNote).catch(() => {})
    }
  }
</script>

<svelte:window on:beforeunload={handleBeforeUnload} />

<div class="notes-view">
  <header class="titlebar">Notes</header>
  <div class="body" style="--sidebar-w: {sidebarWidth}px">
    {#if !collapsed}
      <NotesSidebar
        {noteHeaders}
        {tags}
        {activeNoteId}
        {selectedTagId}
        on:select={e => selectNote(e.detail)}
        on:create={createNote}
        on:delete={e => deleteNote(e.detail)}
        on:selecttag={e => { selectedTagId = e.detail }}
        on:deletetag={e => deleteTag(e.detail)}
        on:collapse={() => { collapsed = true }}
      />
      <PaneDivider on:dragstart={e => startResize(e.detail)} />
    {:else}
      <button class="expand-btn" on:click={() => { collapsed = false }} title="Expand sidebar">
        <svg viewBox="0 0 16 16" fill="currentColor" width="12" height="12">
          <path d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"/>
        </svg>
      </button>
    {/if}
    {#key activeNoteId}
      <NoteEditor
        note={activeNote}
        {tags}
        {isDirty}
        {passageError}
        on:change={e => handleEditorChange(activeNoteId, e.detail)}
        on:save={flushSave}
        on:addtag={e => addTagToNote(e.detail)}
        on:removetag={e => removeTagFromNote(e.detail)}
        on:addpassage={e => addPassageToNote(e.detail)}
        on:removepassage={e => removePassageFromNote(e.detail)}
        on:passageerrorclear={() => { passageError = '' }}
      />
    {/key}
  </div>
</div>

<style>
  .notes-view {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .titlebar {
    height: 50px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    padding-left: 88px;
    font-size: 13px;
    font-weight: 600;
    color: var(--text-muted);
    border-bottom: 1px solid var(--border);
    background: var(--bg-surface);
    -webkit-app-region: drag;
  }

  .body {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .expand-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    flex-shrink: 0;
    background: var(--bg-surface);
    border-right: 1px solid var(--border);
    color: var(--text-muted);
    cursor: pointer;
    transition: color 0.15s, background 0.15s;
  }

  .expand-btn:hover {
    color: var(--text);
    background: var(--bg-hover);
  }
</style>
