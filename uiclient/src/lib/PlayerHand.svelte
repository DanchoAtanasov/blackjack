
<script>
  import SplitButton from "./SplitButton.svelte";


  export let player;

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

    var assetPath = `../assets/svg-cards/${value}_of_${suit.toLowerCase()}.svg`;
    const cardAsset = new URL(assetPath, import.meta.url).href

    return cardAsset;
  }
</script>

<div>
  {#each player.Hands as playerHand}
    <div class="outer">
      {#each playerHand.cards as playerCard}
        <img alt="card" class="playing-card" src={getCardAsset(playerCard.ValueStr, playerCard.Suit)} />
      {/each}
      <p>Sum: {playerHand.sum}</p>
      {#if playerHand.cards.length === 2 && playerHand.cards[0].ValueStr === playerHand.cards[1].ValueStr}
        <SplitButton on:split/>
      {/if}
    </div>
  {/each}
</div>

<style>
  .playing-card {
    width: 75px;
    height: 150px;
  }
  .outer {
    position: relative;
  }
</style>
  