
<script>
  import { currBetStore } from '../stores'
</script>


<label>
	<input type=number bind:value={$currBetStore} min=1>
</label>
  
<style>
  input::-webkit-outer-spin-button,
  input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }
</style>