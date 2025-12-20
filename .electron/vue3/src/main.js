// --------------------------------------------------------------------------------
// File:        main.js
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Main entry point for the Vue3 application.
// --------------------------------------------------------------------------------
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './assets/css/font-awesome.min.css'
import './assets/css/components.css'

const app = createApp(App)
app.use(createPinia())
app.mount('#app')
