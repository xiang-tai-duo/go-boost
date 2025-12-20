// --------------------------------------------------------------------------------
// File:        email.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Email panel component for displaying and managing emails.
// --------------------------------------------------------------------------------
<template>
  <!-- Email Panel -->
  <transition name="modal">
    <div v-if="headerStore.showEmailPanel" class="email-panel-container" @click="headerStore.toggleEmailPanel()">
      <div class="email-panel" @click.stop>
        <div class="panel-header">
          <div class="header-left">
            <i class="fa fa-envelope header-icon"></i>
            <h3>Email</h3>
            <span class="unread-count" v-if="headerStore.unreadEmails.length > 0">
              {{ headerStore.unreadEmails.length }}
            </span>
          </div>
          <button class="btn-modal-close" @click="headerStore.toggleEmailPanel()">
            <i class="fa fa-times"></i>
          </button>
        </div>

        <div class="panel-tabs">
          <button
              class="tab-btn"
              :class="{ active: activeTab === 'inbox' }"
              @click="activeTab = 'inbox'"
          >
            <i class="fa fa-inbox"></i>
            Inbox ({{ headerStore.emailsByStatus['inbox']?.length || 0 }})
          </button>
          <button
              class="tab-btn"
              :class="{ active: activeTab === 'unread' }"
              @click="activeTab = 'unread'"
          >
            <i class="fa fa-circle"></i>
            Unread ({{ headerStore.unreadEmails.length }})
          </button>
          <button
              class="tab-btn"
              :class="{ active: activeTab === 'sent' }"
              @click="activeTab = 'sent'"
          >
            <i class="fa fa-paper-plane"></i>
            Sent ({{ headerStore.emailsByStatus['sent']?.length || 0 }})
          </button>
        </div>

        <div class="panel-actions">
          <button class="action-btn" @click="markAllAsRead" :disabled="headerStore.unreadEmails.length === 0">
            <i class="fa fa-check"></i>
            Mark All as Read
          </button>
          <button class="action-btn" @click="composeNewEmail">
            <i class="fa fa-pencil"></i>
            Compose
          </button>
          <button class="action-btn" @click="refreshEmails">
            <i class="fa fa-refresh"></i>
            Refresh
          </button>
        </div>

        <div class="email-list">
          <transition-group name="email" tag="div">
            <div
                v-for="email in filteredEmails"
                :key="email.id"
                class="email-item"
                :class="{ unread: !email.isRead }"
                @click="showEmailDetail(email)"
            >
              <div class="email-checkbox">
                <input type="checkbox" :id="`email-${email.id}`" @click.stop class="checkbox-primary">
              </div>
              <div class="email-sender"
                   :class="{ 'has-attachment': email.attachments && email.attachments.length > 0 }">
                <i class="fa fa-user-circle"></i>
                <span class="sender-name">{{ email.fromName }}</span>
              </div>
              <div class="email-content">
                <div class="email-header">
                  <span class="email-subject">{{ email.subject }}</span>
                  <span class="email-time">{{ formatTime(email.createdAt) }}</span>
                </div>
                <div class="email-preview">{{ email.content.substring(0, 80) }}...</div>
                <div class="email-meta">
                  <span v-if="email.attachments && email.attachments.length > 0" class="attachment-icon">
                    <i class="fa fa-paperclip"></i>
                    {{ email.attachments.length }} Attachments
                  </span>
                  <span class="email-priority" :class="email.priority">
                    {{ getPriorityText(email.priority) }}
                  </span>
                </div>
              </div>
              <div class="email-status">
                <div v-if="!email.isRead" class="unread-dot" :class="email.priority"></div>
                <button class="delete-btn" @click.stop="deleteEmail(email.id)">
                  <i class="fa fa-trash-o"></i>
                </button>
              </div>
            </div>
          </transition-group>

          <div v-if="filteredEmails.length === 0" class="empty-state">
            <i class="fa fa-envelope-o"></i>
            <p>{{ emptyStateText }}</p>
            <button v-if="activeTab === 'unread' && headerStore.emailsByStatus['inbox']?.length > 0"
                    class="btn-secondary" @click="activeTab = 'inbox'">
              View All Emails
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>

  <!-- Email Detail Modal -->
  <transition name="modal">
    <div v-if="currentEmail" class="email-detail-modal">
      <div class="modal-overlay" @click="closeDetailModal"></div>
      <div class="modal-content" :class="currentEmail.priority">
        <div class="modal-header">
          <div class="header-info">
            <div class="modal-icon">
              <i class="fa fa-envelope"></i>
            </div>
            <div class="modal-title-section">
              <h4>{{ currentEmail.subject }}</h4>
              <div class="modal-meta">
                <span class="detail-from">From: {{ currentEmail.fromName }}</span>
                <span class="detail-time">{{ formatTime(currentEmail.createdAt) }}</span>
                <span class="detail-priority" :class="currentEmail.priority">
                  {{ getPriorityText(currentEmail.priority) }}
                </span>
              </div>
            </div>
          </div>
          <button class="btn-modal-close" @click="closeDetailModal">
            <i class="fa fa-times"></i>
          </button>
        </div>
        <div class="modal-body">
          <div class="email-detail-info">
            <div class="email-detail-row">
              <span class="detail-label">From:</span>
              <span class="detail-value">{{ currentEmail.fromName }} &lt;{{ currentEmail.from }}&gt;</span>
            </div>
            <div class="email-detail-row">
              <span class="detail-label">To:</span>
              <span class="detail-value">{{ currentEmail.to }}</span>
            </div>
            <div class="email-detail-row">
              <span class="detail-label">Date:</span>
              <span class="detail-value">{{ formatDate(currentEmail.createdAt) }}</span>
            </div>
            <div v-if="currentEmail.attachments && currentEmail.attachments.length > 0" class="email-detail-row">
              <span class="detail-label">Attachments:</span>
              <div class="attachments">
                <span v-for="(attachment, index) in currentEmail.attachments" :key="index" class="attachment-item">
                  <i class="fa fa-file-o"></i>
                  {{ attachment }}
                </span>
              </div>
            </div>
          </div>
          <div class="detail-content">
            <div class="email-body" v-html="formatEmailContent(currentEmail.content)"></div>
          </div>
          <div class="detail-actions">
            <button class="action-btn reply-btn" @click="replyEmail">
              <i class="fa fa-reply"></i>
              Reply
            </button>
            <button class="action-btn forward-btn" @click="forwardEmail">
              <i class="fa fa-share"></i>
              Forward
            </button>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="closeDetailModal">
            <i class="fa fa-times"></i>
            Close
          </button>
          <button v-if="!currentEmail.isRead" class="btn-secondary" @click="markAsReadAndClose">
            <i class="fa fa-check"></i>
            Mark as Read
          </button>
          <button class="btn-secondary delete-email-btn" @click="deleteEmail(currentEmail.id)">
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
import {useHeaderStore} from '../../../data-sources/dashboard/header/header.js'

