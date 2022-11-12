<script>
  // TODO merge with PlayerHand
  import { dealerHandStore } from '../stores'

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

{#if $dealerHandStore !== undefined}
  <p>Dealer's hand:</p>
  {#each $dealerHandStore.cards as dealerCard}
    <img alt="card" class="playing-card" src={getCardAsset(dealerCard.ValueStr, dealerCard.Suit)} />
  {/each}
  <p>Sum: {$dealerHandStore.sum}</p>
{/if}

<style>
  .playing-card {
    width: 75px;
    height: 150px;
  }
</style>
