
<script>
  import { v4 as uuidv4 } from 'uuid';
  const play = () => {
    console.log("Button clicked")
    console.log("Send details to api server")
    const apiServerUrl = "http://localhost:3333/play"
    const data = {
      "Name": "UIPlayer",
      "BuyIn": 69,
    };
    var token;
    fetch(apiServerUrl, {
      method: "POST",
      headers: {'Content-Type': 'application/json'}, 
      body: JSON.stringify(data),
    }).then(res => res.json()
    ).then(resData => {
      console.log(resData);
      token = resData.Token;
      console.log(token);
    });

    // // Create WebSocket connection.
    // const socket = new WebSocket('ws://localhost:8080');

    // let randomId = uuidv4(); // ⇨ '9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d'
    // // Connection opened
    // socket.addEventListener('open', (event) => {
    //     socket.send(`{"Token": "${randomId}"}`);
    // });

    // // Listen for messages
    // socket.addEventListener('message', (event) => {
    //     console.log('Message from server ', event.data);
    // });
  }
</script>

<button on:click={play}>
  Play
</button>
  