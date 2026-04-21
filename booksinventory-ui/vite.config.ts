import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  server: {
    // port: 3000, // custom port (default port is 5173)
    open: true, // As soon as you run the command 'npm run dev', your browser will open
    // This is how you can configure HTTPS
    // https:{
    //   key: fileURLToPath(new URL('./ssl/server.key', import.meta.url)),
    //   cert: fileURLToPath(new URL('./ssl/server.crt', import.meta.url))
    // }
  },
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
})
