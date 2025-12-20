// --------------------------------------------------------------------------------
// File:        sidebar.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Sidebar component for navigation menu.
// --------------------------------------------------------------------------------
<template>
  <aside :class="['sidebar', 'background-primary', isHidden ? 'sidebar-hidden' : 'sidebar-visible']">
    <div class="sidebar-header">
      <h1 class="sidebar-title">
        <div class="flex items-center">
          <i class="fa fa-globe sidebar-icon"></i>
          <div>
            <div class="sidebar-main-title text-shadow-sm">TRAE CN</div>
          </div>
        </div>
      </h1>
    </div>
    <nav class="sidebar-nav">
      <ul>
        <li class="nav-group-title text-shadow-sm">
          Dashboard
        </li>
        <li>
          <a href="#" @click.prevent="openOverview" class="nav-item text-shadow-sm" :class="{ 'nav-item-active': !isAnyModuleOpen }"
             draggable="false">
            <i class="fa fa-dashboard icon-blue"></i>
            Overview
          </a>
        </li>
        <li class="nav-group-title text-shadow-sm">
          Assets
        </li>
        <li>
          <a href="#" @click.prevent="openUserManagement" class="nav-item text-shadow-sm"
             :class="{ 'nav-item-active': pageStore.showUserManagement }" draggable="false">
            <i class="fa fa-users icon-purple"></i>
            Users
          </a>
        </li>
        <li>
          <a href="#" @click.prevent="openDeviceManagement" class="nav-item text-shadow-sm"
             :class="{ 'nav-item-active': pageStore.showDeviceManagement }" draggable="false">
            <i class="fa fa-print icon-pink"></i>
            Devices
          </a>
        </li>
        <li>
          <a href="#" @click.prevent="openDeviceGroupManagement" class="nav-item text-shadow-sm"
             :class="{ 'nav-item-active': pageStore.showDeviceGroupManagement }" draggable="false">
            <i class="fa fa-print icon-pink"></i>
            Groups
          </a>
        </li>
        <li>
          <a href="#" @click.prevent="openDeptManagement" class="nav-item text-shadow-sm"
             :class="{ 'nav-item-active': pageStore.showDepartmentManagement }" draggable="false">
            <i class="fa fa-sitemap icon-red"></i>
            Departments
          </a>
        </li>
        <li class="nav-group-title text-shadow-sm">
          Reports
        </li>
        <li>
          <a href="#" @click.prevent="openUserAccountingReport" class="nav-item text-shadow-sm" 
             :class="{ 'nav-item-active': pageStore.showUserAccountingReport }" draggable="false">
            <i class="fa fa-file-text-o icon-orange"></i>
            User Reports
          </a>
        </li>
        <li>
          <a href="#" @click.prevent="openPrintTrackingReport" class="nav-item text-shadow-sm" 
             :class="{ 'nav-item-active': pageStore.showPrintTrackingReport }" draggable="false">
            <i class="fa fa-file-text-o icon-yellow"></i>
            Logs
          </a>
        </li>
        <li>
          <a href="#" @click.prevent="openReportGeneration" class="nav-item text-shadow-sm" 
             :class="{ 'nav-item-active': pageStore.showReportGeneration }" draggable="false">
            <i class="fa fa-file-text-o icon-teal"></i>
            Generate
          </a>
        </li>
        <li>
          <a href="#" @click.prevent="openDepartmentAccountingReport" class="nav-item text-shadow-sm"
             :class="{ 'nav-item-active': pageStore.showDepartmentAccountingReport }" draggable="false">
            <i class="fa fa-file-text-o icon-green"></i>
            Dept Reports
          </a>
        </li>

      </ul>
    </nav>
    <div class="sidebar-user-footer">
      <div class="sidebar-user-info">
        <img src="../assets/images/default-avatar.svg" alt="User Avatar" class="sidebar-avatar">
        <div class="user-details">
          <p class="sidebar-username text-shadow-sm">Admin</p>
          <a href="mailto:service@traecn.com" class="sidebar-user-email text-shadow-sm">
            service@traecn.com
          </a>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>

import {usePageNavigationStore} from '../data-sources/sidebar.js'
import {computed, onUnmounted, ref} from 'vue'

const props = defineProps({
  isHidden: {type: Boolean, default: false}
});

const pageStore = usePageNavigationStore()

const openDeviceGroupManagement = () => {
  pageStore.closeAllPages()
  pageStore.openDeviceGroupManagement()
}

const openOverview = () => {
  pageStore.closeAllPages()
}

const isAnyModuleOpen = computed(() => {
  return pageStore.showUserManagement ||
      pageStore.showAdvancedUserManagement ||
      pageStore.showDeviceManagement ||
      pageStore.showDeviceGroupManagement ||
      pageStore.showDepartmentManagement ||
      pageStore.showDepartmentAccountingReport ||
      pageStore.showUserAccountingReport ||
      pageStore.showPrintTrackingReport ||
      pageStore.showReportGeneration;
});

const openDeptManagement = () => {
  try {

    pageStore.closeAllPages();

    if (pageStore && typeof pageStore.openDepartmentManagement === 'function') {
      pageStore.openDepartmentManagement();
      console.log('Sidebar: Trigger department management open');
    } else {
      console.error('Sidebar: pageStore.openDepartmentManagement method does not exist');
    }
  } catch (error) {
    console.error('Sidebar: Failed to open department management:', error);
  }
}

const openUserManagement = () => {
  try {

    pageStore.closeAllPages();

    if (pageStore && typeof pageStore.openUserManagement === 'function') {
      pageStore.openUserManagement();
      console.log('Sidebar: Trigger user management open');
    } else {
      console.error('Sidebar: pageStore.openUserManagement method does not exist');
    }
  } catch (error) {
    console.error('Sidebar: Failed to open user management:', error);
  }
}

