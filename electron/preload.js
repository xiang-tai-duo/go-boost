// File:        sqlite.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/.electron/preload.js
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Preload script for Electron application, provides browser window with access to Node.js APIs.
// --------------------------------------------------------------------------------
window.addEventListener('error', (event) => {
    var msg = event.message || ''
    var src = event.filename || ''
    var line = event.lineno || 0
    var col = event.colno || 0
    var stack = event.error && event.error.stack ? event.error.stack : ''
    console.error('[UncaughtError] ' + msg + ' at ' + src + ':' + line + ':' + col + (stack ? '\n' + stack : ''))
}, true)

window.addEventListener('unhandledrejection', (event) => {
    var reason = event.reason
    var text = ''
    if (reason && reason.stack) {
        text = reason.stack
    } else if (typeof reason === 'object') {
        try { text = JSON.stringify(reason) } catch (e) { text = String(reason) }
    } else {
        text = String(reason)
    }
    console.error('[UnhandledRejection] ' + text)
})

window.addEventListener('DOMContentLoaded', () => {

})
