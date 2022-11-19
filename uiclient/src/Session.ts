import { 
  dealerHandStore, playersStore, isConnected, hasGameStarted, NewPlayerRequest,
  newPlayerRequestStore, isGameOver,
} from './stores'
import { get } from 'svelte/store'
import type { Player, Hand } from './stores';


const API_SERVER_URL = "http://localhost:3333/play"

type Token = {
  Token: string,
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
  // public gameStarted: boolean;
  // public connected: boolean;

  constructor() {
    this.socket = undefined;
    // this.connected = false;
    // this.gameStarted = false;
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

    this.socket.addEventListener('close', (event) => {
      console.log("Connection closed by server.")
      isGameOver.set(true);
    });

    isConnected.set(true);
    console.log("Connected");
    this.addMessageListeners();
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
          case "ListPlayers":
            this.handleListPlayersMessage(message);
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

  handleListPlayersMessage(message: Message){
    var players: Player[] = JSON.parse(message.message);
    // Remove disconnected players
    playersStore.clear();
    players.forEach(player => {
      if (player.Hand.cards === null) {
        player.Hand.cards = []; 
      }
      playersStore.set(player.Name, player);
    });
  }

  handlePlayerHandMessages(message: Message){
    var player: Player = JSON.parse(message.message);
    console.log(player);
    playersStore.set(player.Name, player);
  }

  handleDealerHandMessages(message: Message){
    var dealerHand: Hand = JSON.parse(message.message);
    dealerHandStore.set(dealerHand);
  }

  handleGameMessages(message: Message){
    switch (message.message) {
      case "Start":
        console.log("Game started");
        hasGameStarted.set(true);
        break;
      case "Over":
        console.log("Game over");
        hasGameStarted.set(false);
        isGameOver.set(true);
        this.socket.close();
        break;
      default:
        console.log("Game message not recognized");
        break;
    }
  }

  async getGameDetails(): Promise<GameDetails> {
    console.log("Send player data to api server")
    const newPlayerRequest: NewPlayerRequest = get(newPlayerRequestStore);

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
