
<script>
  import { createEventDispatcher } from 'svelte';
  import { get } from 'svelte/store';
  import {currTurn, currPlayerName, playDealSound} from '../stores'

  const dispatch = createEventDispatcher();

  function sendHit() {
    dispatch('hit', {});
    playDealSound.update((value) => !value);
  }

  var isDisabled = true;
  currTurn.subscribe((currTurnPlayerName) => {
    isDisabled = currTurnPlayerName !== get(currPlayerName); 
  })
  
</script>


<button disabled={isDisabled} on:click={sendHit}>
  Hit
</button>
  