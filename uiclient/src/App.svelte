<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import { name, buyin, dealerCard, dealerSuit } from './stores'

  import Session from './Session';

  function handleStartGame(event) {
    var session = new Session();
    session.connect();
  }
</script>

<main>
  <h1>Welcome to Blackjack</h1>

  <div class="card">
    <PlayButton on:start-game={handleStartGame}/>
  </div>

  <p>Name is {$name}, buy in: {$buyin}</p>
  <p>Dealer's hand {$dealerCard}, {$dealerSuit}</p>

</main>

<style>
</style>
