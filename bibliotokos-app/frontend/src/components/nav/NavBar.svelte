<script>
  import { Events } from '@wailsio/runtime'
  import NavButton from './NavButton.svelte'
  import SearchInput from './SearchInput.svelte'
  import ThemeToggle from './ThemeToggle.svelte'

  let showInfo = false
  let showHelp = false
  let copied = ''

  function openNotes() {
    Events.Emit('open-notes', null)
  }

  function openInfo() {
    showInfo = true
    showHelp = false
  }

  function openHelp() {
    showHelp = true
    showInfo = false
  }

  function closeModals() {
    showInfo = false
    showHelp = false
    copied = ''
  }

  async function copyQuery(q) {
    try {
      await navigator.clipboard.writeText(q)
      copied = q
      setTimeout(() => { if (copied === q) copied = '' }, 1500)
    } catch (_) {
      // clipboard may fail in some envs; ignore
    }
  }

  function handleKey(e) {
    if (e.key === 'Escape') closeModals()
  }
</script>

<svelte:window on:keydown={handleKey} />

<nav class="navbar">
  <div class="left">
    <NavButton title="Scripture">
      <svg viewBox="0 0 16 16" fill="currentColor" width="16" height="16">
        <path d="M1 2.828c.885-.37 2.154-.769 3.388-.893 1.33-.134 2.458.063 3.112.752v9.746c-.935-.53-2.12-.603-3.213-.493-1.187.12-2.147.35-2.287.498V2.828zm7.5-.141c.654-.689 1.782-.886 3.112-.752 1.234.124 2.503.523 3.388.893v9.985c-.14-.148-1.1-.378-2.287-.498-1.093-.11-2.278-.037-3.213.493V2.687zM8 1.783C7.015.936 5.587.81 4.287.94c-1.514.153-3.042.672-3.994 1.105A.5.5 0 0 0 0 2.5v11a.5.5 0 0 0 .707.455c.882-.4 2.303-.881 3.68-1.02 1.409-.142 2.59.087 3.223.877a.5.5 0 0 0 .78 0c.633-.79 1.814-1.019 3.222-.877 1.378.139 2.8.62 3.681 1.02A.5.5 0 0 0 16 13.5v-11a.5.5 0 0 0-.293-.455c-.952-.433-2.48-.952-3.994-1.105C10.413.809 8.985.936 8 1.783z"/>
      </svg>
    </NavButton>
    <NavButton title="Notes" on:click={openNotes}>
      <svg viewBox="0 0 16 16" fill="currentColor" width="16" height="16">
        <path d="M2.5 1A1.5 1.5 0 0 0 1 2.5v11A1.5 1.5 0 0 0 2.5 15h6.086a1.5 1.5 0 0 0 1.06-.44l4.915-4.914A1.5 1.5 0 0 0 15 8.586V2.5A1.5 1.5 0 0 0 13.5 1h-11zm6 8.5a1 1 0 0 1 1-1h4.396l-5.396 5.396V9.5z"/>
      </svg>
    </NavButton>
  </div>

  <div class="center">
    <SearchInput />
  </div>

  <div class="right">
    <ThemeToggle />
    <button class="icon-btn" on:click={openInfo} title="Info" aria-label="Info">
      <svg viewBox="0 0 16 16" fill="currentColor" width="16" height="16">
        <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
        <path d="M8.93 6.588-.84-.84.84-.84.84.84-.84.84z"/>
        <path d="M8 10.5a.5.5 0 0 1-.5-.5v-3a.5.5 0 0 1 1 0v3a.5.5 0 0 1-.5.5z"/>
      </svg>
    </button>
    <button class="icon-btn" on:click={openHelp} title="Help" aria-label="Help">
      <svg viewBox="0 0 16 16" fill="currentColor" width="16" height="16">
        <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zM8 1a6 6 0 0 0-6 6 6 6 0 0 0 3.5 5.4.75.75 0 1 0 .75-1.3A4.5 4.5 0 1 1 8 3.5 4.5 4.5 0 0 1 11.5 8a.75.75 0 0 0 1.5 0A6 6 0 0 0 8 1z"/>
        <path d="M8 12a1 1 0 1 0 0-2 1 1 0 0 0 0 2z"/>
      </svg>
    </button>
  </div>
</nav>