const openDeviceManagement = () => {
  try {

    pageStore.closeAllPages();

    if (pageStore && typeof pageStore.openDeviceManagement === 'function') {
      pageStore.openDeviceManagement();
      console.log('Sidebar: Trigger device management open');
    } else {
      console.error('Sidebar: pageStore.openDeviceManagement method does not exist');
    }
  } catch (error) {
    console.error('Sidebar: Failed to open device management:', error);
  }
}

const openDepartmentAccountingReport = () => {
  try {

    pageStore.closeAllPages();
    pageStore.openDepartmentAccountingReport()
  } catch (error) {
    console.error('Sidebar: Failed to open department accounting report:', error);
  }
}

const openUserAccountingReport = () => {
  try {

    pageStore.closeAllPages();
    pageStore.openUserAccountingReport()
  } catch (error) {
    console.error('Sidebar: Failed to open user accounting report:', error);
  }
}

const openPrintTrackingReport = () => {
  try {

    pageStore.closeAllPages();
    pageStore.openPrintTrackingReport()
  } catch (error) {
    console.error('Sidebar: Failed to open operation log audit:', error);
  }
}

const openReportGeneration = () => {
  try {

    pageStore.closeAllPages();
    pageStore.openReportGeneration()
  } catch (error) {
    console.error('Sidebar: Failed to open comprehensive report generation:', error);
  }
}
</script>

<style scoped>

.sidebar-title {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
}

.sidebar-subtitle {
  font-size: 14px;
  color: #bbbbbb;
  margin: 4px 0 0 0;
}

.flex {
  display: flex;
}

.items-center {
  align-items: center;
}

.justify-between {
  justify-content: space-between;
}

.mr-2 {
  margin-right: 0.5rem;
}

.ml-2 {
  margin-left: 0.5rem;
}

.whitespace-nowrap {
  white-space: nowrap;
}

.text-secondary {
  color: #64748b;
}

.align-middle {
  vertical-align: middle;
}

.select-none {
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.transition {
  transition: all 0.3s ease;
}

.duration-200 {
  transition-duration: 200ms;
}

.ease-in {
  transition-timing-function: ease-in;
}

.sidebar {
  width: 250px;
  height: 100%;
  border-right: 1px solid var(--border-color);
  padding: 11px 16px;
  overflow-y: auto;
  transition: all 0.3s ease;
  flex-shrink: 0;
  overflow: hidden;
}

.sidebar-hidden {
  width: 0;
  padding: 0;
  border-right: none;
  visibility: hidden;
  opacity: 0;
}

.sidebar-visible {
  width: 250px;
  padding: 11px 16px;
  border-right: 1px solid var(--border-color);
  visibility: visible;
  opacity: 1;
}

.sidebar-header {
  border-bottom: 1px solid var(--border-color);
}

.sidebar-title {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
}

.sidebar-nav {
  margin-bottom: 32px;
}

.sidebar-nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-group-title {
  font-size: 14px;
  font-weight: 600;
  color: #000000;
  text-transform: none;
  letter-spacing: normal;
  margin-top: 20px;
  margin-bottom: 8px;
  font-family: Arial, Helvetica, 'Helvetica Neue', sans-serif;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 10px 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 15px;
  font-weight: 400;
  color: #475569;
  text-decoration: none;
  font-family: Arial, Helvetica, 'Helvetica Neue', sans-serif;
  line-height: 1.5;
}

.nav-item i {
  width: 14px;
  height: 14px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 14px;
  margin-right: 6px;
}

.nav-item:hover {
  color: #1e293b;

  font-weight: 900;
}

.nav-item-active {
  background-color: var(--nav-item-active-background-primary);
  color: var(--color-primary);
  font-weight: 600;

}

.nav-item-active:has(.icon-blue) {
  color: #475569;
}

.sidebar-user-footer {
  margin-top: auto;
  padding: 24px 0 24px 0;
  border-top: 1px solid var(--border-color);
}

.sidebar-user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.sidebar-avatar {
  width: 40px;
  height: 40px;
  object-fit: cover;

  transition: all 0.3s ease;
  cursor: pointer;
  border-radius: 50%;
}

.sidebar-avatar:hover {
  transform: scale(1.05);
}

.user-details {
  flex: 1;
}

.sidebar-username {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 0px 0;
}

.sidebar-user-email {
  font-size: 12px;
  color: #64748b;
  text-decoration: none;
}

.sidebar-user-email:hover {
  text-decoration: underline;
}

.toggle-checkbox {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-label {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 20px;
  background-color: #e2e8f0;

  transition: background-color 0.3s;
}

.toggle-checkbox:checked + .toggle-label {
  background-color: var(--color-primary);
}

.toggle-label::before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background-color: white;

  transition: transform 0.3s;
}

.toggle-checkbox:checked + .toggle-label::before {
  transform: translateX(20px);
}

.sidebar-main-title {
  font-size: 24px;
  font-weight: 900;
  font-family: Impact, Haettenschweiler, 'Arial Narrow Bold', sans-serif;
  color: var(--color-danger);
  margin-top: 3px;
  margin-left: 8px;
  margin-bottom: 5px;
  cursor: pointer;
}

.sidebar-icon {
  color: var(--color-danger);
  font-size: 26px;
  position: relative;
  top: -1px;
  left: 0px;
}



.sidebar-title {
  /* 使用通用的text-shadow-sm类 */
}

.icon-blue,
.icon-purple,
.icon-pink,
.icon-red,
.icon-orange,
.icon-yellow,
.icon-green,
.icon-teal,
.icon-cyan {
  color: #000000;
}

.icon-indigo {
  color: #8B5CF6;
}
</style>
