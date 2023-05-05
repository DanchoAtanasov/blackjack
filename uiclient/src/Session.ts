import { 
  dealerHandStore, playersStore, isConnected, hasGameStarted, NewPlayerRequest,
  newPlayerRequestStore, isGameOver, currBetStore, currTurn, LoginRequest, isLoggedIn, currPlayerName, showLogin, showSignup
} from './stores'
import { get } from 'svelte/store'
import type { Player, Hand } from './stores';


const API_SERVER_URL = "https://blackjack.gg/api/play"
const API_SERVER_COOKIE = "https://blackjack.gg/api/cookie"
const API_SERVER_LOGIN = "https://blackjack.gg/api/login"
const API_SERVER_SIGNUP = "https://blackjack.gg/api/signup"

type Token = {
  Token: string,
}

type GameDetails = {
  Token: string,
  GameServer: string,
}

type CookieLoginResponse = {
  Userid: string,
  Username: string,
}

type Message = {
  type: string,
  message: string,
}


// TODO: session currently handles messaging, factor that out and separate message formats
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
    this.socket = new WebSocket(`wss://${gameDetails.GameServer}`);

    // Send token when connection opened
    this.socket.addEventListener('open', (event) => {
      var token: Token = {"Token": gameDetails.Token};
      this.socket.send(JSON.stringify(token));
    });

    this.socket.addEventListener('error', (event) => {
      console.log(event);
    });

    this.socket.addEventListener('close', (event) => {
      console.log("Connection closed by server.")
      console.log(event);
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

  sendSplit() {
    console.log("Sending split");
    var splitMessage = {"type": "PlayerAction", "message": "Split"}
    this.socket.send(JSON.stringify(splitMessage));
  }

  sendLeave() {
    console.log("Sending leave");
    var leaveMessage = {"type": "PlayerAction", "message": "Leave"}
    this.socket.send(JSON.stringify(leaveMessage));
    hasGameStarted.set(false);
    isGameOver.set(true);
    // Doesn't seem to be closing the connection
    this.socket.close();
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
      if (player.Hands[0].cards === null) {
        player.Hands[0].cards = []; 
      }
      playersStore.set(player.Name, player);
    });
  }

  handlePlayerHandMessages(message: Message){
    var player: Player = JSON.parse(message.message);
    console.log(player);
    currTurn.set(player.Name);
    playersStore.set(player.Name, player);
  }

  handleDealerHandMessages(message: Message){
    var dealerHand: Hand = JSON.parse(message.message);
    currTurn.set('Dealer');
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
        // Doesn't seem to be closing the connection
        this.socket.close();
        break;
      case "IN":
        console.log("GAME IN: asking if playing the hand");
        this.sendIn();
        hasGameStarted.set(true);
        break;
      default:
        console.log("Game message not recognized");
        break;
    }
  }

  sendIn() {
    console.log("Sending in");
    console.log(get(currBetStore));
    
    var message = JSON.stringify({"Playing": true, "CurrentBet": get(currBetStore)});
    var inMessage = {"type": "PlayerAction", "message": message}
    this.socket.send(JSON.stringify(inMessage));
  }

  async getGameDetails(): Promise<GameDetails> {
    console.log("Send player data to api server")
    const newPlayerRequest: NewPlayerRequest = get(newPlayerRequestStore);
    console.log(newPlayerRequest);
    

    var gameDetails: GameDetails = await fetch(API_SERVER_URL, {
      method: "POST",
      headers: {'Content-Type': 'application/json'}, 
      body: JSON.stringify(newPlayerRequest),
      credentials: "include",
    }).then(res => res.json()
    ).catch(err => console.log(`Error getting details ${err}`)
    );
    return gameDetails
  }

  async cookieLogin() {
    console.log("Found cookie, sending it to /cookie")
    var resp: CookieLoginResponse = await fetch(API_SERVER_COOKIE, {
      method: "POST",
      headers: {'Content-Type': 'application/json'}, 
      credentials: "include",
    }).then(resp => {
      if (resp.status === 200) {
        console.log("Login Successful");
        isLoggedIn.set(true);
        showLogin.set(false);
      } else if (resp.status === 401) {
        console.log("Wrong password");
        throw new Error("Cookie login error")
      } else {
        console.log(resp);
        throw new Error("Cookie login error")
      }
      return resp
    }).then(resp => resp.json()).then(respJson => {
      console.log(respJson);
      currPlayerName.set(respJson.Username);
      return respJson
    }).catch(err => console.log(`Error getting details ${err}`)
    );
    return
  }

  async login(loginRequest: LoginRequest): Promise<GameDetails> {
    var resp: void | Response = await fetch(API_SERVER_LOGIN, {
      method: "POST",
      headers: {'Content-Type': 'application/json'}, 
      body: JSON.stringify(loginRequest),
      credentials: "include",
    }).then(resp => {
      if (resp.status === 200) {
        console.log("Login Successful");
        currPlayerName.set(loginRequest.username);
        isLoggedIn.set(true);
        showLogin.set(false);
      } else if (resp.status === 401) {
        console.log("Wrong password");
      } else {
        console.log(resp);
      }
    }).catch(err => console.log(`Error getting details ${err}`)
    );
    return
  }

  async signup(loginRequest: LoginRequest): Promise<GameDetails> {
    var resp: void | Response = await fetch(API_SERVER_SIGNUP, {
      method: "POST",
      headers: {'Content-Type': 'application/json'}, 
      body: JSON.stringify(loginRequest),
      credentials: "include",
    }).then(resp => {
      if (resp.status === 200) {
        console.log("Signup Successful");
        showLogin.set(true);
        showSignup.set(false);
      } else {
        console.log(resp);
      }
    }).catch(err => console.log(`Error getting details ${err}`)
    );
    return
  }
}