const headerStore = useHeaderStore()
const activeTab = ref('inbox')
const currentEmail = ref(null)

const filteredEmails = computed(() => {
  if (activeTab.value === 'unread') {
    return headerStore.unreadEmails
  } else if (activeTab.value === 'sent') {
    return headerStore.emailsByStatus['sent'] || []
  } else {
    return headerStore.emailsByStatus['inbox'] || []
  }
})

const emptyStateText = computed(() => {
  if (activeTab.value === 'unread') {
    return 'No unread emails'
  } else if (activeTab.value === 'sent') {
    return 'No sent emails'
  } else {
    return 'Inbox is empty'
  }
})

const getPriorityText = (priority) => {
  const priorityMap = {
    'high': 'High',
    'medium': 'Normal',
    'normal': 'Normal',
    'low': 'Low Priority'
  }
  return priorityMap[priority] || 'Normal'
}

const formatEmailContent = (content) => {
  // Convert text content to HTML format, preserving line breaks
  if (!content) return ''
  return content.replace(/\n/g, '<br>')
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
      day: 'numeric'
    })
  }
}

const formatDate = (timestamp) => {
  const time = new Date(timestamp)
  return time.toLocaleDateString('en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const showEmailDetail = (email) => {
  currentEmail.value = email
  if (!email.isRead) {
    headerStore.markEmailAsRead(email.id)
  }
}

const closeDetailModal = () => {
  currentEmail.value = null
}

const markAsReadAndClose = () => {
  if (currentEmail.value) {
    headerStore.markEmailAsRead(currentEmail.value.id)
    currentEmail.value = null
  }
}

const deleteEmail = (id) => {
  if (confirm('Are you sure you want to delete this email?')) {
    headerStore.deleteEmail(id)
    if (currentEmail.value && currentEmail.value.id === id) {
      currentEmail.value = null
    }
  }
}

const markAllAsRead = () => {
  if (confirm('Are you sure you want to mark all emails as read?')) {
    headerStore.markAllEmailsAsRead()
  }
}

const composeNewEmail = () => {
  // Compose email functionality can be implemented here, temporarily using alert
  alert('Compose email feature to be implemented')
}

const replyEmail = () => {
  if (currentEmail.value) {
    // Reply functionality can be implemented here, temporarily using alert
    alert(`Reply to: ${currentEmail.value.fromName}`)
  }
}

const forwardEmail = () => {
  if (currentEmail.value) {
    // Forward functionality can be implemented here, temporarily using alert
    alert(`Forward email: ${currentEmail.value.subject}`)
  }
}

const refreshEmails = () => {
  // notificationStore没有refreshEmails方法，暂时移除
}
</script>

<style scoped>
.email-panel-container {
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

.modal-enter-active .email-panel,
.modal-leave-active .email-panel {
  transition: opacity 0.3s ease;
}

.modal-enter-from .email-panel,
.modal-leave-to .email-panel {
  opacity: 0;
}

.email-panel {
  position: fixed;
  top: 60px;
  right: 20px;
  width: 450px;
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

.email-list {
  flex: 1;
  overflow-y: auto;
  max-height: calc(92vh - 200px);
}

.email-item {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  align-items: center;
}

.email-checkbox {
  width: 20px;
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
  padding-top: 4px;
}

.email-item:hover {
  background: #f8fafc;
}

.email-item.unread {
  background: #f0f9ff;
  border-left: 3px solid var(--color-primary);
}

.email-sender {
  width: 80px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.email-sender i {
  font-size: 24px;
  color: #64748b;
}

.sender-name {
  font-size: 12px;
  font-weight: 500;
  color: #475569;
  text-align: center;
}

.has-attachment i {
  color: #f59e0b;
}

.email-content {
  flex: 1;
  min-width: 0;
}

.email-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.email-subject {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-time {
  font-size: 12px;
  color: #64748b;
  flex-shrink: 0;
}

.email-preview {
  font-size: 13px;
  color: #64748b;
  line-height: 1.4;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.email-meta {
  display: flex;
  gap: 8px;
  align-items: center;
}

.attachment-icon {
  font-size: 12px;
  color: #f59e0b;
  display: flex;
  align-items: center;
  gap: 2px;
}

.email-priority {
  font-size: 11px;
  padding: 2px 6px;
  font-weight: 500;
}

.email-priority.high {
  background: #fef2f2;
  color: #dc2626;
}

.email-priority.medium,
.email-priority.normal {
  background: #fffbeb;
  color: #d97706;
}

.email-priority.low {
  background: #f0f9ff;
  color: #0284c7;
}

.email-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.unread-dot {
  width: 8px;
  height: 8px;
  background: var(--color-primary);
  border-radius: 50%;
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

.email-item:hover .delete-btn {
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


/* 邮件详情模态框样式 */
.email-detail-modal {
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
  width: 700px;
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
  background: linear-gradient(135deg, var(--color-primary), #6366f1);
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

.detail-from,
.detail-time {
  color: #666;
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

.detail-priority.medium,
.detail-priority.normal {
  background: linear-gradient(135deg, #faad14, #ffc53d);
  color: white;
}

.detail-priority.low {
  background: linear-gradient(135deg, #52c41a, #73d13d);
  color: white;
}

.modal-body {
  padding: 24px;
  flex: 1;
  overflow-y: auto;
  background: var(--background-primary);
}

.email-detail-info {
  background: white;
  padding: 20px;
  margin-bottom: 24px;
  border-radius: 4px;
}

.email-detail-row {
  display: flex;
  margin-bottom: 12px;
}

.email-detail-row:last-child {
  margin-bottom: 0;
}

.detail-label {
  width: 80px;
  font-weight: 600;
  color: #666;
  flex-shrink: 0;
}

.detail-value {
  flex: 1;
  color: #333;
}

.attachments {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  flex: 1;
}

.attachment-item {
  background: #f1f5f9;
  padding: 6px 12px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #475569;
}

.detail-content {
  font-size: 15px;
  line-height: 1.7;
  color: #333;
  margin-bottom: 24px;
  background: white;
  padding: 20px;
  border-radius: 4px;
}

.email-body {
  white-space: pre-wrap;
  word-break: break-word;
}

.detail-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.reply-btn,
.forward-btn {
  background: var(--color-primary);
  color: white;
  border: none;
}

.reply-btn:hover,
.forward-btn:hover {
  background: #2563eb;
  color: white;
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
  border-color: var(--color-primary);
  color: var(--color-primary);
  transform: translateY(-1px);
}


.delete-email-btn:hover {
  border-color: var(--border-color);
  color: #dc2626;
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  opacity: 0;
  transform: translateY(-50px) scale(0.9);
}
</style>
