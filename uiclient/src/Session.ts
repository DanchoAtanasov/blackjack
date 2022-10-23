import { name, buyin } from './stores'
import { get } from 'svelte/store'
import { v4 as uuidv4 } from 'uuid';

export async function startSession() {
    var currName = get(name);
    var currBuyIn = get(buyin);
    console.log(`Starting session ${currName},  ${currBuyIn}`);

    console.log("Send details to api server")
    const apiServerUrl = "http://localhost:3333/play"
    const data = {
      "Name": currName,
      "BuyIn": Number(currBuyIn),
    };
    console.log(data);
    var token;
    await fetch(apiServerUrl, {
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
    const socket = new WebSocket('ws://localhost:8080');

    let randomId = uuidv4();
    // Connection opened
    socket.addEventListener('open', (event) => {
        socket.send(`{"Token": "${token}"}`);
    });

    // Listen for messages
    socket.addEventListener('message', (event) => {
        console.log('Message from server ', event.data);
    });
}

name.subscribe(newName => console.log("name change"))