{#if showInfo}
  <div class="modal-backdrop" on:click={closeModals}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <span>Info</span>
        <button class="close-btn" on:click={closeModals}>×</button>
      </div>
      <div class="modal-body">
        <p><strong>Developer:</strong> Insula Labs, LLC</p>
        <p><strong>Contact Email:</strong> bosley@insulalabs.com</p>
        <p><strong>Bibles Sourced from:</strong></p>
        <ul>
          <li><a href="https://sacred-texts.com/bib/osrc/index.htm" target="_blank" rel="noopener">https://sacred-texts.com/bib/osrc/index.htm</a></li>
          <li><a href="https://github.com/seven1m/open-bibles" target="_blank" rel="noopener">https://github.com/seven1m/open-bibles</a></li>
        </ul>
        <p>All Bible texts are either Public Domain or licensed under Creative Commons (CC).</p>
      </div>
    </div>
  </div>
{/if}

{#if showHelp}
  <div class="modal-backdrop" on:click={closeModals}>
    <div class="modal help-modal" on:click|stopPropagation>
      <div class="modal-header">
        <span>Help</span>
        <button class="close-btn" on:click={closeModals}>×</button>
      </div>
      <div class="modal-body">
        <h4>Bible Search Queries</h4>
        <p class="help-sub">Type these in the top search bar (examples are copyable):</p>

        <div class="query-list">
          {#each [
            'John 3:16',
            'Gen 1',
            'Psalm 23:1-6',
            'Rom 8:28',
            '1cor 13',
            'mk 1:1',
            'song 2',
            'Gen 1 - Rev 22'
          ] as q}
            <div class="query-row">
              <code>{q}</code>
              <button class="copy-pill" on:click={() => copyQuery(q)}>
                {copied === q ? 'Copied!' : 'Copy'}
              </button>
            </div>
          {/each}
        </div>

        <p class="help-note">Queries support book names/aliases, chapters, verses, ranges, and cross-book ranges with " - ".</p>

        <h4>Linked Passages</h4>
        <p>In the Notes window, use the "Link passage" button to associate a Bible reference with a note. When you view that passage (or section) in the main Scripture window, any linked notes will appear in a bar at the bottom of the pane. Click a linked note to jump back to it.</p>
      </div>
    </div>
  </div>
{/if}

<style>
  .navbar {
    height: var(--nav-height);
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 16px 0 80px;
    border-bottom: 1px solid var(--border);
    background: var(--bg);
    flex-shrink: 0;
    -webkit-app-region: drag;
  }

  .left {
    display: flex;
    align-items: center;
    gap: 4px;
    -webkit-app-region: no-drag;
  }

  .center {
    flex: 1;
    display: flex;
    justify-content: center;
    -webkit-app-region: no-drag;
  }

  .right {
    display: flex;
    align-items: center;
    gap: 4px;
    -webkit-app-region: no-drag;
  }

  .icon-btn {
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
    background: transparent;
    border: none;
    cursor: pointer;
  }

  .icon-btn:hover {
    color: var(--text);
    background: var(--bg-hover);
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm, 6px);
    width: 460px;
    max-width: 92vw;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
    color: var(--text);
  }

  .help-modal {
    width: 520px;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border);
    font-weight: 600;
    font-size: 14px;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    font-size: 20px;
    line-height: 1;
    cursor: pointer;
    padding: 0 4px;
  }

  .close-btn:hover {
    color: var(--text);
  }

  .modal-body {
    padding: 16px;
    font-size: 13px;
    line-height: 1.5;
  }

  .modal-body p {
    margin: 8px 0;
  }

  .modal-body ul {
    margin: 8px 0 12px 20px;
    padding: 0;
  }

  .modal-body a {
    color: var(--accent);
    text-decoration: none;
  }

  .modal-body a:hover {
    text-decoration: underline;
  }

  .help-sub {
    color: var(--text-muted);
    font-size: 12px;
  }

  .query-list {
    margin: 12px 0;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .query-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--bg-hover);
    border-radius: 4px;
    padding: 6px 10px;
  }

  .query-row code {
    font-family: var(--font-mono, monospace);
    font-size: 12px;
  }

  .copy-pill {
    font-size: 11px;
    padding: 2px 10px;
    border-radius: 999px;
    background: var(--accent);
    color: white;
    border: none;
    cursor: pointer;
    transition: opacity 0.1s;
  }

  .copy-pill:hover {
    opacity: 0.85;
  }

  .help-note {
    font-size: 12px;
    color: var(--text-muted);
    margin-top: 8px;
  }

  h4 {
    margin: 16px 0 8px;
    font-size: 13px;
    color: var(--text-muted);
  }
</style>
