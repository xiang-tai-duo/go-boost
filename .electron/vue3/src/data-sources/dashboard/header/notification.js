// --------------------------------------------------------------------------------
// File:        notification.js
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Pinia store for managing notifications.
// --------------------------------------------------------------------------------
import {defineStore} from 'pinia'

export const useNotificationStore = defineStore('notification', {
    state: () => ({
        notifications: [
            {
                id: 1,
                type: 'system',
                title: 'System Update Notification',
                content: 'The system will undergo routine maintenance tonight from 22:00-23:00, which may affect normal usage. Please prepare in advance',
                isRead: false,
                priority: 'high',
                createdAt: new Date(Date.now() - 1000 * 60 * 30),
                category: 'System Maintenance'
            },
            {
                id: 2,
                type: 'print',
                title: 'Print Job Completed',
                content: 'Your document "Q3 Quarterly Report.pdf" has been printed on a Laser Printer, total 12 pages',
                isRead: false,
                priority: 'normal',
                createdAt: new Date(Date.now() - 1000 * 60 * 15),
                category: 'Print Jobs'
            },
            {
                id: 3,
                type: 'device',
                title: 'Device Status Warning',
                content: 'Printer ink cartridge level is below 20%, please replace consumables in time',
                isRead: true,
                priority: 'medium',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2),
                category: 'Device Status'
            },
            {
                id: 4,
                type: 'quota',
                title: 'Print Quota Reminder',
                content: 'You have used 85% of your monthly print quota, 234 pages remaining. Please use it reasonably',
                isRead: false,
                priority: 'medium',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 4),
                category: 'Quota Management'
            },
            {
                id: 5,
                type: 'security',
                title: 'Security Alert',
                content: 'Abnormal login behavior detected, IP address 192.168.1.100 attempted to access the system. Please confirm if this is your operation',
                isRead: true,
                priority: 'high',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24),
                category: 'Security Alerts'
            },
            {
                id: 6,
                type: 'system',
                title: 'New Feature Launch',
                content: 'Batch printing function has been added to the print management system, supporting one-time processing of multiple documents. Welcome to experience',
                isRead: false,
                priority: 'normal',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 6),
                category: 'Feature Updates'
            },
            {
                id: 7,
                type: 'print',
                title: 'Print Job Failed',
                content: 'Your document "Contract Template.docx" printing failed. Error reason: Paper size mismatch. Please check printer settings',
                isRead: false,
                priority: 'high',
                createdAt: new Date(Date.now() - 1000 * 60 * 45),
                category: 'Print Jobs'
            },
            {
                id: 8,
                type: 'device',
                title: 'Device Offline Notification',
                content: 'A printer has been offline for more than 30 minutes. Please check network connection and device status',
                isRead: true,
                priority: 'medium',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 3),
                category: 'Device Status'
            },
            {
                id: 9,
                type: 'quota',
                title: 'Quota About to Run Out',
                content: 'Your monthly print quota has only 56 pages left. It is recommended to apply for an increase or optimize print settings',
                isRead: false,
                priority: 'high',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 8),
                category: 'Quota Management'
            },
            {
                id: 10,
                type: 'security',
                title: 'Password About to Expire',
                content: 'Your system password will expire in 7 days. Please change your password in time to ensure account security',
                isRead: true,
                priority: 'medium',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 12),
                category: 'Security Alerts'
            },
            {
                id: 11,
                type: 'system',
                title: 'Data Backup Completed',
                content: 'Today\'s data backup task has been successfully completed, all important data has been safely saved to the cloud',
                isRead: true,
                priority: 'low',
                createdAt: new Date(Date.now() - 1000 * 60 * 60 * 1),
                category: 'Data Management'
            },
            {
                id: 12,
                type: 'print',
                title: 'Batch Printing Completed',
                content: 'Your batch printing task of 50 documents has been completed, taking 15 minutes, average 3.3 pages per minute',
                isRead: false,
                priority: 'normal',
                createdAt: new Date(Date.now() - 1000 * 60 * 20),
                category: 'Print Jobs'
            }
        ],
        showNotificationPanel: false,
        currentNotification: null
    }),
    getters: {
        unreadNotifications: (state) => state.notifications.filter(n => !n.isRead),
        notificationCount: (state) => state.notifications.filter(n => !n.isRead).length,
        notificationsByType: (state) => {
            return state.notifications.reduce((acc, notification) => {
                if (!acc[notification.type]) {
                    acc[notification.type] = []
                }
                acc[notification.type].push(notification)
                return acc
            }, {})
        }
    },
    actions: {
        markAsRead(id) {
            const notification = this.notifications.find(n => n.id === id)
            if (notification) {
                notification.isRead = true
            }
        },
        markAllAsRead() {
            this.notifications.forEach(notification => {
                notification.isRead = true
            })
        },
        deleteNotification(id) {
            const index = this.notifications.findIndex(n => n.id === id)
            if (index > -1) {
                this.notifications.splice(index, 1)
            }
        },
        clearAllNotifications() {
            this.notifications = []
        },
        toggleNotificationPanel() {
            this.showNotificationPanel = !this.showNotificationPanel
        },
        closeAllPanels() {
            this.showNotificationPanel = false
            this.currentNotification = null
        },
        setCurrentNotification(notification) {
            this.currentNotification = notification
            if (notification && !notification.isRead) {
                this.markAsRead(notification.id)
            }
        },
        refreshNotifications() {
// Simulate refresh operation
            console.log('Refresh notification list')
        }
    }
})
