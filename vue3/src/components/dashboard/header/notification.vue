// --------------------------------------------------------------------------------
// File:        notification.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Notification center panel component for displaying and managing system notifications.
// --------------------------------------------------------------------------------
<template>
  <!-- Notification Center Panel -->
  <transition name="modal">
    <div v-if="notificationStore.showNotificationPanel" class="notification-panel-container"
         @click="notificationStore.closeAllPanels()">
      <div class="notification-panel" @click.stop>
        <div class="panel-header">
          <div class="header-left">
            <i class="fa fa-bell header-icon"></i>
            <h3>Notification Center</h3>
            <span class="unread-count" v-if="notificationStore.unreadNotifications.length > 0">
              {{ notificationStore.unreadNotifications.length }}
            </span>
          </div>
          <button class="btn-modal-close" @click="notificationStore.closeAllPanels()">
            <i class="fa fa-times"></i>
          </button>
        </div>

        <div class="panel-tabs">
          <button
              class="tab-btn"
              :class="{ active: activeTab === 'all' }"
              @click="activeTab = 'all'"
          >
            <i class="fa fa-list"></i>
            All ({{ notificationStore.notifications.length }})
          </button>
          <button
              class="tab-btn"
              :class="{ active: activeTab === 'unread' }"
              @click="activeTab = 'unread'"
          >
            <i class="fa fa-circle"></i>
            Unread ({{ notificationStore.unreadNotifications.length }})
          </button>
        </div>

        <div class="panel-actions">
          <button class="action-btn" @click="markAllAsRead"
                  :disabled="notificationStore.unreadNotifications.length === 0">
            <i class="fa fa-check"></i>
            Mark All as Read
          </button>
          <button class="action-btn" @click="clearAll" :disabled="notificationStore.notifications.length === 0">
            <i class="fa fa-trash"></i>
            Clear All
          </button>
          <button class="action-btn" @click="refreshNotifications">
            <i class="fa fa-refresh"></i>
            Refresh
          </button>
        </div>

        <div class="notification-list">
          <transition-group name="notification" tag="div">
            <div
                v-for="notification in filteredNotifications"
                :key="notification.id"
                class="notification-item"
                :class="{
                unread: !notification.isRead,
                high: notification.priority === 'high',
                medium: notification.priority === 'medium'
              }"
                @click="showNotificationDetail(notification)"
            >
              <div class="notification-icon" :class="notification.type">
                <i :class="getNotificationIcon(notification.type)"></i>
              </div>
              <div class="notification-content">
                <div class="notification-header">
                  <span class="notification-title">{{ notification.title }}</span>
                  <span class="notification-time">{{ formatTime(notification.createdAt) }}</span>
                </div>
                <div class="notification-text">{{ notification.content }}</div>
                <div class="notification-meta">
                  <span class="notification-category" :class="notification.category.toLowerCase()">
                    {{ notification.category }}
                  </span>
                  <span class="notification-priority" :class="notification.priority">
                    {{ getPriorityText(notification.priority) }}
                  </span>
                  <span class="notification-source">{{ notification.source || 'System' }}</span>
                </div>
              </div>
              <div class="notification-status">
                <div v-if="!notification.isRead" class="unread-dot" :class="notification.priority"></div>
                <button class="delete-btn" @click.stop="deleteNotification(notification.id)">
                  <i class="fa fa-times"></i>
                </button>
              </div>
            </div>
          </transition-group>

          <div v-if="filteredNotifications.length === 0" class="empty-state">
            <i class="fa fa-bell-slash"></i>
            <p>{{ emptyStateText }}</p>
            <button v-if="activeTab === 'unread' && notificationStore.notifications.length > 0"
                    class="btn-secondary" @click="activeTab = 'all'">
              View All Notifications
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>

  <!-- Notification Detail Modal -->
  <transition name="modal">
    <div v-if="currentNotification" class="notification-detail-modal">
      <div class="modal-overlay" @click="closeDetailModal"></div>
      <div class="modal-content" :class="currentNotification.priority">
        <div class="modal-header">
          <div class="header-info">
            <div class="modal-icon" :class="currentNotification.type">
              <i :class="getNotificationIcon(currentNotification.type)"></i>
            </div>
            <div class="modal-title-section">
              <h4>{{ currentNotification.title }}</h4>
              <div class="modal-meta">
                <span class="detail-category" :class="currentNotification.category.toLowerCase()">
                  {{ currentNotification.category }}
                </span>
                <span class="detail-priority" :class="currentNotification.priority">
                  {{ getPriorityText(currentNotification.priority) }}
                </span>
                <span class="detail-time">{{ formatTime(currentNotification.createdAt) }}</span>
              </div>
            </div>
          </div>
          <button class="btn-modal-close" @click="closeDetailModal">
            <i class="fa fa-times"></i>
          </button>
        </div>
        <div class="modal-body">
          <div class="detail-content">
            {{ currentNotification.content }}
          </div>
          <div v-if="currentNotification.details" class="detail-extra">
            <div class="detail-section">
              <h5>Details</h5>
              <p>{{ currentNotification.details }}</p>
            </div>
          </div>
          <div v-if="currentNotification.actions" class="detail-actions">
            <button
                v-for="action in currentNotification.actions"
                :key="action.id"
                class="action-btn"
                :class="action.type"
                @click="handleAction(action)"
            >
              <i :class="action.icon"></i>
              {{ action.text }}
            </button>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="closeDetailModal">
            <i class="fa fa-times"></i>
            Close
          </button>
          <button v-if="!currentNotification.isRead" class="btn-secondary" @click="markAsReadAndClose">
            <i class="fa fa-check"></i>
            Mark as Read
          </button>
          <button class="btn-secondary" @click="deleteNotification(currentNotification.id)">
            <i class="fa fa-trash"></i>
            Delete
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import {computed, ref} from 'vue'
import {useNotificationStore} from '../../../data-sources/dashboard/header/notification.js'

