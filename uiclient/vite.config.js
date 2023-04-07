import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import fs from 'fs';

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: 'blackjack.gg',
    https: {
      key: fs.readFileSync('./keys/key.pem'),
      cert: fs.readFileSync('./certs/cert.crt'),
    },
    proxy: {},
  },
  plugins: [svelte()]
})
