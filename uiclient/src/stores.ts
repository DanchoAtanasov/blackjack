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

export const name = writable('');
export const buyin = writable(0);

export const dealerHandStore = writable<Hand>({
    cards: [{
        ValueStr: "",
        Suit: "",
    }],
    sum: 0,
});
export const playerHandStore = writable<Hand>({
    cards: [{
        ValueStr: "",
        Suit: "",
    }],
    sum: 0,
});

type Players = { [name: string]: Player };

// Got from https://svelte.dev/repl/ccbc94cb1b4c493a9cf8f117badaeb31?version=3.53.1
// [name: string]: Player
function createMapStore(initial) {
  const store = writable<Players>(initial);
  const set = (key: string, value: Player) => store.update(m => Object.assign({}, m, {[key]: value}));
  const results = derived(store, s => ({
    keys: Object.keys(s),
    values: Object.values(s),
    entries: Object.entries(s),
    set(k: string, v: Player) {
      store.update(s => Object.assign({}, s, {[k]: v}))
    },
    remove(k: string) {
      store.update(s => {
        delete s[k];
        return s;
      });
    }
  }));
  return {
    subscribe: results.subscribe,
    set: set,
  }
}

export const playersStore = createMapStore({});
