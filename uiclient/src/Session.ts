import { name, buyin, dealerCard, dealerSuit } from './stores'
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
    var blackjackHost;
    await fetch(apiServerUrl, {
      method: "POST",
      headers: {'Content-Type': 'application/json'}, 
      body: JSON.stringify(data),
    }).then(res => res.json()
    ).then(resData => {
      console.log(resData);
      token = resData.Token;
      blackjackHost = resData.GameServer;
      console.log(token);
    });

    // // Create WebSocket connection.
    const socket = new WebSocket(`ws://${blackjackHost}`);

    let randomId = uuidv4();
    // Connection opened
    socket.addEventListener('open', (event) => {
        socket.send(`{"Token": "${token}"}`);
    });

    // Listen for messages
    socket.addEventListener('message', (event) => {
        // TODO: improve code quality here
        console.log('Message from server ', event.data);
        var message = JSON.parse(event.data);
        if (message.type === "Game") {
          if (message.message === "Start") {
            console.log("Game started");
          } else if (message.message === "Over") {
            console.log("Game over");
          } else {
            console.log("Wrong game msg");
          }
          return
        } else if (message.type === "DealerHand") {
          console.log("Dealer hand");
          console.log(message.message);
          var dealerHand = JSON.parse(message.message).cards[0];
          dealerCard.set(dealerHand.ValueStr);
          dealerSuit.set(dealerHand.Suit);
          return
        } else {
          console.log(`Got weird event ${event.data}`);
        }
    });
}

name.subscribe(newName => console.log("name change"))