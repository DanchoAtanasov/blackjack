<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import { 
    currPlayerName, playersStore, isConnected, isLoggedIn, hasGameStarted, isGameOver, showLogin,
    showSignup,
  } from './stores';

  import { onMount } from 'svelte';
  import Session from './Session';
  import LoginButton from './lib/LoginButton.svelte';
  import LandingPage from './lib/LandingPage.svelte';
  import SignupButton from './lib/SignupButton.svelte';
  import GameView from './lib/GameView.svelte';
  import AudioPlayer from './lib/AudioPlayer.svelte'

  var session = new Session();

  function checkTokenCookie() {
    console.log(document.cookie);
  }

  onMount(async () => {
    checkTokenCookie();
    if (document.cookie != "") {
      session.cookieLogin();
    } else {
      console.log(document.cookie);
    }
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

  function handleSendHit(event) {
    session.sendHit();
  }

  function handleSendStand(event) {
    session.sendStand();
  }

  function handleSendLeave(event) {
    session.sendLeave();
  }

  function handleSendSplit(event) {
    session.sendSplit();
  }
</script>

<main>
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
    <GameView on:hit={handleSendHit} on:stand={handleSendStand} on:leave={handleSendLeave} on:split={handleSendSplit}></GameView>
  {/if}
  <AudioPlayer></AudioPlayer>
</main>

<style>
</style>
