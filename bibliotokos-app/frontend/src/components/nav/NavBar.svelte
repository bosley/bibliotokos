<script>
  import { Events, Window } from '@wailsio/runtime'
  import NavButton from './NavButton.svelte'
  import SearchInput from './SearchInput.svelte'
  import ThemeToggle from './ThemeToggle.svelte'

  let showInfo = false
  let showHelp = false
  let copied = ''
  let copiedUrl = ''

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
    } catch (_) {}
  }

  async function copyUrl(url) {
    try {
      await navigator.clipboard.writeText(url)
      copiedUrl = url
      setTimeout(() => { if (copiedUrl === url) copiedUrl = '' }, 1500)
    } catch (_) {}
  }

  function handleKey(e) {
    if (e.key === 'Escape') closeModals()
  }

  function handleTitlebarDblClick() {
    Window.ToggleMaximise()
  }
</script>

<svelte:window on:keydown={handleKey} />

<div class="titlebar" on:dblclick={handleTitlebarDblClick}></div>
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
        <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
        <path d="M5.255 5.786a.237.237 0 0 0 .241.247h.825c.138 0 .248-.113.266-.25.09-.656.54-1.134 1.342-1.134.686 0 1.314.343 1.314 1.168 0 .635-.374.927-.965 1.371-.673.489-1.206 1.06-1.168 1.987l.003.217a.25.25 0 0 0 .25.246h.811a.25.25 0 0 0 .25-.25v-.105c0-.718.273-.927 1.01-1.486.609-.463 1.244-.977 1.244-2.056 0-1.511-1.276-2.241-2.673-2.241-1.267 0-2.655.59-2.75 2.286zm1.557 5.763c0 .533.425.927 1.01.927.609 0 1.028-.394 1.028-.927 0-.552-.42-.94-1.029-.94-.584 0-1.009.388-1.009.94z"/>
      </svg>
    </button>
  </div>
</nav>

