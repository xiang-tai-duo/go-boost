// --------------------------------------------------------------------------------
// File:        settings.vue
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: System settings panel component for managing user preferences and system configuration.
// --------------------------------------------------------------------------------
<template>
  <transition name="modal">
    <div v-if="headerStore.showSettingsPanel" class="modal-overlay" @click.self="headerStore.closeAllPanels()">
      <div class="settings-panel shadow-settings">
        <div class="modal-header">
          <h2 class="modal-title">System Settings</h2>
          <button class="btn-modal-close" @click="headerStore.closeAllPanels()">
            <i class="close-icon"></i>
          </button>
        </div>

        <div class="settings-container">
          <!-- Left Navigation -->
          <div class="settings-sidebar">
            <div class="sidebar-section">
              <h4>Personal Settings</h4>
              <ul>
                <li
                    v-for="item in personalSettings"
                    :key="item.key"
                    :class="{ active: activeSection === item.key }"
                    @click="activeSection = item.key"
                >
                  <i :class="item.icon"></i>
                  {{ item.title }}
                </li>
              </ul>
            </div>

            <div class="sidebar-section">
              <h4>System Settings</h4>
              <ul>
                <li
                    v-for="item in systemSettings"
                    :key="item.key"
                    :class="{ active: activeSection === item.key }"
                    @click="activeSection = item.key"
                >
                  <i :class="item.icon"></i>
                  {{ item.title }}
                </li>
              </ul>
            </div>

            <div class="sidebar-section">
              <h4>Advanced Settings</h4>
              <ul>
                <li
                    v-for="item in advancedSettings"
                    :key="item.key"
                    :class="{ active: activeSection === item.key }"
                    @click="activeSection = item.key"
                >
                  <i :class="item.icon"></i>
                  {{ item.title }}
                </li>
              </ul>
            </div>
          </div>

          <!-- Right Content -->
          <div class="settings-content">
            <div class="content-header">
              <h2>{{ currentSection.title }}</h2>
              <p class="section-description">{{ currentSection.description }}</p>
            </div>

            <div class="settings-items">
              <div
                  v-for="setting in currentSection.settings"
                  :key="setting.key"
                  class="setting-item"
              >
                <div class="setting-header">
                  <h4>{{ setting.title }}</h4>
                  <p class="setting-description">{{ setting.description }}</p>
                </div>

                <div class="setting-control">
                  <!-- 输入框 -->
                  <input
                      v-if="['text', 'email', 'tel', 'password', 'number'].includes(setting.type)"
                      :type="setting.type"
                      :value="getSettingValue(currentSection.key, setting.key)"
                      @input="saveSetting(currentSection.key, setting.key, $event.target.value)"
                      class="form-input input-primary"
                  />

                  <!-- 复选框 -->
                  <label v-else-if="setting.type === 'checkbox'" class="checkbox-label">
                    <input
                        type="checkbox"
                        :checked="getSettingValue(currentSection.key, setting.key)"
                        @change="saveSetting(currentSection.key, setting.key, $event.target.checked)"
                        class="checkbox-primary"
                    />
                    <span class="checkmark"></span>
                  </label>

                  <!-- 下拉选择 -->
                  <select
                      v-else-if="setting.type === 'select'"
                      :value="getSettingValue(currentSection.key, setting.key)"
                      @change="saveSetting(currentSection.key, setting.key, $event.target.value)"
                      class="select-primary"
                  >
                    <option v-for="option in setting.options" :key="option" :value="option">
                      {{ option }}
                    </option>
                  </select>

                  <!-- 文本域 -->
                  <textarea
                      v-else-if="setting.type === 'textarea'"
                      :value="getSettingValue(currentSection.key, setting.key)"
                      @input="saveSetting(currentSection.key, setting.key, $event.target.value)"
                      class="textarea-primary"
                  ></textarea>

                  <!-- 进度条 -->
                  <div v-else-if="setting.type === 'progress'" class="progress-container">
                    <progress
                        :value="getSettingValue(currentSection.key, setting.key)"
                        max="100"
                        class="form-progress"
                    ></progress>
                    <span class="progress-text">{{ getSettingValue(currentSection.key, setting.key) }}%</span>
                  </div>

                  <!-- 时间范围 -->
                  <div v-else-if="setting.type === 'time-range'" class="time-range-control">
                    <input
                        type="time"
                        :value="getSettingValue(currentSection.key, setting.key)?.start || '22:00'"
                        @input="saveSetting(currentSection.key, setting.key, {
                        ...getSettingValue(currentSection.key, setting.key),
                        start: $event.target.value
                      })"
                        class="form-time input-primary"
                    />
                    <span class="time-separator"> to </span>
                    <input
                        type="time"
                        :value="getSettingValue(currentSection.key, setting.key)?.end || '08:00'"
                        @input="saveSetting(currentSection.key, setting.key, {
                        ...getSettingValue(currentSection.key, setting.key),
                        end: $event.target.value
                      })"
                        class="form-time input-primary"
                    />
                  </div>

                  <!-- 只读信息 - 新风格 -->
                  <div v-else-if="setting.type === 'readonly'" class="info-item">
                    <div class="info-content">
                      <div class="info-title">{{ setting.title }}</div>
                      <div class="info-value">{{ setting.value }}</div>
                    </div>
                  </div>

                  <!-- 许可证信息 -->
                  <div v-else-if="setting.type === 'license'" class="license-list">
                    <div v-for="license in setting.value" :key="license.name" class="license-item">
                      <div class="license-info">
                        <div class="license-name">{{ license.name }} {{ license.version }}</div>
                        <div class="license-version">{{ license.license }}</div>
                      </div>
                      <a :href="license.url" target="_blank" class="license-link">View License</a>
                    </div>
                  </div>

                  <!-- 支持信息 -->
                  <div v-else-if="setting.type === 'support'" class="support-info">
                    <div class="support-item">
                      <span class="support-label">Official Website</span>
                      <span class="support-value">
                        <a :href="setting.value.website" target="_blank">{{ setting.value.website }}</a>
                      </span>
                    </div>
                    <div class="support-item">
                      <span class="support-label">Support Email</span>
                      <span class="support-value">
                        <a :href="`mailto:${setting.value.email}`">{{ setting.value.email }}</a>
                      </span>
                    </div>
                    <div class="support-item">
                      <span class="support-label">Support Phone</span>
                      <span class="support-value">{{ setting.value.phone }}</span>
                    </div>
                    <div class="support-item">
                      <span class="support-label">Company Address</span>
                      <span class="support-value">{{ setting.value.address }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import {computed, onMounted, onUnmounted, ref} from 'vue'
import {useHeaderStore} from '../../../data-sources/dashboard/header/header.js'

const headerStore = useHeaderStore()
const activeSection = ref('basicInfo')

const personalSettings = [
  {
    key: 'basicInfo',
    title: 'Basic Info',
    icon: 'fa fa-user',
    description: 'Manage your personal profile and account information',
    settings: [
      {
        key: 'username',
        title: 'Username',
        description: 'Your login username, used for system login',
        type: 'text',
        value: 'admin'
      },
      {
        key: 'displayName',
        title: 'Display Name',
        description: 'Name displayed in the system',
        type: 'text',
        value: 'Admin'
      },
      {
        key: 'email',
        title: 'Email Address',
        description: 'Used for receiving system notifications and emails',
        type: 'email',
        value: 'admin@company.com'
      },
      {
        key: 'phone',
        title: 'Phone Number',
        description: 'Used for receiving SMS notifications',
        type: 'tel',
        value: '13800138000'
      },
      {
        key: 'department',
        title: 'Department',
        description: 'Your department information',
        type: 'select',
        options: ['IT Department', 'Human Resources', 'Finance', 'Marketing', 'Sales'],
        value: 'IT Department'
      },
      {
        key: 'position',
        title: 'Position',
        description: 'Your position in the company',
        type: 'text',
        value: 'System Administrator'
      }
    ]
  },
  {
    key: 'avatar',
    title: 'Avatar Settings',
    icon: 'fa fa-camera',
    description: 'Upload and manage your personal avatar',
    settings: [
      {
        key: 'avatarImage',
        title: 'Avatar Image',
        description: 'Supports JPG, PNG formats, recommended size 200x200 pixels',
        type: 'file',
        value: null
      },
      {
        key: 'avatarStyle',
        title: 'Avatar Style',
        description: 'Select the display style of your avatar',
        type: 'select',
        options: ['Circle', 'Square', 'Rounded'],
        value: 'Circle'
      },
      {
        key: 'useGravatar',
        title: 'Use Gravatar',
        description: 'Use Gravatar service for your avatar',
        type: 'checkbox',
        value: false
      }
    ]
  },
  {
    key: 'password',
    title: 'Password Security',
    icon: 'fa fa-lock',
    description: 'Change password and security settings',
    settings: [
      {
        key: 'currentPassword',
        title: 'Current Password',
        description: 'Enter your current password to verify identity',
        type: 'password',
        value: ''
      },
      {
        key: 'newPassword',
        title: 'New Password',
        description: 'Set a new password, at least 8 characters',
        type: 'password',
        value: ''
      },
      {
        key: 'confirmPassword',
        title: 'Confirm New Password',
        description: 'Enter new password again',
        type: 'password',
        value: ''
      },
      {
        key: 'passwordStrength',
        title: 'Password Strength',
        description: 'Shows password strength level',
        type: 'progress',
        value: 75
      }
    ]
  },
  {
    key: 'privacy',
    title: 'Privacy Settings',
    icon: 'fa fa-eye-slash',
    description: 'Control the visibility of your personal information',
    settings: [
      {
        key: 'profileVisibility',
        title: 'Profile Visibility',
        description: 'Control whether other users can see your profile',
        type: 'select',
        options: ['Public', 'Colleagues Only', 'Admins Only', 'Private'],
        value: 'Colleagues Only'
      },
      {
        key: 'showEmail',
        title: 'Show Email Address',
        description: 'Display email address in your profile',
        type: 'checkbox',
        value: true
      },
      {
        key: 'showPhone',
        title: 'Show Phone Number',
        description: 'Display phone number in your profile',
        type: 'checkbox',
        value: false
      },
      {
        key: 'allowMessages',
        title: 'Allow Messages',
        description: 'Allow other users to send you messages',
        type: 'checkbox',
        value: true
      }
    ]
  }
]

const systemSettings = [
  {
    key: 'general',
    title: 'General Settings',
    icon: 'fa fa-cog',
    description: 'Basic system configuration and preferences',
    settings: [
      {
        key: 'language',
        title: 'System Language',
        description: 'Select system display language',
        type: 'select',
        options: ['简体中文', '繁體中文', 'English', '日本語'],
        value: 'English'
      },
      {
        key: 'timezone',
        title: 'Timezone',
        description: 'Select system timezone',
        type: 'select',
        options: ['UTC+8 Beijing Time', 'UTC+9 Tokyo Time', 'UTC-5 New York Time', 'UTC+0 London Time'],
        value: 'UTC+8 Beijing Time'
      },
      {
        key: 'dateFormat',
        title: 'Date Format',
        description: 'Select date display format',
        type: 'select',
        options: ['YYYY-MM-DD', 'MM/DD/YYYY', 'DD/MM/YYYY', 'YYYY年MM月DD日'],
        value: 'YYYY-MM-DD'
      },
      {
        key: 'timeFormat',
        title: 'Time Format',
        description: 'Select time display format',
        type: 'select',
        options: ['24-hour format', '12-hour format'],
        value: '24-hour format'
      },
      {
        key: 'theme',
        title: 'Interface Theme',
        description: 'Select system interface theme',
        type: 'select',
        options: ['Light Theme', 'Dark Theme', 'Follow System'],
        value: 'Light Theme'
      },
      {
        key: 'animations',
        title: 'Animations',
        description: 'Enable or disable interface animations',
        type: 'checkbox',
        value: true
      }
    ]
  },
  {
    key: 'notification',
    title: 'Notification Settings',
    icon: 'fa fa-bell',
    description: 'Configure system notifications and reminders',
    settings: [
      {
        key: 'emailNotifications',
        title: 'Email Notifications',
        description: 'Receive important notifications via email',
        type: 'checkbox',
        value: true
      },
      {
        key: 'smsNotifications',
        title: 'SMS Notifications',
        description: 'Receive emergency notifications via SMS',
        type: 'checkbox',
        value: false
      },
      {
        key: 'browserNotifications',
        title: 'Browser Notifications',
        description: 'Show notifications in browser',
        type: 'checkbox',
        value: true
      },
      {
        key: 'notificationSound',
        title: 'Notification Sound',
        description: 'Play sound when receiving notifications',
        type: 'checkbox',
        value: true
      },
      {
        key: 'quietHours',
        title: 'Quiet Hours',
        description: 'Set quiet hours',
        type: 'time-range',
        value: {start: '22:00', end: '08:00'}
      }
    ]
  },
  {
    key: 'print',
    title: 'Print Settings',
    icon: 'fa fa-print',
    description: 'Configure printing related parameters and options',
    settings: [
      {
        key: 'defaultPrinter',
        title: 'Default Printer',
        description: 'Set default printer',
        type: 'select',
        options: ['Laser Printer 1', 'Laser Printer 2', 'Laser Printer 3', 'Auto Select'],
        value: 'Laser Printer 1'
      },
      {
        key: 'defaultPrintMode',
        title: 'Default Print Mode',
        description: 'Set default print mode',
        type: 'select',
        options: ['Single-sided', 'Double-sided', 'Economy Mode', 'High Quality Mode'],
        value: 'Double-sided'
      },
      {
        key: 'printQuota',
        title: 'Print Quota',
        description: 'Set monthly print quota (pages)',
        type: 'number',
        min: 100,
        max: 10000,
        value: 1000
      },
      {
        key: 'printPreview',
        title: 'Print Preview',
        description: 'Show preview window before printing',
        type: 'checkbox',
        value: true
      },
      {
        key: 'printLogs',
        title: 'Print Logs',
        description: 'Record print operation logs',
        type: 'checkbox',
        value: true
      }
    ]
  },
  {
    key: 'security',
    title: 'Security Settings',
    icon: 'fa fa-shield',
    description: 'System security related configuration',
    settings: [
      {
        key: 'twoFactorAuth',
        title: 'Two-Factor Authentication',
        description: 'Enable two-factor authentication to improve account security',
        type: 'checkbox',
        value: false
      },
      {
        key: 'sessionTimeout',
        title: 'Session Timeout',
        description: 'Set session timeout time (minutes)',
        type: 'number',
        min: 5,
        max: 480,
        value: 30
      },
      {
        key: 'ipWhitelist',
        title: 'IP Whitelist',
        description: 'List of allowed IP addresses, separated by commas',
        type: 'textarea',
        value: ''
      },
      {
        key: 'loginAttempts',
        title: 'Login Attempts',
        description: 'Maximum number of allowed login attempts',
        type: 'number',
        min: 3,
        max: 10,
        value: 5
      }
    ]
  }
]

const advancedSettings = [
  {
    key: 'data',
    title: 'Data Management',
    icon: 'fa fa-database',
    description: 'Data backup, restore and cleanup settings',
    settings: [
      {
        key: 'autoBackup',
        title: 'Auto Backup',
        description: 'Enable automatic data backup functionality',
        type: 'checkbox',
        value: true
      },
      {
        key: 'backupFrequency',
        title: 'Backup Frequency',
        description: 'Set data backup frequency',
        type: 'select',
        options: ['Daily', 'Weekly', 'Monthly'],
        value: 'Daily'
      },
      {
        key: 'dataRetention',
        title: 'Data Retention Period',
        description: 'Set data retention period (days)',
        type: 'number',
        min: 7,
        max: 3650,
        value: 365
      },
      {
        key: 'exportFormat',
        title: 'Export Format',
        description: 'Set default data export format',
        type: 'select',
        options: ['Excel', 'CSV', 'JSON', 'XML'],
        value: 'Excel'
      }
    ]
  },
  {
    key: 'api',
    title: 'API Settings',
    icon: 'fa fa-code',
    description: 'API interface and developer related settings',
    settings: [
      {
        key: 'enableAPI',
        title: 'Enable API',
        description: 'Enable system API interface',
        type: 'checkbox',
        value: true
      },
      {
        key: 'apiKey',
        title: 'API Key',
        description: 'API access key',
        type: 'password',
        value: 'sk-1234567890abcdef'
      },
      {
        key: 'rateLimit',
        title: 'Rate Limit',
        description: 'API request rate limit (times/hour)',
        type: 'number',
        min: 100,
        max: 10000,
        value: 1000
      },
      {
        key: 'corsOrigin',
        title: 'CORS Origin',
        description: 'Allowed cross-origin request sources',
        type: 'text',
        value: '*'
      }
    ]
  },
  {
    key: 'logs',
    title: 'Log Settings',
    icon: 'fa fa-file-text',
    description: 'System log and error report configuration',
    settings: [
      {
        key: 'logLevel',
        title: 'Log Level',
        description: 'Set log recording level',
        type: 'select',
        options: ['Debug', 'Info', 'Warning', 'Error'],
        value: 'Info'
      },
      {
        key: 'logRetention',
        title: 'Log Retention Period',
        description: 'Set log file retention period (days)',
        type: 'number',
        min: 1,
        max: 365,
        value: 30
      },
      {
        key: 'errorAlerts',
        title: 'Error Alerts',
        description: 'Send alert notifications when errors occur',
        type: 'checkbox',
        value: true
      }
    ]
  },
  {
    key: 'about',
    title: 'About System',
    icon: 'fa fa-info-circle',
    description: 'System information and version details',
    settings: [
      {
        key: 'version',
        title: 'System Version',
        description: 'Current system version number',
        type: 'readonly',
        value: 'v2.1.0'
      },
      {
        key: 'buildDate',
        title: 'Build Date',
        description: 'System build date',
        type: 'readonly',
        value: '2024-01-15'
      },
      {
        key: 'licenseInfo',
        title: 'License Information',
        description: 'Open source component license information',
        type: 'license',
        value: [
          {
            name: 'Vue.js',
            version: '3.3.4',
            license: 'MIT License',
            url: 'https://github.com/vuejs/vue/blob/main/LICENSE'
          },
          {
            name: 'Pinia',
            version: '2.1.7',
            license: 'MIT License',
            url: 'https://github.com/vuejs/pinia/blob/v2/LICENSE'
          },
          {
            name: 'Chart.js',
            version: '4.4.0',
            license: 'MIT License',
            url: 'https://github.com/chartjs/Chart.js/blob/master/LICENSE.md'
          },
          {
            name: 'Font Awesome',
            version: '6.4.0',
            license: 'Font Awesome Free License',
            url: 'https://github.com/FortAwesome/Font-Awesome/blob/6.x/LICENSE.txt'
          },
          {
            name: 'ECharts',
            version: '5.4.3',
            license: 'Apache License 2.0',
            url: 'https://github.com/apache/echarts/blob/master/LICENSE'
          }
        ]
      },
      {
        key: 'thirdParty',
        title: 'Third-Party Components',
        description: 'Third-party open source components used',
        type: 'readonly',
        value: 'Vue.js, Pinia, Chart.js, Font Awesome, ECharts, Element Plus'
      },
      {
        key: 'support',
        title: 'Technical Support',
        description: 'Ways to get technical support',
        type: 'support',
        value: {
          website: 'https://www.traecn.com',
          email: 'support@traecn.com',
          phone: '400-123-4567',
          address: 'Room 1201, TRAE Building, High Tech Zone, Beijing, China'
        }
      }
    ]
  }
]

const currentSection = computed(() => {
  const allSettings = [...personalSettings, ...systemSettings, ...advancedSettings]
  return allSettings.find(section => section.key === activeSection.value) || personalSettings[0]
})

const getSettingValue = (category, key) => {
  return headerStore.getSetting(category, key)
}

const saveSetting = (category, key, value) => {
  headerStore.saveSetting(category, key, value)
  console.log(`Save setting: ${category}.${key} = ${value}`)
}

const saveSettings = () => {
  const settings = currentSection.value.settings.reduce((acc, setting) => {
    acc[setting.key] = getSettingValue(currentSection.value.key, setting.key)
    return acc
  }, {})

  headerStore.saveAllSettings(currentSection.value.key, settings)
  console.log('Save all settings')
  alert('Settings saved successfully!')
}

const resetSettings = () => {
  if (confirm('Are you sure you want to reset all settings for the current category?')) {
    headerStore.resetSettings(currentSection.value.key)
    console.log('Reset settings')
    alert('Settings have been reset to default values!')
  }
}

// Listen for ESC key to close settings panel
const handleEscKey = (event) => {
  if (event.key === 'Escape' && headerStore.showSettingsPanel) {
    headerStore.closeAllPanels()
  }
}

// Add event listener when component mounts
onMounted(() => {
  document.addEventListener('keydown', handleEscKey)
})

// Remove event listener when component unmounts to avoid memory leaks
onUnmounted(() => {
  document.removeEventListener('keydown', handleEscKey)
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .settings-panel,
.modal-leave-active .settings-panel {
  transition: opacity 0.3s ease;
}

.modal-enter-from .settings-panel,
.modal-leave-to .settings-panel {
  opacity: 0;
}

.settings-panel {
  width: 90vw;
  height: 90vh;
  background: var(--background-modal);
  border: 1px solid var(--border-color);
  z-index: 1000;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  overflow: hidden;
}

.modal-title {
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  display: flex;
  align-items: center;
}

.modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: #64748b;
  padding: 2px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-close:hover {
  background-color: #f1f5f9;
  color: #334155;
}

.close-icon {
  width: 20px;
  height: 20px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='currentColor'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M6 18L18 6M6 6l12 12' /%3E%3C/svg%3E");
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

.settings-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.settings-sidebar {
  width: 280px;
  background: var(--background-modal);
  border-right: 1px solid var(--border-color);
  overflow-y: auto;
  padding: 20px 0;
}

.sidebar-section {
  margin-bottom: 24px;
}

.sidebar-section h4 {
  margin: 0 0 12px 0;
  padding: 0 20px;
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.sidebar-section ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.sidebar-section li {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #374151;
  border-left: 3px solid transparent;
}

.sidebar-section li:hover {
  background: #e2e8f0;
  color: #1e293b;
}

.sidebar-section li.active {
  background: #dbeafe;
  color: #1d4ed8;
  border-left-color: var(--color-primary);
  font-weight: 500;
}

.sidebar-section li i {
  width: 16px;
  text-align: center;
}

.settings-content {
  flex: 1;
  padding: 24px;
  overflow-y: visible; /* 移除滚动条 */
}

.content-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-color);
}

.content-header h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
}

.section-description {
  margin: 0;
  font-size: 14px;
  color: #64748b;
}

.settings-items {
  display: flex;
  flex-direction: column;
  padding: 8px;
  background-color: var(--background-modal);
  overflow-y: auto;
  max-height: calc(92vh - 200px); /* 减去头部和间距的高度 */
}

.setting-item {
  display: flex;
  margin-bottom: 4px;
}

.setting-header {
  flex: 1;
}

.setting-header h4 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.setting-description {
  font-size: 12px;
  color: #64748b;
}

.setting-control {
  display: flex;
  align-items: center;
  width: 240px;
}

.form-input {
  width: 100%;
  height: 30px;
  border: 1px solid var(--border-color);
  font-size: 14px;
  padding: 0px 12px;
  transition: border-color 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
}

.form-input.readonly {
  background: #f1f5f9;
  color: #64748b;
  cursor: not-allowed;
}

.form-progress {
  width: 100%;
  height: 8px;
  border: none;
  background: #e2e8f0;
}

.form-progress::-webkit-progress-bar {
  background: #e2e8f0;
}

.form-progress::-webkit-progress-value {
  background: var(--color-primary);
}

.progress-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.progress-text {
  font-size: 12px;
  color: #64748b;
  min-width: 35px;
}

.time-range-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-separator {
  font-size: 14px;
  color: #64748b;
}

.form-time {
  padding: 6px 10px;
  border: 1px solid var(--border-color);
  font-size: 14px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 14px;
}

.checkbox-label input[type="checkbox"] {
  margin: 0;
  width: 16px;
  height: 16px;
  cursor: pointer;
  border-radius: 0 !important;
  -webkit-appearance: none !important;
  -moz-appearance: none !important;
  appearance: none !important;
  border: 1px solid var(--border-color) !important;
  background-color: white !important;
  position: relative !important;
  box-sizing: border-box !important;
}

.checkbox-label input[type="checkbox"]:checked {
  background-color: var(--color-primary) !important;
  border-color: var(--color-primary) !important;
}

.checkbox-label input[type="checkbox"]:checked::after {
  content: "" !important;
  position: absolute !important;
  top: 44% !important;
  left: 50% !important;
  width: 3px !important;
  height: 6px !important;
  border: solid white !important;
  border-width: 0 1.5px 1.5px 0 !important;
  transform: translate(-50%, -50%) rotate(45deg) !important;
}

.license-list {
  display: flex;
  flex-direction: column;
}

.license-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f1f5f9;
  font-size: 12px;
}

.license-item .license-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.license-item .license-name {
  font-weight: 600;
  color: #374151;
}

.license-item .license-version {
  color: #64748b;
}

.license-item .license-link {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 500;
}

.license-item .license-link:hover {
  text-decoration: underline;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f1f5f9;
  border-radius: 6px;
}

.info-content {
  display: flex;
  flex-direction: column;
}

.info-title {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #374151;
  font-weight: 600;
}

.support-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.support-item {
  display: flex;
  align-items: center;
}

.support-item:last-child {
  border-bottom: none;
}

.support-label {
  font-weight: 600;
  color: #374151;
  font-size: 13px;
  width: 100px;
}

.support-value {
  color: #64748b;
  font-size: 13px;
}

.support-value a {
  color: var(--color-primary);
  text-decoration: none;
}

.support-value a:hover {
  text-decoration: underline;
}

.btn-secondary {
  padding: 10px 20px;
  background: #f1f5f9;
  color: #475569;
  border: 1px solid var(--border-color);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #e2e8f0;
  color: #1e293b;
}


@media (max-width: 1024px) {
  .settings-panel {
    width: 95vw;
    height: 85vh;
  }

  .settings-sidebar {
    width: 240px;
  }

  .settings-content {
    padding: 16px;
  }
}

@media (max-width: 768px) {
  .settings-panel {
    width: 100vw;
    height: 100vh;
    border-radius: 0;
  }

  .settings-container {
    flex-direction: column;
  }

  .settings-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
    max-height: 200px;
  }

  .setting-item {
    flex-direction: column;
    gap: 12px;
    align-items: center;
  }

  .setting-control {
    width: 100%;
  }
}
</style>
