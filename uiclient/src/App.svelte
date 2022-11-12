<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import HitButton from './lib/HitButton.svelte';
  import StandButton from './lib/StandButton.svelte';
  import PlayerHand from './lib/PlayerHand.svelte';
  import { name, buyin, dealerHandStore, playersStore } from './stores'

  import Session from './Session';
  import DealerHand from './lib/DealerHand.svelte';

  var session = new Session();

  var active = false;

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

  {#if !active}
    <PlayButton on:start-game={handleStartGame}/>

  {:else}
    <p>Name is {$name}, buy in: {$buyin}</p>
    <DealerHand></DealerHand>
    <PlayerHand></PlayerHand>


    <HitButton on:hit={sendHit}/>
    <StandButton on:stand={sendStand}/>

  {/if}
</main>

<style>
  /* color: #003806 */
</style>
