<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import HitButton from './lib/HitButton.svelte';
  import StandButton from './lib/StandButton.svelte';
  import PlayerHand from './lib/PlayerHand.svelte';
  import { currPlayerName, playersStore, isConnected, hasGameStarted, isGameOver } from './stores';

  import { onMount } from 'svelte';
  import Session from './Session';
  import DealerHand from './lib/DealerHand.svelte';
  import CurrentBet from './lib/CurrentBet.svelte';

  function checkTokenCookie() {
    console.log(document.cookie);
  }

  onMount(async () => {
    checkTokenCookie();
  });

  var session = new Session();

  function handleStartGame(event) {
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

  {#if !$isConnected}
    <PlayButton on:start-game={handleStartGame}/>
  {:else if !$hasGameStarted && !$isGameOver}
    <p>Name is {$currPlayerName}, waiting for game to begin...</p>
  {:else if $isGameOver}
    <p> Game is over, final buyin: {$playersStore.get($currPlayerName).BuyIn} </p>
  {:else}
    <p>Name is {$currPlayerName}, 
      <!--TODO player data in store isn't quite ready by the time the game starts
      reorder messages and remove the if check -->
      {#if $playersStore.get($currPlayerName)}
      buy in: {$playersStore.get($currPlayerName).BuyIn}
      {/if}
    </p>
    <DealerHand></DealerHand>
    <PlayerHand></PlayerHand>


    <HitButton on:hit={sendHit}/>
    <StandButton on:stand={sendStand}/>
    <CurrentBet></CurrentBet>

  {/if}
</main>

<style>
</style>
