
<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  var username, password;

  function signup() {
    dispatch('signup', {username, password});
  }

</script>


<form class="content">
  <input type="text" bind:value={username} placeholder="Username" />
  <input type="password" bind:value={password} placeholder="Password"/>
</form>
<button on:click={signup}>
  Signup
</button>

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
</style>
  