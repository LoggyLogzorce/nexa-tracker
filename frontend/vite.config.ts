import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    allowedHosts: ['tasks.kivoapps.ru'],
    proxy: {
      '/api': {
        target: 'https://tasks.kivoapps.ru:8443',

      '/uploads': {
        target: 'https://tasks.kivoapps.ru:8443',
        changeOrigin: true,
      },
    },
  },
})
