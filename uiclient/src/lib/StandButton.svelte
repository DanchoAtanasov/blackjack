
<script>
  import { createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import {currTurn, currPlayerName} from '../stores'

  const dispatch = createEventDispatcher();

  function sendStand() {
    dispatch('stand', {});
  }

  var isDisabled = true;
  currTurn.subscribe((currTurnPlayerName) => {
    isDisabled = currTurnPlayerName !== get(currPlayerName); 
  })
  
</script>


<button disabled={isDisabled} on:click={sendStand}>
  Stand
</button>
  