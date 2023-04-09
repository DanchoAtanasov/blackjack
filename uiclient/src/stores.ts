import {writable, derived, get} from 'svelte/store';

type Card = {
  ValueStr: string,
  Suit: string,
};

export type Hand = {
  cards: Card[],
  sum: number,
}

export type Player = {
  Name: string,
	BuyIn: number,
  Hand: Hand,
	CurrentBet: number,
}

export type NewPlayerRequest = {
  BuyIn: number,
  CurrBet: number,
}

export type LoginRequest = {
  username: string,
  password: string,
}

export const currPlayerName = writable('');
export const currTurn = writable('');
export const buyin = writable(0);
export const newPlayerRequestStore = writable<NewPlayerRequest>();
export const isConnected = writable(false);
export const isLoggedIn = writable(false);
export const showLogin = writable(false);
export const showSignup = writable(false);
export const hasGameStarted = writable(false);
export const isGameOver = writable(false);
export const currBetStore = writable(0);

export const dealerHandStore = writable<Hand>();

type Players = { [name: string]: Player };

// Got from https://svelte.dev/repl/ccbc94cb1b4c493a9cf8f117badaeb31?version=3.53.1
// [name: string]: Player
function createMapStore(initial) {
  const store = writable<Players>(initial);
  const set = (key: string, value: Player) => store.update(m => Object.assign({}, m, {[key]: value}));
  const clear = () => {
    store.update(players => {
      Object.keys(players).forEach(playerName => {
        delete players[playerName];
      });
      return players;
    })
  }
  const results = derived(store, s => ({
    keys: Object.keys(s),
    values: Object.values(s),
    entries: Object.entries(s),
    set(k: string, v: Player) {
      store.update(s => Object.assign({}, s, {[k]: v}))
    },
    get(k: string) {
      var p = s[k];
      return p;
    },
    remove(k: string) {
      store.update(s => {
        delete s[k];
        return s;
      });
    },
  }));
  return {
    subscribe: results.subscribe,
    set: set,
    // This is actually svelte get
    get: get,
    clear: clear,
  }
}

export const playersStore = createMapStore({});
