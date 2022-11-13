
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { currPlayerName, newPlayerRequestStore } from '../stores';

  const dispatch = createEventDispatcher();

  // var nameInput: string, buyinInput: string;
  var nameInput, buyinInput;

  function play() {
    currPlayerName.set(nameInput);

    newPlayerRequestStore.set({
      Name: nameInput,
      BuyIn: Number(buyinInput),
    });

    dispatch('start-game', {});
  }

</script>


<form class="content">
  <input type="text" bind:value={nameInput} placeholder="Name" />
  <input type="text" bind:value={buyinInput} placeholder="Buy In"/>
</form>
<button on:click={play}>
  Play
</button>
  