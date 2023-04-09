
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { showLogin, showSignup } from '../stores'

  const dispatch = createEventDispatcher();

  function handleShowLogin() {
    showLogin.set(true);
    showSignup.set(false);
  }

  function handleShowSignup() {
    showSignup.set(true);
    showLogin.set(false);
  }

</script>


<button on:click={handleShowLogin}>
  Login
</button>
<button on:click={handleShowSignup}>
  Signup
</button>

<style></style>
  