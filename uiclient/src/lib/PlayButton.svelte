
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import { currPlayerName, newPlayerRequestStore, currBetStore } from '../stores';

  const dispatch = createEventDispatcher();

  var nameInput, buyinInput;

  function play() {
    currPlayerName.set(nameInput);

    console.log(get(currBetStore));
    
    newPlayerRequestStore.set({
      Name: nameInput,
      BuyIn: Number(buyinInput),
      CurrBet: get(currBetStore),
    });

    dispatch('start-game', {});
  }

</script>


<form class="content">
  <input type="text" bind:value={nameInput} placeholder="Name" />
  <input type="text" bind:value={buyinInput} placeholder="Buy In"/>
  <input type="number" bind:value={$currBetStore} placeholder="Current Bet"/>
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
</style>
  