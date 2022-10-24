
<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function sendStand() {
    dispatch('stand', {});
  }
</script>


<button on:click={sendStand}>
  Stand
</button>
  