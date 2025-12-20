// --------------------------------------------------------------------------------
// File:        main.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Main content component for the application.
// --------------------------------------------------------------------------------
<template>
  <main class="main-content">
    <div class="content-placeholder">
      <h2>Welcome</h2>
      <p>This is a Vue3 project with sidebar and header functionality</p>
    </div>
  </main>
</template>

<script setup>
// Main content component
</script>

<style scoped>
.main-content {
  flex: 1;
  padding: 24px;
  overflow: hidden;
  background-color: var(--background-primary);
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content-placeholder {
  text-align: center;
  padding: 40px;
  background-color: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.content-placeholder h2 {
  font-size: 24px;
  margin-bottom: 16px;
  color: #333;
}

.content-placeholder p {
  font-size: 16px;
  color: #666;
}
</style>