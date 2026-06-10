import { writable } from 'svelte/store'

const theme = writable('dark')

theme.subscribe(t => {
  if (typeof document !== 'undefined') {
    document.body.dataset.theme = t
  }
})

export { theme }
