
<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function sendHit() {
    dispatch('hit', {});
  }
</script>


<button on:click={sendHit}>
  Hit
</button>
  