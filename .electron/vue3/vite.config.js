// --------------------------------------------------------------------------------
// File:        vite.config.js
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Vite configuration file for the Vue3 project.
// --------------------------------------------------------------------------------
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    open: true
  }
})
