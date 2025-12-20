// --------------------------------------------------------------------------------
// File:        header.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Main header component with sidebar toggle and notification icons.
// --------------------------------------------------------------------------------
<template>
  <header class="main-header">
    <button class="sidebar-toggle text-shadow" @click="toggleSidebar">
      <i class="fa fa-bars"></i>
    </button>
    <div class="header-icons">
      <div class="version-info text-shadow">
        This page is only optimized for 3840x2160 resolution
      </div>
      <div class="version-info text-shadow">
        v2.3.1-rc.5+build.20251119
      </div>
      <!-- Real-time time display -->
      <div class="current-time text-shadow">
        <span class="date-time">{{ currentDateTime }}</span>
        <span class="weekday-greeting" :class="{ 'sunday-red': currentWeekday === 'Sunday' }">{{
            currentWeekday
          }}, <span class="greeting-part">{{ currentGreeting }}</span></span>
      </div>
      <button class="header-icon-btn text-shadow" @click="notificationStore.toggleNotificationPanel()">
        <i class="fa fa-bell"></i>
        <span v-if="notificationStore.notificationCount > 0" class="notification-dot"></span>
      </button>
      <button class="header-icon-btn text-shadow" @click="headerStore.toggleEmailPanel()">
        <i class="fa fa-envelope"></i>
        <span v-if="headerStore.emailCount > 0" class="notification-badge">{{ headerStore.emailCount }}</span>
      </button>
      <button class="header-icon-btn text-shadow" @click="headerStore.toggleSettingsPanel()">
        <i class="fa fa-cog"></i>
        <span class="new-feature-tag">New</span>
      </button>
      <div class="version-info logo">
        TRAE CN
      </div>
    </div>
  </header>

  <!-- Notification Panel -->
  <BellNotification/>

  <!-- Email Panel -->
  <EmailMessage/>

  <!-- System Settings Panel -->
  <SystemSettings/>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from 'vue'
import {useNotificationStore} from '../../data-sources/dashboard/header/notification.js'
import {useHeaderStore} from '../../data-sources/dashboard/header/header.js'
import BellNotification from './header/notification.vue'
import EmailMessage from './header/email.vue'
import SystemSettings from './header/settings.vue'

const notificationStore = useNotificationStore()
const headerStore = useHeaderStore()

const emitEvent = defineEmits(['toggle-sidebar'])

const toggleSidebar = () => {
  emitEvent('toggle-sidebar')
}

const currentDateTime = ref('')
const currentWeekday = ref('')
const currentGreeting = ref('')
const showColon = ref(true)

const weekdays = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']

// Get time greeting
const getGreeting = () => {
  const now = new Date()
  const hour = now.getHours()

  if (hour >= 5 && hour < 12) {
    return 'Good morning'
  } else if (hour >= 12 && hour < 14) {
    return 'Good noon'
  } else if (hour >= 14 && hour < 18) {
    return 'Good afternoon'
  } else if (hour >= 18 && hour < 22) {
    return 'Good evening'
  } else if (hour >= 22 && hour < 24) {
    return "It's late, time to rest"
  } else if (hour >= 0 && hour < 3) {
    return "It's late, time to rest"
  } else {
    return "It's dawn, the rooster will crow"
  }
}

const updateDateTime = () => {
  const now = new Date()

  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')

  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  const seconds = String(now.getSeconds()).padStart(2, '0')

  const weekdayIndex = now.getDay()
  currentWeekday.value = weekdays[weekdayIndex]
  currentGreeting.value = getGreeting()

  const colon = showColon.value ? ':' : ' '
  currentDateTime.value = `${year}/${month}/${day} ${hours}${colon}${minutes}${colon}${seconds}`
}

const toggleColon = () => {
  showColon.value = !showColon.value
  updateDateTime()
}

let timer = null
let colonTimer = null

onMounted(() => {
  updateDateTime()
  timer = setInterval(updateDateTime, 1000)
  colonTimer = setInterval(toggleColon, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.notification-badge {
  position: absolute;
  top: 4px;
  right: 2px;
  background-color: red;
  color: white;
  font-size: 6px;
  font-weight: bold;
  padding: 2px 2px;
  min-width: 14px;
  white-space: nowrap;
}

.new-feature-tag {
  position: absolute;
  top: 4px;
  right: 0px;
  background-color: red;
  color: white;
  font-size: 6px;
  font-weight: bold;
  padding: 2px 2px;
  min-width: 14px;
  white-space: nowrap;
}

/* Ensure consistent styling */
.main-header {
  background-color: var(--background-primary);
  border-bottom: 1px solid var(--border-color);
  padding: 8px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.version-info {
  font-size: 14px;
  color: var(--color-primary);
  margin-right: 8px;
  padding: 4px 8px;
  font-weight: 500;
}

.logo {
  color: #ff0000;
  font-weight: 900;
  margin-right: 0;
  background-color: transparent;
  font-family: Impact, Haettenschweiler, 'Arial Narrow Bold', sans-serif
}

.sidebar-toggle {
  background: none;
  border: none;
  font-size: 20px;
  color: #64748b;
  cursor: pointer;
  padding: 8px;
  transition: all 0.3s ease;
}

.sidebar-toggle:hover {
  background-color: var(--background-primary);
  color: #1e293b;
}

.header-icons {
  display: flex;
  align-items: center;
}

.header-icon-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #64748b;
  cursor: pointer;
  padding: 8px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.header-icon-btn:hover {
  background-color: var(--background-primary);
  color: #1e293b;
}

.notification-dot {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 8px;
  height: 8px;
  background-color: #ef4444;
  border-radius: 50%;
}

.current-time {
  font-size: 14px;
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 8px;
}

.date-time {
  color: var(--color-primary);
  font-weight: 600;
  font-family: 'Courier New', Courier, monospace;
}

.weekday-greeting {
  color: var(--color-primary);
  font-weight: 500;
}

.weekday-greeting .greeting-part {
  color: #ec4899;
  font-weight: 600;
}

.weekday-greeting.sunday-red {
  color: #ef4444;
}
</style>
