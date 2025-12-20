// --------------------------------------------------------------------------------
// File:        main.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Main content component for the application.
// --------------------------------------------------------------------------------
<template>
  <main class="main-content">
    <div class="content-wrapper">
      <!-- Page Header -->
      <div class="page-header">
        <h1>Dashboard</h1>
        <p>This dashboard provides an overview of system performance and resource utilization.</p>
      </div>

      <!-- Action Bar -->
      <div class="action-bar">
        <div class="search-bar">
          <input type="text" placeholder="Find resources..." class="input-primary">
          <button class="btn-third">Search</button>
        </div>
      </div>

      <!-- Data Table -->
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th><input type="checkbox" class="checkbox-primary"></th>
              <th>Resource Group</th>
              <th>Access Level</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="group in deviceGroups" :key="group.id">
              <td><input type="checkbox" class="checkbox-primary"></td>
              <td>{{ group.name }}</td>
              <td>{{ group.authMethod }}</td>
              <td>
                <button class="btn-third">Modify</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="pagination">
        <div class="page-info">
          <span>1-10/10</span>
        </div>
        <div class="page-controls">
          <button class="btn-secondary" disabled>First</button>
          <button class="btn-third" disabled>Previous</button>
          <span>1</span>
          <button class="btn-third" disabled>Next</button>
          <button class="btn-secondary" disabled>Last</button>
        </div>
        <div class="page-size">
          <label for="pageSize">Rows:</label>
          <select id="pageSize" class="select-primary">
            <option value="50">50</option>
            <option value="100">100</option>
            <option value="200">200</option>
          </select>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref } from 'vue';

// Mock data for resource groups
const deviceGroups = ref([
  { id: 1, name: 'Alpha Group', authMethod: 'Basic Access' },
  { id: 2, name: 'Beta Cluster', authMethod: 'Enhanced Security' },
  { id: 3, name: 'Gamma Unit', authMethod: 'Advanced Protocol' },
  { id: 4, name: 'Delta System', authMethod: 'Standard Access' },
  { id: 5, name: 'Epsilon Network', authMethod: 'Basic Access' },
  { id: 6, name: 'Zeta Framework', authMethod: 'Enhanced Security' },
  { id: 7, name: 'Eta Service', authMethod: 'Advanced Protocol' },
  { id: 8, name: 'Theta Environment', authMethod: 'Restricted Access' },
  { id: 9, name: 'Iota Division', authMethod: 'Basic Access' },
  { id: 10, name: 'Kappa Platform', authMethod: 'Enhanced Security' }
]);
</script>

<style scoped>
.main-content {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  background-color: var(--background-modal);
  width: 100%;
}

.content-wrapper {
  width: 100%;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.page-header {
  margin: 0;
  background-color: var(--background-modal-title);
  padding: 8px 12px;
  width: 100%;
  box-sizing: border-box;
}

.page-header h1 {
  font-size: 20px;
  color: #1e293b;
  margin: 0;
  line-height: 1.5;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.page-header p {
  font-size: 14px;
  color: #64748b;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: var(--background-secondary);
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.search-bar {
  display: flex;
  gap: 8px;
}

.table-container {
  background-color: var(--background-secondary);
  border: 1px solid var(--border-color);
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th, .data-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

/* Keep headers visible */
.data-table th {
  background-color: var(--background-primary);
  font-weight: 600;
  color: #1e293b;
}

/* Normal styling for data cells */
.data-table td {
  background-color: transparent;
  color: #1e293b;
  opacity: 1;
  pointer-events: auto;
  user-select: auto;
}

.data-table tr:last-child td {
  border-bottom: none;
}

/* Add btn-third style based on btn-second with --color-second background */
.btn-third {
  background-color: var(--color-second);
  color: white;
  border: 0px;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: background-color 0.2s ease;
  justify-content: center;
  border-radius: 0;
}

.btn-third:hover {
  background-color: var(--color-second);
  opacity: 0.9;
  border-color: transparent;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: var(--background-secondary);
}

.page-info {
  font-size: 14px;
  color: #64748b;
}

.page-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-controls button {
  margin: 0;
}

.page-controls span {
  padding: 4px 8px;
  font-weight: 600;
}

.page-size {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Add virtual text styling for sensitive information */
.virtual-text {
  /* Virtualize text for sensitive information */
  background: repeating-linear-gradient(
    90deg,
    #6b7280 0px,
    #6b7280 10px,
    #9ca3af 10px,
    #9ca3af 20px
  );
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  pointer-events: none;
  user-select: none;
}

/* Remove virtual text effect from all elements */


/* Remove all border-radius */
.btn-third,
.btn-secondary,
.input-primary,
.select-primary,
.checkbox-primary,
.action-btn,
.radiobox-primary::before,
.radiobox-primary::after,
.info-icon,
.modal-footer,
.modal-container,
.form-section-title-primary {
  border-radius: 0;
}
</style>