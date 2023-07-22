
<script>
    import {playDealSound, playWinSound} from '../stores'
  
    const dealCardUrl = '/deal-card.mp3';
    const winUrl = '/instant-win.wav';

    let dealCardAudioElement;
    let winAudioElement;
  
    playDealSound.subscribe((value) => {
      playAudio(dealCardAudioElement);
    })

    playWinSound.subscribe((value) => {
      playAudio(winAudioElement);
    })
  
    function playAudio(audioElement) {
      if (audioElement === undefined) return;
      audioElement.play();
    }
    
</script>
    
<audio id="audio-deal" bind:this={dealCardAudioElement} style="display: none;">
    <source src={dealCardUrl} type="audio/mp3">
</audio>

<audio id="audio-win" bind:this={winAudioElement} style="display: none;">
    <source src={winUrl} type="audio/wav">
</audio>