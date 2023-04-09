<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import HitButton from './lib/HitButton.svelte';
  import StandButton from './lib/StandButton.svelte';
  import PlayerHand from './lib/PlayerHand.svelte';
  import { 
    currPlayerName, playersStore, isConnected, isLoggedIn, hasGameStarted, isGameOver, showLogin,
    showSignup,
  } from './stores';

  import { onMount } from 'svelte';
  import Session from './Session';
  import DealerHand from './lib/DealerHand.svelte';
  import CurrentBet from './lib/CurrentBet.svelte';
  import LoginButton from './lib/LoginButton.svelte';
  import LandingPage from './lib/LandingPage.svelte';
  import SignupButton from './lib/SignupButton.svelte';

  var session = new Session();

  function checkTokenCookie() {
    console.log(document.cookie);
  }

  onMount(async () => {
    checkTokenCookie();
    session.callNew();
  });


  function handleStartGame(event) {
    session.connect();
  }

  function handleLogin(event) {
    session.login(event.detail);
  }

  function handleSignUp(event) {
    session.signup(event.detail);
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

  {#if $showLogin}
    <LoginButton on:login={handleLogin}/>
  {:else if $showSignup}
    <SignupButton on:signup={handleSignUp}/>
  {:else if !$isLoggedIn}
    <LandingPage on:login={handleLogin}/>
  {:else if !$isConnected}
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
