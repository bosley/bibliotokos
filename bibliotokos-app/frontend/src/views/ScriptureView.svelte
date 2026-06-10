<script>
  import NavBar from '../components/nav/NavBar.svelte'
  import ScriptureWorkspace from '../components/scripture/ScriptureWorkspace.svelte'
</script>

<div class="scripture-view">
  <NavBar />
  <main class="content">
    <ScriptureWorkspace />
  </main>
</div>

<style>
  .scripture-view {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .content {
    flex: 1;
    overflow: hidden;
    display: flex;
  }
</style>
