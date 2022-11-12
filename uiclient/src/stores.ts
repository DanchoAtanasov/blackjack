import { writable } from 'svelte/store';

type Card = {
  ValueStr: string,
  Suit: string,
};

export type Hand = {
  cards: Card[],
  sum: number,
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

export const playersStore = writable<string[]>([]);