{#if showInfo}
  <button class="modal-backdrop" on:click={closeModals} aria-label="Close modal">
    <div class="modal info-modal" role="dialog" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <span>About Bibliotokos</span>
        <button class="close-btn" on:click={closeModals}>×</button>
      </div>
      <div class="modal-body">
        <div class="info-section">
          <div class="info-row">
            <span class="info-label">Developer</span>
            <span class="info-value">Insula Labs, LLC</span>
          </div>
          <div class="info-row">
            <span class="info-label">Contact</span>
            <a href="mailto:bosley@insulalabs.com" class="info-value info-link">bosley@insulalabs.com</a>
          </div>
        </div>

        <div class="info-card">
          <div class="info-card-header">ENV &amp; VENV</div>
          <p>The <em>Emergent Noetic Version</em> (ENV) and <em>Vulgate Emergent Noetic Version</em> (VENV) are exclusive to Bibliotokos. Created by iteratively translating original-language texts through multiple LLMs to preserve meaning as closely as possible.</p>
          <div class="env-details">
            <div class="env-item"><span class="env-tag">ENV</span> Hebrew Tanach + Greek NT</div>
            <div class="env-item"><span class="env-tag">VENV</span> Latin Vulgate</div>
          </div>
          <p class="info-note">Source texts are available in the translator.</p>
        </div>

        <div class="info-card">
          <div class="info-card-header">Other Bible Texts</div>
          <p>Sourced from the following (Public Domain / Creative Commons):</p>
          <div class="source-links">
            <button class="source-btn" class:copied={copiedUrl === 'https://sacred-texts.com/bib/osrc/index.htm'} on:click={() => copyUrl('https://sacred-texts.com/bib/osrc/index.htm')}>
              <span>sacred-texts.com</span>
              <span class="copy-hint">{copiedUrl === 'https://sacred-texts.com/bib/osrc/index.htm' ? 'Copied!' : 'Copy'}</span>
            </button>
            <button class="source-btn" class:copied={copiedUrl === 'https://github.com/seven1m/open-bibles'} on:click={() => copyUrl('https://github.com/seven1m/open-bibles')}>
              <span>github.com/seven1m/open-bibles</span>
              <span class="copy-hint">{copiedUrl === 'https://github.com/seven1m/open-bibles' ? 'Copied!' : 'Copy'}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </button>
{/if}

{#if showHelp}
  <button class="modal-backdrop" on:click={closeModals} aria-label="Close modal">
    <div class="modal help-modal" role="dialog" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <span>Quick Reference</span>
        <button class="close-btn" on:click={closeModals}>×</button>
      </div>
      <div class="modal-body">
        <div class="help-card">
          <div class="help-card-header">
            <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
              <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"/>
            </svg>
            <span>Search Queries</span>
          </div>
          <p class="help-card-sub">Type in the top search bar — click to copy:</p>
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
              <button class="query-chip" on:click={() => copyQuery(q)} class:copied={copied === q}>
                <code>{q}</code>
                {#if copied === q}<span class="copied-label">Copied!</span>{/if}
              </button>
            {/each}
          </div>
          <p class="help-card-note">Supports book names, aliases, chapters, verses, ranges, and cross-book ranges.</p>
        </div>

        <div class="help-card">
          <div class="help-card-header">
            <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
              <path d="M4.715 6.542 3.343 7.914a3 3 0 1 0 4.243 4.243l1.828-1.829A3 3 0 0 0 8.586 5.5L8 6.086a1.002 1.002 0 0 0-.154.199 2 2 0 0 1 .861 3.337L6.88 11.45a2 2 0 1 1-2.83-2.83l.793-.792a4.018 4.018 0 0 1-.128-1.287z"/>
              <path d="M6.586 4.672A3 3 0 0 0 7.414 9.5l.775-.776a2 2 0 0 1-.896-3.346L9.12 3.55a2 2 0 1 1 2.83 2.83l-.793.792c.112.42.155.855.128 1.287l1.372-1.372a3 3 0 1 0-4.243-4.243L6.586 4.672z"/>
            </svg>
            <span>Linked Passages</span>
          </div>
          <p>Use the <strong>"Link passage"</strong> button in the Notes window to associate a Bible reference with a note.</p>
          <p>When viewing that passage in Scripture, linked notes appear in a bar at the bottom. Click to jump to the note.</p>
        </div>
      </div>
    </div>
  </button>
{/if}

<style>
  .titlebar {
    height: 38px;
    background: var(--bg-surface);
    flex-shrink: 0;
  }

  .navbar {
    height: var(--nav-height);
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 16px;
    border-bottom: 1px solid var(--border);
    background: var(--bg);
    flex-shrink: 0;
  }

  .left {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .center {
    flex: 1;
    display: flex;
    justify-content: center;
  }

  .right {
    display: flex;
    align-items: center;
    gap: 4px;
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
    border: none;
    cursor: default;
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
    width: 440px;
  }

  .help-modal .modal-header {
    background: var(--bg-hover);
    border-bottom: 1px solid var(--border);
    border-radius: var(--radius-sm, 6px) var(--radius-sm, 6px) 0 0;
    padding: 16px 20px;
  }

  .help-card {
    background: var(--bg-hover);
    border-radius: 8px;
    padding: 14px 16px;
    margin-bottom: 12px;
    border-left: 3px solid var(--accent);
  }

  .help-card:last-child {
    margin-bottom: 0;
  }

  .help-card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    font-size: 13px;
    margin-bottom: 10px;
    color: var(--text);
  }

  .help-card-header svg {
    color: var(--accent);
    flex-shrink: 0;
  }

  .help-card p {
    font-size: 12px;
    color: var(--text-muted);
    margin: 0 0 8px;
    line-height: 1.5;
  }

  .help-card p:last-child {
    margin-bottom: 0;
  }

  .help-card p strong {
    color: var(--text);
    font-weight: 500;
  }

  .help-card-sub {
    font-size: 11px !important;
    opacity: 0.8;
  }

  .help-card-note {
    font-size: 11px !important;
    opacity: 0.7;
    margin-top: 10px !important;
    margin-bottom: 0 !important;
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

  .query-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin: 8px 0;
  }

  .query-chip {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 5px 10px;
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
  }

  .query-chip:hover {
    border-color: var(--accent);
    background: color-mix(in srgb, var(--accent) 8%, var(--bg-surface));
  }

  .query-chip.copied {
    border-color: var(--accent);
    background: color-mix(in srgb, var(--accent) 15%, var(--bg-surface));
  }

  .query-chip code {
    font-family: var(--font-mono, monospace);
    font-size: 11px;
    color: var(--text);
  }

  .copied-label {
    font-size: 10px;
    color: var(--accent);
    font-weight: 500;
  }

  h4 {
    margin: 16px 0 8px;
    font-size: 13px;
    color: var(--text-muted);
  }

  .info-modal {
    width: 420px;
  }

  .info-modal .modal-header {
    background: var(--bg-hover);
    border-bottom: 1px solid var(--border);
    border-radius: var(--radius-sm, 6px) var(--radius-sm, 6px) 0 0;
    padding: 16px 20px;
  }


  .info-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: var(--bg-hover);
    border-radius: 6px;
  }

  .info-label {
    font-size: 12px;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .info-value {
    font-size: 13px;
    font-weight: 500;
  }

  .info-link {
    color: var(--accent);
    text-decoration: none;
  }

  .info-link:hover {
    text-decoration: underline;
  }

  .info-card {
    background: var(--bg-hover);
    border-radius: 8px;
    padding: 14px 16px;
    margin-bottom: 12px;
    border-left: 3px solid var(--accent);
  }

  .info-card-header {
    font-weight: 600;
    font-size: 13px;
    margin-bottom: 8px;
    color: var(--text);
  }

  .info-card p {
    font-size: 12px;
    color: var(--text-muted);
    margin: 0 0 10px;
    line-height: 1.5;
  }

  .info-card p em {
    color: var(--text);
    font-style: normal;
    font-weight: 500;
  }

  .env-details {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin: 10px 0;
  }

  .env-item {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 12px;
    color: var(--text-muted);
  }

  .env-tag {
    background: var(--accent);
    color: white;
    font-size: 10px;
    font-weight: 600;
    padding: 2px 8px;
    border-radius: 4px;
    letter-spacing: 0.5px;
  }

  .info-note {
    font-size: 11px !important;
    opacity: 0.8;
    margin-bottom: 0 !important;
  }

  .source-links {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .source-btn {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding: 6px 10px;
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: 6px;
    font-size: 12px;
    color: var(--text);
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
  }

  .source-btn:hover {
    border-color: var(--accent);
    background: color-mix(in srgb, var(--accent) 8%, var(--bg-surface));
  }

  .source-btn.copied {
    border-color: var(--accent);
  }

  .copy-hint {
    font-size: 10px;
    color: var(--text-muted);
    opacity: 0;
    transition: opacity 0.15s;
  }

  .source-btn:hover .copy-hint,
  .source-btn.copied .copy-hint {
    opacity: 1;
  }

  .source-btn.copied .copy-hint {
    color: var(--accent);
  }
</style>
