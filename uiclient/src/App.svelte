<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import HitButton from './lib/HitButton.svelte';
  import StandButton from './lib/StandButton.svelte';
  import { name, buyin, dealerHandStore, playerHandStore } from './stores'

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

  function sendStand() {
    session.sendStand();
  }

</script>

<main>
  <h1>Welcome to Blackjack</h1>

  <div class="card">
    <PlayButton on:start-game={handleStartGame}/>
  </div>
  {#if active}
    <HitButton on:hit={sendHit}/>
    <StandButton on:stand={sendStand}/>
  {/if}

  <p>Name is {$name}, buy in: {$buyin}</p>
  <p>Dealer's hand:</p>
  {#each $dealerHandStore as dealerCard}
    <p>{dealerCard.ValueStr} {dealerCard.Suit}</p>
	{/each}

  <p>Player's hand:</p>
  {#each $playerHandStore as playerCard}
    <p> {playerCard.ValueStr} {playerCard.Suit}</p>
	{/each}

</main>

<style>
</style>
