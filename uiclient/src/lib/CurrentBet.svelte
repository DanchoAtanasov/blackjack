
<script>
  import { currBetStore } from '../stores'
</script>


<label>
  <h4>Next bet:
  	<input type=number class="input-box" bind:value={$currBetStore} min=1>
  </h4>
</label>
  
<style>
  input::-webkit-outer-spin-button,
  input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }
   .input-box {
     width: 40px;
     text-align: center;
     border-radius: 100px;
     border: 2px solid #73AD21;
     padding: 5px;
   }
</style>