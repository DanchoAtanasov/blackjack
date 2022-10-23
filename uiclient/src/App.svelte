<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import { name, buyin } from './stores'

  import { startSession } from './Session';

  function handleStartGame(event) {
    console.log(event);
    startSession(event.detail.name, event.detail.buyin);
  }
</script>

<main>
  <h1>Welcome to Blackjack</h1>

  <div class="card">
    <PlayButton on:start-game={handleStartGame}/>
  </div>

  <p>Name is {$name}, buy in: {$buyin}</p>

</main>

<style>
</style>
