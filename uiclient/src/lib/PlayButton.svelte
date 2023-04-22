
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import { currPlayerName, newPlayerRequestStore, currBetStore } from '../stores';

  const dispatch = createEventDispatcher();

  function play() {
    newPlayerRequestStore.set({
      CurrBet: get(currBetStore),
    });

    dispatch('start-game', {});
  }

</script>

<h2>Nice to see you again, {$currPlayerName}</h2>
<h2> Fancy a game?</h2>

<form class="content">
  <label for="current-bet">Starting bet:</label>
  <input type="number" id="current-bet" bind:value={$currBetStore}/>
</form>
<button on:click={play}>
  Play
</button>

<style>
  input::-webkit-outer-spin-button,
  input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }
  #current-bet {
    width: 40px;
  }
</style>
  