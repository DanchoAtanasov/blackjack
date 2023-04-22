
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import {showLogin, showSignup} from '../stores'

  const dispatch = createEventDispatcher();

  var username, password;

  function login() {
    dispatch('login', {username, password});
  }

  function goToSignUp() {
    showLogin.set(false);
    showSignup.set(true);
  }

</script>

<div class="parent">
  <div class="login-form">
  <h1>Welcome to Blackjack</h1>
    <h2> Log in:</h2>
    <form class="content">
      <input type="text" bind:value={username} placeholder="Username" />
      <input type="password" bind:value={password} placeholder="Password"/>
    </form>
    <button on:click={login}>
      Login
    </button>
  </div>

  <div class="go-to-signup">
    <button on:click={goToSignUp}>
      Go to Signup
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
  