<script lang="ts">
  import PlayButton from './lib/PlayButton.svelte';
  import HitButton from './lib/HitButton.svelte';
  import StandButton from './lib/StandButton.svelte';
  import { name, buyin, dealerHandStore, playersStore } from './stores'

  import Session from './Session';

  var session = new Session();

  var active = false;

  function handleStartGame(event) {
    active = true;
    session.connect();
  }

  function sendHit() {
    session.sendHit();
  }

  function sendStand() {
    session.sendStand();
  }

  function getCardAsset(value, suit) {
    // TODO fix asset name mismatch
    switch (value) {
      case 'J':
        value = 'jack'
        break;
      case 'Q':
        value = 'queen'
        break;
      case 'K':
        value = 'king'
        break;
      case 'A':
        value = 'ace'
        break;
    }

    var assetPath = `./assets/svg-cards/${value}_of_${suit.toLowerCase()}.svg`;
    const cardAsset = new URL(assetPath, import.meta.url).href

    return cardAsset;
  }

</script>

<main>
  <h1>Welcome to Blackjack</h1>

  {#if !active}
    <PlayButton on:start-game={handleStartGame}/>

  {:else}
    <p>Name is {$name}, buy in: {$buyin}</p>
    <div>
      {#if $dealerHandStore !== undefined}
        <p>Dealer's hand:</p>
        {#each $dealerHandStore.cards as dealerCard}
          <img alt="card" class="playing-card" src={getCardAsset(dealerCard.ValueStr, dealerCard.Suit)} />
        {/each}
        <p>Sum: {$dealerHandStore.sum}</p>
      {/if}
    </div>

    {#each $playersStore.values as player}
      <p>{player.Name}'s hand:</p>
      <div>
        {#each player.Hand.cards as playerCard}
          <img alt="card" class="playing-card" src={getCardAsset(playerCard.ValueStr, playerCard.Suit)} />
        {/each}
        <p>Sum: {player.Hand.sum}</p>
      </div>
    {/each}

    <HitButton on:hit={sendHit}/>
    <StandButton on:stand={sendStand}/>

  {/if}
</main>

<style>
  .inline-block {
    display: inline-block;
    /* color: #003806 */
  }

  .playing-card {
    width: 75px;
    height: 150px;
  }

</style>
