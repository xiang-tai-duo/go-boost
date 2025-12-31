// --------------------------------------------------------------------------------
// File:        sidebar.js
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Pinia store for managing page navigation state.
// --------------------------------------------------------------------------------
import {defineStore} from 'pinia'

export const usePageNavigationStore = defineStore('pageNavigation', {
    state: () => ({
        showUserManagement: false,
        showAdvancedUserManagement: false,
        showDeviceManagement: false,
        showDeviceGroupManagement: false,
        showDepartmentManagement: false,
        showDepartmentAccountingReport: false,
        showUserAccountingReport: false,
        showPrintTrackingReport: false,
        showReportGeneration: false
    }),
    actions: {
        openUserManagement() {
            this.showUserManagement = true
        },
        closeUserManagement() {
            this.showUserManagement = false
        },
        openDepartmentManagement() {
            this.showDepartmentManagement = true
        },
        closeDepartmentManagement() {
            this.showDepartmentManagement = false
        },
        openDeviceManagement() {
            this.showDeviceManagement = true
        },
        closeDeviceManagement() {
            this.showDeviceManagement = false
        },
        openDepartmentAccountingReport() {
            this.showDepartmentAccountingReport = true
        },
        closeDepartmentAccountingReport() {
            this.showDepartmentAccountingReport = false
        },
        openUserAccountingReport() {
            this.showUserAccountingReport = true
        },
        closeUserAccountingReport() {
            this.showUserAccountingReport = false
        },
        openPrintTrackingReport() {
            this.showPrintTrackingReport = true
        },
        closePrintTrackingReport() {
            this.showPrintTrackingReport = false
        },
        openReportGeneration() {
            this.showReportGeneration = true
        },
        closeReportGeneration() {
            this.showReportGeneration = false
        },
// Device Group Management
        openDeviceGroupManagement() {
            this.showDeviceGroupManagement = true
        },
        closeDeviceGroupManagement() {
            this.showDeviceGroupManagement = false
        },
        closeAllPages() {
            this.showUserManagement = false
            this.showDeviceManagement = false
            this.showDeviceGroupManagement = false
            this.showDepartmentManagement = false
            this.showDepartmentAccountingReport = false
            this.showUserAccountingReport = false
            this.showPrintTrackingReport = false
            this.showReportGeneration = false
        }
    }
})
