import { writable } from 'svelte/store';

type Hand = {
    ValueStr: string,
    Suit: string,
}[];

export const name = writable('');
export const buyin = writable(0);
export const dealerHandStore = writable<Hand>([{ValueStr: "", Suit: ""}]);
export const playerHandStore = writable<Hand>([{ValueStr: "", Suit: ""}]);