const notificationStore = useNotificationStore()
const activeTab = ref('all')
const currentNotification = ref(null)

const filteredNotifications = computed(() => {
  if (activeTab.value === 'unread') {
    return notificationStore.unreadNotifications
  }
  return notificationStore.notifications
})

const emptyStateText = computed(() => {
  if (activeTab.value === 'unread') {
    return 'No unread notifications'
  }
  return 'No notifications'
})

const getNotificationIcon = (type) => {
  const iconMap = {
    'system': 'fa fa-cog',
    'security': 'fa fa-shield',
    'message': 'fa fa-envelope',
    'task': 'fa fa-tasks',
    'warning': 'fa fa-exclamation-triangle',
    'info': 'fa fa-info-circle',
    'success': 'fa fa-check-circle',
    'error': 'fa fa-times-circle',
    'update': 'fa fa-download',
    'reminder': 'fa fa-clock-o',
    'print': 'fa fa-print',
    'device': 'fa fa-desktop',
    'quota': 'fa fa-pie-chart'
  }
  return iconMap[type] || 'fa fa-bell'
}

const getPriorityText = (priority) => {
  const priorityMap = {
    'high': 'High Priority',
    'medium': 'Medium Priority',
    'low': 'Low Priority'
  }
  return priorityMap[priority] || 'Normal'
}

const formatTime = (timestamp) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time

  if (diff < 60000) {
    return 'Just now'
  } else if (diff < 3600000) {
    return `${Math.floor(diff / 60000)}m ago`
  } else if (diff < 86400000) {
    return `${Math.floor(diff / 3600000)}h ago`
  } else if (diff < 604800000) {
    return `${Math.floor(diff / 86400000)}d ago`
  } else {
    return time.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}

const showNotificationDetail = (notification) => {
  currentNotification.value = notification
  if (!notification.isRead) {
    notificationStore.markAsRead(notification.id)
  }
}

const closeDetailModal = () => {
  currentNotification.value = null
}

const markAsReadAndClose = () => {
  if (currentNotification.value) {
    notificationStore.markAsRead(currentNotification.value.id)
    currentNotification.value = null
  }
}

const deleteNotification = (id) => {
  notificationStore.deleteNotification(id)
  if (currentNotification.value && currentNotification.value.id === id) {
    currentNotification.value = null
  }
}

const markAllAsRead = () => {
  notificationStore.markAllAsRead()
}

const clearAll = () => {
  if (confirm('Are you sure you want to clear all notifications? This action cannot be undone.')) {
    notificationStore.clearAllNotifications()
  }
}

const refreshNotifications = () => {
  notificationStore.refreshNotifications()
}

const handleAction = (action) => {
  if (action.handler) {
    action.handler()
  } else {

    switch (action.id) {
      case 'view':

        break
      case 'accept':

        break
      case 'reject':

        break
      case 'remind':

        break
      default:
        console.log('执行动作:', action.text)
    }
  }
  closeDetailModal()
}
</script>

<style scoped>
.notification-panel-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: transparent;
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .notification-panel,
.modal-leave-active .notification-panel {
  transition: opacity 0.3s ease;
}

.modal-enter-from .notification-panel,
.modal-leave-to .notification-panel {
  opacity: 0;
}

.notification-panel {
  position: fixed;
  top: 60px;
  right: 20px;
  width: 400px;
  height: 92vh;
  background: white;

  border: 1px solid var(--border-color);
  z-index: 10000;
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  background: #f8fafc;

}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.close-btn {
  background: none;
  border: none;
  font-size: 16px;
  color: #64748b;
  cursor: pointer;
  padding: 4px;

  transition: all 0.2s;
}

.close-btn:hover {
  background: #e2e8f0;
  color: #1e293b;
}

.panel-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
}

.tab-btn {
  flex: 1;
  padding: 12px 16px;
  background: none;
  border: none;
  font-size: 14px;
  color: #64748b;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
}

.tab-btn.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  font-weight: 600;
}

