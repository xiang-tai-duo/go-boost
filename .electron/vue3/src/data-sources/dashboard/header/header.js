// --------------------------------------------------------------------------------
// File:        header.js
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Pinia store for managing header state including emails and settings.
// --------------------------------------------------------------------------------
import {defineStore} from 'pinia'

export const useHeaderStore = defineStore('header', {
    state: () => ({
        emails: [
            {
                id: 1,
                from: 'admin@company.com',
                fromName: 'System Administrator',
                to: 'user@company.com',
                subject: 'System Maintenance Notice',
                content: 'Dear Colleagues,\n\nThe system will undergo routine maintenance this Saturday from 2:00-4:00 AM, which may affect normal usage.\n\nMaintenance content:\n• Database optimization and upgrade\n• Security patch update\n• Performance tuning\n\nPlease prepare in advance, we apologize for any inconvenience caused.\n\nTechnical Support Team',
                isRead: false,
                priority: 'high',
                createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000),
                attachments: ['Maintenance Plan.pdf', 'Rollback Plan.docx'],
                status: 'inbox'
            },
            {
                id: 2,
                from: 'hr@company.com',
                fromName: 'Human Resources',
                to: 'user@company.com',
                subject: 'Notice on Office Hours Adjustment',
                content: 'Dear Colleagues,\n\nEffective next Monday, office hours will be adjusted to:\nMorning: 9:00-12:00\nAfternoon: 13:30-18:00\n\nPlease be informed and inform each other.\n\nHuman Resources Department',
                isRead: true,
                priority: 'normal',
                createdAt: new Date(Date.now() - 24 * 60 * 60 * 1000),
                attachments: [],
                status: 'inbox'
            },
            {
                id: 3,
                from: 'it@company.com',
                fromName: 'IT Department',
                to: 'colleague@company.com',
                subject: 'Network Security Reminder',
                content: 'Dear Colleagues,\n\nRecently, multiple phishing email attacks have been detected. Please note:\n1. Do not click on unknown links\n2. Do not download suspicious attachments\n3. Update antivirus software in a timely manner\n\nPlease contact the IT department immediately if any abnormalities are found.\n\nIT Support Team',
                isRead: false,
                priority: 'high',
                createdAt: new Date(Date.now() - 3 * 60 * 60 * 1000),
                attachments: ['Security Manual.pdf'],
                status: 'sent'
            },
            {
                id: 4,
                from: 'finance@company.com',
                fromName: 'Finance Department',
                to: 'user@company.com',
                subject: 'Monthly Expense Reimbursement Reminder',
                content: 'Dear Colleagues,\n\nThe deadline for this month\'s expense reimbursement is the 5th of next month. Please submit the relevant documents in a timely manner.\n\nReimbursement process:\n1. Fill in the electronic reimbursement form\n2. Department manager review\n3. Finance department review\n4. Financial director approval\n\nIf you have any questions, please contact Xiao Wang from the Finance Department.\n\nFinance Department',
                isRead: true,
                priority: 'normal',
                createdAt: new Date(Date.now() - 48 * 60 * 60 * 1000),
                attachments: [],
                status: 'inbox'
            },
            {
                id: 5,
                from: 'marketing@company.com',
                fromName: 'Marketing Department',
                to: 'all@company.com',
                subject: 'New Product Launch Notice',
                content: 'Dear Colleagues,\n\nWe will hold a new product launch in the multifunctional hall at 2:00 PM next Wednesday. Please attend on time.\n\nLaunch process:\n• Product introduction\n• Demo session\n• Q&A interaction\n\nPlease ask department heads to organize employees to attend.\n\nMarketing Department',
                isRead: true,
                priority: 'high',
                createdAt: new Date(Date.now() - 72 * 60 * 60 * 1000),
                attachments: ['Launch Agenda.pdf'],
                status: 'sent'
            }
        ],
        settings: {
            personal: {
                language: 'en-US',
                theme: 'light',
                notifications: true,
                emailNotifications: true
            },
            system: {
                autoSave: true,
                backupFrequency: 'daily',
                dataRetention: '3months'
            },
            advanced: {
                apiTimeout: 30000,
                maxRetries: 3,
                cacheEnabled: true
            }
        },
        showEmailPanel: false,
        showSettingsPanel: false,
        currentEmail: null
    }),
    getters: {
        unreadEmails: (state) => {
            return state.emails.filter(email => !email.isRead && email.status === 'inbox')
        },
        emailCount: (state) => {
            return state.emails.filter(email => !email.isRead && email.status === 'inbox').length
        },
        emailsByStatus: (state) => {
            const result = {}
            state.emails.forEach(email => {
                if (!result[email.status]) {
                    result[email.status] = []
                }
                result[email.status].push(email)
            })
            return result
        }
    },
    actions: {
        markEmailAsRead(id) {
            const email = this.emails.find(e => e.id === id)
            if (email) {
                email.isRead = true
            }
        },
        markAllEmailsAsRead() {
            this.emails.forEach(email => {
                email.isRead = true
            })
        },
        deleteEmail(id) {
            const index = this.emails.findIndex(e => e.id === id)
            if (index > -1) {
                this.emails.splice(index, 1)
            }
        },
        moveEmailToTrash(id) {
            const email = this.emails.find(e => e.id === id)
            if (email) {
                email.status = 'trash'
            }
        },
        toggleEmailPanel() {
            this.showEmailPanel = !this.showEmailPanel
        },
        toggleSettingsPanel() {
            this.showSettingsPanel = !this.showSettingsPanel
        },
        setCurrentEmail(email) {
            this.currentEmail = email
        },
        closeAllPanels() {
            this.showEmailPanel = false
            this.showSettingsPanel = false
        },
        saveSetting(category, key, value) {
            if (this.settings[category]) {
                this.settings[category][key] = value
            }
        },
        saveAllSettings(category, settings) {
            if (this.settings[category]) {
                this.settings[category] = {...this.settings[category], ...settings}
            }
        },
        resetSettings(category) {
            if (category === 'personal') {
                this.settings.personal = {
                    language: 'zh-CN',
                    theme: 'light',
                    notifications: true,
                    emailNotifications: true
                }
            } else if (category === 'system') {
                this.settings.system = {
                    autoSave: true,
                    backupFrequency: 'daily',
                    dataRetention: '3months'
                }
            } else if (category === 'advanced') {
                this.settings.advanced = {
                    apiTimeout: 30000,
                    maxRetries: 3,
                    cacheEnabled: true
                }
            }
        },
        getSetting(category, key) {
            if (this.settings[category] && this.settings[category][key] !== undefined) {
                return this.settings[category][key]
            }
            return null
        }
    }
})
