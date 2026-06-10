<script>
  import { onMount } from 'svelte'
  import { theme } from './stores/theme.js'
  import { versions, books } from './stores/scripture.js'
  import { GetTheme } from '../bindings/bibliotokos/services/system/systemservice.js'
  import { GetVersions, GetBooks } from '../bindings/bibliotokos/services/bible/bibleservice.js'
  import ScriptureView from './views/ScriptureView.svelte'
  import NotesView from './views/NotesView.svelte'

  const view = new URLSearchParams(window.location.search).get('view') ?? 'scripture'

  onMount(async () => {
    const [t, v, b] = await Promise.all([
      GetTheme().catch(() => 'dark'),
      GetVersions().catch(() => []),
      GetBooks().catch(() => []),
    ])
    theme.set(t)
    versions.set(v)
    books.set(b)
  })
</script>

{#if view === 'notes'}
  <NotesView />
{:else}
  <ScriptureView />
{/if}
