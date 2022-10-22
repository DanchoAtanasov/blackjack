
<script>
  import { v4 as uuidv4 } from 'uuid';
  const play = () => {
    console.log("Button clicked")
    // Create WebSocket connection.
    const socket = new WebSocket('ws://localhost:8080');

    let randomId = uuidv4(); // ⇨ '9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d'
    // Connection opened
    socket.addEventListener('open', (event) => {
        socket.send(`{"Token": "${randomId}"}`);
    });

    // Listen for messages
    socket.addEventListener('message', (event) => {
        console.log('Message from server ', event.data);
    });
  }
</script>

<button on:click={play}>
  Play
</button>
  