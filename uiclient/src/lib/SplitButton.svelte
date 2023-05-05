
<script>
  import { createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import {currTurn, currPlayerName} from '../stores'

  const dispatch = createEventDispatcher();

  function sendSplit() {
    dispatch('split', {});
  }

  var isDisabled = true;
  currTurn.subscribe((currTurnPlayerName) => {
    isDisabled = currTurnPlayerName !== get(currPlayerName); 
  })
  
</script>


<button disabled={isDisabled} on:click={sendSplit}>
  Split
</button>
  