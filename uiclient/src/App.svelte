<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import HitButton from './lib/HitButton.svelte';
  import StandButton from './lib/StandButton.svelte';
  import { name, buyin, dealerHandStore, playerHandStore, playersStore } from './stores'

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
  <p>All players:</p>
  <div>
    {#each $playersStore as player}
      {player}
    {/each}
  </div>

  <p>Name is {$name}, buy in: {$buyin}</p>
  <p>Dealer's hand:</p>
  <div>
    {#each $dealerHandStore.cards as dealerCard}
      <p class="inline-block">{dealerCard.ValueStr} {dealerCard.Suit} | </p> <p></p>
    {/each}
    <p>Sum: {$dealerHandStore.sum}</p>
  </div>

  <p>Player's hand:</p>
  <div>
    {#each $playerHandStore.cards as playerCard}
      <p class="inline-block">{playerCard.ValueStr} {playerCard.Suit} | </p>
    {/each}
    <p>Sum: {$playerHandStore.sum}</p>
  </div>

</main>

<style>
  .inline-block {
    display: inline-block;
    /* color: #003806 */
  }

</style>
