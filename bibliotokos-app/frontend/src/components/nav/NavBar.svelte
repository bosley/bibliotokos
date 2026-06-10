<script>
  import { Events } from '@wailsio/runtime'
  import NavButton from './NavButton.svelte'
  import SearchInput from './SearchInput.svelte'
  import ThemeToggle from './ThemeToggle.svelte'

  function openNotes() {
    Events.Emit('open-notes', null)
  }
</script>

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
  </div>
</nav>

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
</style>
