// File:        sqlite.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/electron/main.js
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Main entry point for Electron application, handles window creation and app lifecycle.
// --------------------------------------------------------------------------------
// noinspection ES6ConvertVarToLetConst,JSUnusedAssignment,JSUnusedLocalSymbols,NpmUsedModulesInstalled,SpellCheckingInspection,JSIgnoredPromiseFromCall,JSUnresolvedReference
// Debuging with WebStorm
//   1. Add npm project
//   2. Set package.json: ..\_go-boost\electron\package.json
//   3. Set Command: start
const {app, BrowserWindow, Menu, session} = require('electron')
const path = require('path')
const fs = require('fs')
const process = require('process')
const WebSocket = require('ws')

const GREEN = '\x1b[92m'
const KEEP_WEBSOCKET_ALIVE = false
const MACOS = 'darwin'
const PURPLE = '\x1b[95m'
const RED = '\x1b[31m'
const RESET = '\x1b[0m'
const WEBSOCKET_DATA_TYPE_EVAL = 'eval'
const WEBSOCKET_DATA_TYPE_TOKEN = 'token'

var authenticationToken = ''
var authenticationTokenKey = ''
var homePage = ''
var idDebugMode = false
var mainWindow = null
var params = null
var websocket = null
var websocketBuffer = new Uint8Array(0)
var websocketPort = 0
var windowSize = null

logger.info = (...args) => logger('INFO', ...args)
logger.error = (...args) => logger('ERROR', ...args)

function attachRendererLogging(win) {
    var wc = win.webContents;
    var levelNames = ['DEBUG', 'INFO', 'WARNING', 'ERROR'];
    wc.on('console-message', (event, level, message, line, sourceId) => {
        var levelName = levelNames[level] || 'INFO';
        var location = sourceId ? ' (' + sourceId + ':' + line + ')' : '';
        if (levelName === 'ERROR') {
            logger.error('[Renderer][Console][' + levelName + '] ' + message + location);
        } else {
            logger.info('[Renderer][Console][' + levelName + '] ' + message + location);
        }
    });
    wc.on('render-process-gone', (event, details) => {
        logger.error('[Renderer][ProcessGone] reason=' + details.reason + ', exitCode=' + details.exitCode);
    });
    wc.on('unresponsive', () => {
        logger.error('[Renderer][Unresponsive] The page has become unresponsive');
    });
    wc.on('responsive', () => {
        logger.info('[Renderer][Responsive] The page has become responsive again');
    });
    wc.on('preload-error', (event, preloadPath, error) => {
        logger.error('[Renderer][PreloadError] path=' + preloadPath + ', error=' + (error && error.stack ? error.stack : error));
    });
    wc.on('did-fail-load', (event, errorCode, errorDescription, validatedURL, isMainFrame) => {
        logger.error('[Renderer][DidFailLoad] code=' + errorCode + ', desc=' + errorDescription + ', url=' + validatedURL + ', mainFrame=' + isMainFrame);
    });
    wc.on('did-fail-provisional-load', (event, errorCode, errorDescription, validatedURL, isMainFrame) => {
        logger.error('[Renderer][DidFailProvisionalLoad] code=' + errorCode + ', desc=' + errorDescription + ', url=' + validatedURL + ', mainFrame=' + isMainFrame);
    });
    wc.on('plugin-crashed', (event, name, version) => {
        logger.error('[Renderer][PluginCrashed] name=' + name + ', version=' + version);
    });
    wc.on('certificate-error', (event, url, error) => {
        logger.error('[Renderer][CertificateError] url=' + url + ', error=' + error);
    });
    wc.on('did-finish-load', () => {
        logger.info('[Renderer][DidFinishLoad] Page finished loading');
    });
    wc.on('crashed', (event, killed) => {
        logger.error('[Renderer][Crashed] killed=' + killed);
    });
    win.on('unresponsive', () => {
        logger.error('[Window][Unresponsive] The window has become unresponsive');
    });
    win.on('closed', () => {
        logger.info('[Window][Closed] The window has been closed');
    });
}

function createWindow() {
    var windowOptions = {
        width: 800,
        height: 600,
        icon: path.join(__dirname, 'app.png'),
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            nodeIntegration: true,
            contextIsolation: false,
        }
    };
    if (windowSize) {
        windowOptions.width = windowSize.width;
        windowOptions.height = windowSize.height;
        windowOptions.resizable = false;
    }
    mainWindow = new BrowserWindow(windowOptions);
    if (!windowSize) {
        mainWindow.maximize();
    }
    Menu.setApplicationMenu(null)
    if (idDebugMode) {
        mainWindow.webContents.openDevTools();
    }
    attachRendererLogging(mainWindow);
}

