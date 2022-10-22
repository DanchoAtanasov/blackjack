<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import { name, buyin } from './stores'
</script>

<main>
  <h1>Welcome to Blackjack</h1>

  <div class="card">
    <PlayButton on:start-game/>
  </div>

  <p>Name is {$name}, buy in: {$buyin}</p>

</main>

<style>
</style>
