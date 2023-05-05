
<script>
  import HitButton from '../lib/HitButton.svelte';
  import StandButton from '../lib/StandButton.svelte';
  import PlayerHand from '../lib/PlayerHand.svelte';
  import { currPlayerName, playersStore } from '../stores';

  import DealerHand from '../lib/DealerHand.svelte';
  import CurrentBet from '../lib/CurrentBet.svelte';
  import LeaveButton from './LeaveButton.svelte';
    import SplitButton from './SplitButton.svelte';
</script>

<h3>Game Table</h3>

<p>Name is {$currPlayerName}, 
    <!--TODO player data in store isn't quite ready by the time the game starts
    reorder messages and remove the if check -->
    {#if $playersStore.get($currPlayerName)}
    buy in: {$playersStore.get($currPlayerName).BuyIn}
    {/if}
</p>
<DealerHand></DealerHand>
<PlayerHand></PlayerHand>


<HitButton on:hit/>
<StandButton on:stand/>
<SplitButton on:split/>
<CurrentBet/>
<LeaveButton on:leave/>
