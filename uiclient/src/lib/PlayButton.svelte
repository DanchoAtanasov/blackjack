
<script>
  import { createEventDispatcher } from 'svelte';
  import { name, buyin } from '../stores';

  const dispatch = createEventDispatcher();

  var nameInput, buyinInput;

  function startGame() {
    dispatch('start-game', {
      name: name,
      buyin: buyin,
    });
  }

  function putInStore() {
    name.set(nameInput);
    buyin.set(buyinInput);
    startGame();
  }

  const play = async () => {
  }
</script>


<form class="content">
  <input type="text" bind:value={nameInput} placeholder="Name" />
  <input type="text" bind:value={buyinInput} placeholder="Buy In"/>
</form>
<button on:click={putInStore}>
  Play
</button>
  