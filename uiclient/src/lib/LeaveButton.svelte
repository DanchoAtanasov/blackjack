
<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function sendLeave() {
    dispatch('leave', {});
  }
</script>


<button on:click={sendLeave}>
  Leave
</button>
  