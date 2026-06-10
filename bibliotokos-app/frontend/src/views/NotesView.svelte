<script>
  import { onMount } from 'svelte'
  import { GetNotes, SaveNote, DeleteNote } from '../../bindings/bibliotokos/services/notes/notesservice.js'
  import NotesSidebar from '../components/notes/NotesSidebar.svelte'
  import NoteEditor from '../components/notes/NoteEditor.svelte'

  let notes = []
  let activeNoteId = null

  $: activeNote = notes.find(n => n.id === activeNoteId) ?? null

  onMount(async () => {
    notes = await GetNotes().catch(() => [])
    activeNoteId = notes[0]?.id ?? null
  })

  async function selectNote(id) {
    activeNoteId = id
  }

  async function createNote() {
    const note = {
      id: crypto.randomUUID(),
      title: 'Untitled Note',
      content: '',
      updatedAt: '',
    }
    await SaveNote(note)
    notes = await GetNotes().catch(() => notes)
    activeNoteId = note.id
  }

  async function updateNote(id, changes) {
    const note = notes.find(n => n.id === id)
    if (!note) return
    await SaveNote({ ...note, ...changes })
    notes = await GetNotes().catch(() => notes)
  }

  async function deleteNote(id) {
    await DeleteNote(id)
    notes = await GetNotes().catch(() => notes)
    if (activeNoteId === id) {
      activeNoteId = notes[0]?.id ?? null
    }
  }
</script>

<div class="notes-view">
  <header class="titlebar">Notes</header>
  <div class="body">
    <NotesSidebar
      {notes}
      {activeNoteId}
      on:select={e => selectNote(e.detail)}
      on:create={createNote}
      on:delete={e => deleteNote(e.detail)}
    />
    {#key activeNoteId}
      <NoteEditor
        note={activeNote}
        on:change={e => updateNote(activeNoteId, e.detail)}
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
</style>
