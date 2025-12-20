// --------------------------------------------------------------------------------
// File:        App.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Main application component with sidebar, header, and main content.
// --------------------------------------------------------------------------------
<template>
  <div id="app" class="font-common">
    <div class="app-container">
      <!-- Sidebar -->
      <Sidebar :isHidden="sidebarHidden"/>

      <!-- Main Content Area -->
      <div class="main-content-area">
        <!-- Header -->
        <Header @toggle-sidebar="sidebarHidden = !sidebarHidden"/>

        <!-- Main Content -->
        <MainContent />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from './components/sidebar.vue'
import Header from './components/dashboard/header.vue'
import MainContent from './components/main.vue'

const sidebarHidden = ref(false)
</script>

<style>
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    user-select: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
  }

  input, select, textarea, code, pre {
    user-select: text;
    -webkit-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
  }

  html, body {
    height: 100%;
    font-family: Arial, Helvetica, 'Helvetica Neue', sans-serif;
    color: #1e293b;
    background-color: var(--background-primary);
    line-height: 1.5;
  }

  #app {
    height: 100vh;
    display: flex;
    overflow: hidden;
    width: 100%;
    margin: 0;
    padding: 0;
  }

  .app-container {
    display: flex;
    width: 100%;
    height: 100%;
    overflow: hidden;
    flex: 1;
  }

  .main-content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  width: calc(100% - 250px);
  transition: width 0.3s ease;
  position: relative;
  z-index: 1;
}

  .sidebar-hidden ~ .main-content-area {
    width: 100%;
  }
</style>
