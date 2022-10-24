<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import HitButton from './lib/HitButton.svelte';
  import { name, buyin, dealerCard, dealerSuit, playerCard, playerSuit } from './stores'

  import Session from './Session';

  var session = new Session();

  var active = false

  function handleStartGame(event) {
    active = true;
    session.connect();
  }

  function sendHit() {
    session.sendHit();
  }

</script>

<main>
  <h1>Welcome to Blackjack</h1>

  <div class="card">
    <PlayButton on:start-game={handleStartGame}/>
  </div>
  {#if active}
    <HitButton on:hit={sendHit}/>
  {/if}

  <p>Name is {$name}, buy in: {$buyin}</p>
  <p>Dealer's hand {$dealerCard}, {$dealerSuit}</p>
  <p>Player's hand {$playerCard}, {$playerSuit}</p>

</main>

<style>
</style>
