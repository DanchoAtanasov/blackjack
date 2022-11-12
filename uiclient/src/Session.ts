import { name, buyin, dealerHandStore, playerHandStore } from './stores'
import { get } from 'svelte/store'
import type { Hand } from './stores';
import { v4 as uuidv4 } from 'uuid';


const API_SERVER_URL = "http://localhost:3333/play"

type Token = {
  Token: string,
}

type NewPlayerRequest = {
  Name: string,
  BuyIn: number,
}

type GameDetails = {
  Token: string,
  GameServer: string,
}

type Message = {
  type: string,
  message: string,
}


export default class Session {
  private socket: WebSocket;
  private active: boolean;

  constructor() {
    this.socket = undefined;
    this.active = false;
  }

  async connect() {
    console.log("Connecting");
    var gameDetails = await this.getGameDetails();

    // // Create WebSocket connection.
    this.socket = new WebSocket(`ws://${gameDetails.GameServer}`);

    // Send token when connection opened
    this.socket.addEventListener('open', (event) => {
      var token: Token = {"Token": gameDetails.Token};
      this.socket.send(JSON.stringify(token));
    });

    this.addMessageListeners();
    this.active = true;
    console.log(this.socket);
    
  }

  isActive(): boolean {
    console.log(this.active);
    return this.active
  }

  sendHit() {
    console.log("Sending hit");
    var hitMessage = {"type": "PlayerAction", "message": "Hit"}
    this.socket.send(JSON.stringify(hitMessage));
  }

  sendStand() {
    console.log("Sending stand");
    var standMessage = {"type": "PlayerAction", "message": "Stand"}
    this.socket.send(JSON.stringify(standMessage));
  }

  addMessageListeners() {
    // Listen for messages
    this.socket.addEventListener('message', (event) => {
        // TODO: improve code quality here
        console.log('Message from server ', event.data);
        var message = JSON.parse(event.data);
        switch (message.type) {
          case "Game":
            this.handleGameMessages(message);
            break;
          case "DealerHand":
            this.handleDealerHandMessages(message);
            break;
          case "PlayerHand":
            this.handlePlayerHandMessages(message);
            break;
          case "HandState":
            this.handleHandStateMessages(message);
            break;
          default:
            console.log("Message type not recognized");
            break;
        }
    });
  }

  handleHandStateMessages(message: Message){
    switch (message.message) {
      case "Bust":
        console.log("Bust");
        break;
      case "Blackjack":
        console.log("Blackjack");
        break;
      default:
        console.log("Hand state message not recognized");
        break;
    }
  }

  handlePlayerHandMessages(message: Message){
    var playerHand: Hand = JSON.parse(message.message);
    playerHandStore.set(playerHand);
  }

  handleDealerHandMessages(message: Message){
    var dealerHand: Hand = JSON.parse(message.message);
    dealerHandStore.set(dealerHand);
  }

  handleGameMessages(message: Message){
    switch (message.message) {
      case "Start":
        console.log("Game started");
        break;
      case "Over":
        console.log("Game over");
        this.active = false
        break;
      default:
        console.log("Game message not recognized");
        break;
    }
  }

  async getGameDetails(): Promise<GameDetails> {
    console.log("Send player data to api server")
    var newPlayerRequest: NewPlayerRequest = {
      Name: get(name),
      BuyIn: Number(get(buyin)),
    }

    var gameDetails: GameDetails = await fetch(API_SERVER_URL, {
      method: "POST",
      headers: {'Content-Type': 'application/json'}, 
      body: JSON.stringify(newPlayerRequest),
    }).then(res => res.json()
    ).catch(err => console.log(`Error getting details ${err}`)
    );
    return gameDetails
  }
}

name.subscribe(newName => console.log("name change"))