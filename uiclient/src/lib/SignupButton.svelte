
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import {showLogin, showSignup} from '../stores'

  const dispatch = createEventDispatcher();

  var username, password;

  function signup() {
    dispatch('signup', {username, password});
  }

  function goToLogin() {
    showLogin.set(true);
    showSignup.set(false);
  }

</script>

<div class="parent">
  <div class="signup-form">
    <h1>Welcome to Blackjack</h1>
    <h2> Sign in:</h2>
      <form class="content">
        <input type="text" bind:value={username} placeholder="Username" />
        <input type="password" bind:value={password} placeholder="Password"/>
      </form>
      <button on:click={signup}>
        Signup
      </button>
  </div>

  <div class="go-to-login">
    <button on:click={goToLogin}>
      Go to Login
    </button>
  </div>
</div>

<style>
  input::-webkit-outer-spin-button,
  input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }
  input {
    display: block;
    margin: auto;
    margin-bottom: 6px;
    width: 35%;
    justify-content: center;
  }
  div {
    margin-bottom: 10px;
  }
</style>
  