function createWebSocket() {
    websocket = new WebSocket(`ws://127.0.0.1:${websocketPort}`)
    websocket.binaryType = 'arraybuffer'
    websocket.onopen = () => { }
    websocket.onmessage = (event) => {
        try {
            var newData = new Uint8Array(event.data);
            var combined = new Uint8Array(websocketBuffer.length + newData.length);
            combined.set(websocketBuffer);
            combined.set(newData, websocketBuffer.length);
            websocketBuffer = combined;
            while (websocketBuffer.length >= 4) {
                var length = (websocketBuffer[0] << 24) | (websocketBuffer[1] << 16) | (websocketBuffer[2] << 8) | websocketBuffer[3]
                if (websocketBuffer.length >= 4 + length) {
                    var messageData = websocketBuffer.slice(4, 4 + length);
                    websocketBuffer = websocketBuffer.slice(4 + length);
                    var message = new TextDecoder().decode(messageData);
                    try {
                        var json = JSON.parse(message);
                        if (json.type === WEBSOCKET_DATA_TYPE_EVAL && typeof json.data === 'string' && json.data !== '') {
                            logger.info('Websocket data recevied: JSON, ', json.data);
                            eval(json.data);
                        } else if (json.type === WEBSOCKET_DATA_TYPE_TOKEN && typeof json.data === 'object' && json.data !== null) {
                            logger.info('Websocket data recevied: Token');
                            authenticationTokenKey = json.data.key;
                            authenticationToken = json.data.value;
                            loadHomePage();
                        }
                    } catch (e) {
                        logger.error('Failed to parse or execute message:', e);
                    }
                } else {
                    break
                }
            }
        } catch (error) {
            logger.error('WebSocket error:', error);
            websocketBuffer = new Uint8Array(0);
        }
    }
    websocket.onclose = () => {
        if (KEEP_WEBSOCKET_ALIVE && idDebugMode) {
            setTimeout(createWebSocket, 1000);
        } else {
            quitApp();
        }
    }
    websocket.onerror = () => {
        if (KEEP_WEBSOCKET_ALIVE && idDebugMode) {
            setTimeout(createWebSocket, 1000);
        } else {
            quitApp();
        }
    }
}

if (process.platform === MACOS) {
    app.dock.hide()
}

process.on('uncaughtException', (error) => {
    logger.error('[Main][UncaughtException] ' + (error && error.stack ? error.stack : error));
});
process.on('unhandledRejection', (reason, promise) => {
    logger.error('[Main][UnhandledRejection] ' + (reason && reason.stack ? reason.stack : reason));
});

app.whenReady().then(() => {
    if (parseCommandLine()) {
        createWindow();
        createWebSocket();
        app.on('activate', () => {
            if (BrowserWindow.getAllWindows().length === 0) {
                createWindow();
            }
        })
        app.on('window-all-closed', () => {
            quitApp();
        })
    } else {
        quitApp();
    }
})

function loadHomePage() {
    if (mainWindow && authenticationToken && authenticationTokenKey && homePage) {
        var cookie = {
            url: homePage,
            name: authenticationTokenKey,
            value: authenticationToken,
            httpOnly: true,
            secure: false
        };
        mainWindow.webContents.session.cookies.set(cookie).then(() => {
            logger.info('Cookie set successfully');
            mainWindow.loadURL(homePage);
        }).catch((error) => {
            logger.error('Failed to set cookie:', error);
            mainWindow.loadURL(homePage);
        });
    }
}

function logger(level, ...args) {
    var levelColor = level === 'INFO' ? GREEN : RED
    var message = args.map(arg => typeof arg === 'object' ? JSON.stringify(arg) : arg).join(' ')
    console.log(`[${PURPLE}ELECTRON${RESET}][${levelColor}${level}${RESET}]${message}`)
}

function parseCommandLine() {
    var result = false;
    try {
        for (var i = 0; i < process.argv.length; i++) {
            if (process.argv[i] === '--params' && i + 1 < process.argv.length) {
                try {
                    var base64Data = process.argv[i + 1];
                    var jsonString = Buffer.from(base64Data, 'base64').toString('utf-8');
                    params = JSON.parse(jsonString);
                    logger.info('Params loaded from base64:', params);
                } catch (e) {
                    logger.error('Failed to parse --params:', e);
                }
            }
        }
        if (params) {
            if (params.websocketPort) {
                websocketPort = parseInt(params.websocketPort, 10);
            }
            if (params.homePage) {
                homePage = params.homePage;
            }
            if (params.isDebugMode) {
                idDebugMode = true;
            }
            if (params.windowSize) {
                var width = parseInt(params.windowSize.width, 10);
                var height = parseInt(params.windowSize.height, 10);
                if (!isNaN(width) && !isNaN(height)) {
                    if (width === -1 || height === -1) {
                        windowSize = null;
                    } else if (width > 0 && height > 0) {
                        windowSize = {width: width, height: height};
                        logger.info('Window size:', width + 'x' + height);
                    }
                }
            }
        }
        if (websocketPort === 0 || homePage === '') {
            logger.error('Missing required arguments: websocketPort or homePage in params');
            result = false;
            quitApp();
        } else {
            result = true;
        }
    } catch (error) {
        logger.error('parseCommandLine error:', error);
    }
    return result;
}

function quitApp() {
    app.quit();
}
