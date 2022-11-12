
<script>
  import { playersStore } from '../stores'

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

{#each $playersStore.values as player}
  <p>{player.Name}'s hand:</p>
  <div>
    {#each player.Hand.cards as playerCard}
      <img alt="card" class="playing-card" src={getCardAsset(playerCard.ValueStr, playerCard.Suit)} />
    {/each}
    <p>Sum: {player.Hand.sum}</p>
  </div>
{/each}

<style>
  .playing-card {
    width: 75px;
    height: 150px;
  }
</style>
  