.tab-btn:hover:not(.active) {
  color: #1e293b;
  background: #f8fafc;
}

.panel-actions {
  display: flex;
  gap: 8px;
  padding: 12px 20px;
  border-bottom: 1px solid var(--border-color);
}

.action-btn {
  padding: 6px 12px;
  background: #f1f5f9;
  border: 1px solid var(--border-color);

  font-size: 12px;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.action-btn:hover:not(:disabled) {
  background: #e2e8f0;
  color: #1e293b;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.notification-list {
  flex: 1;
  overflow-y: auto;
  max-height: calc(92vh - 200px); /* 统一使用与系统设置相同的内容区域高度 */
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.notification-item:hover {
  background: #f8fafc;
}

.notification-item.unread {
  background: #f0f9ff;
  border-left: 3px solid var(--color-primary);
}

.notification-item.high {
  border-left-color: var(--border-color);
}

.notification-icon {
  width: 40px;
  height: 40px;
  background: #f1f5f9;

  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notification-icon i {
  font-size: 16px;
  color: #64748b;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.notification-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-time {
  font-size: 12px;
  color: #64748b;
  flex-shrink: 0;
}

.notification-text {
  font-size: 13px;
  color: #475569;
  line-height: 1.4;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-meta {
  display: flex;
  gap: 8px;
  align-items: center;
}

.notification-category {
  font-size: 11px;
  color: #64748b;
  background: #f1f5f9;
  padding: 2px 6px;

}

.notification-priority {
  font-size: 11px;
  padding: 2px 6px;

  font-weight: 500;
}

.notification-priority.high {
  background: #fef2f2;
  color: #dc2626;
}

.notification-priority.medium {
  background: #fffbeb;
  color: #d97706;
}

.notification-priority.normal {
  background: #f0f9ff;
  color: #0284c7;
}

.notification-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.unread-dot {
  width: 8px;
  height: 8px;
  background: var(--color-primary);

}

.delete-btn {
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;

  transition: all 0.2s;
  opacity: 0;
}

.notification-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #64748b;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  font-size: 14px;
  margin: 0;
}

.notification-detail-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1001;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  animation: fadeIn 0.3s ease-out;
  z-index: 1000;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-content {
  position: relative;
  background: var(--background-primary);

  width: 560px;
  max-width: 90vw;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  animation: modalSlideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1001;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-50px) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, #fafafa, var(--background-primary));

}

.header-info {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  flex: 1;
}

.modal-icon {
  width: 56px;
  height: 56px;

  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.modal-icon i {
  color: white;
  font-size: 24px;
}

.modal-title-section {
  flex: 1;
}

.modal-header h4 {
  margin: 0 0 12px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
  line-height: 1.4;
}

.modal-meta {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
  font-size: 13px;
}

.detail-category {
  background: #f0f0f0;
  padding: 6px 14px;

  color: #666;
  font-weight: 500;
  font-size: 12px;
  transition: all 0.2s;
}

.detail-category:hover {
  background: #e0e0e0;
  transform: translateY(-1px);
}

.detail-category.system {
  background: #e6f7ff;
  color: #1890ff;
}

.detail-category.security {
  background: #fff2f0;
  color: #ff4d4f;
}

.detail-category.message {
  background: #f6ffed;
  color: #52c41a;
}

.detail-category.task {
  background: #fff7e6;
  color: #faad14;
}

.detail-priority {
  padding: 6px 14px;

  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  transition: all 0.2s;
}

.detail-priority.high {
  background: linear-gradient(135deg, #ff4d4f, #ff7875);
  color: white;
}

.detail-priority.medium {
  background: linear-gradient(135deg, #faad14, #ffc53d);
  color: white;
}

.detail-priority.low {
  background: linear-gradient(135deg, #52c41a, #73d13d);
  color: white;
}

.detail-time {
  color: #999;
  font-size: 12px;
  margin-left: auto;
}

.modal-body {
  padding: 24px;
  flex: 1;
  overflow-y: auto;
  background: var(--background-primary);
}

.detail-content {
  font-size: 15px;
  line-height: 1.7;
  color: #333;
  white-space: pre-wrap;
  margin-bottom: 24px;
  background: white;
  padding: 20px;

}

.detail-extra {
  margin-bottom: 24px;
}

.detail-section h5 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-section p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: #555;
  background: white;
  padding: 20px;

  border-left: 4px solid var(--border-color);
}

.detail-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
  background: linear-gradient(135deg, #fafafa, var(--background-primary));

}

.btn-secondary {
  padding: 10px 20px;
  background: white;
  border: 1px solid var(--border-color);

  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-secondary:hover {
  border-color: var(--border-color);
  color: #667eea;
  transform: translateY(-1px);
}


.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-content {
  opacity: 0;
  transform: translateY(-50px) scale(0.9);
}

.modal-leave-to .modal-content {
  opacity: 0;
  transform: translateY(-50px) scale(0.9);
}
</style>
