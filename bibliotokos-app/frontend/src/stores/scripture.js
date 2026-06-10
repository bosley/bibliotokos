import { writable } from 'svelte/store'

export const query = writable('')
export const versions = writable([])
export const books = writable([])
export const scrollLock = writable(true)
