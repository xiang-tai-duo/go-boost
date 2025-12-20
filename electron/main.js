// File:        main.js
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/electron/main.js
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Main entry point for Electron application, handles window creation, app lifecycle and WebSocket communication.
// --------------------------------------------------------------------------------
// noinspection ES6ConvertVarToLetConst,JSUnusedAssignment,JSUnusedLocalSymbols,NpmUsedModulesInstalled,SpellCheckingInspection,JSIgnoredPromiseFromCall,JSUnresolvedReference,HttpUrlsUsage,JSDeprecatedSymbols
// Debuging with JetBrains
//   1. Add npm project
//   2. Set package.json: ..\_go-boost\electron\package.json
//   3. Set Command: start
const { app, BrowserWindow, Menu, session, screen } = require('electron')
const { execFile } = require('child_process')
const path = require('path')
const fs = require('fs')
const process = require('process')
const WebSocket = require('ws')

var ANSI_COLOR_PATTERN = /\x1b\[[0-9;]*m/g
var CONTENT_TYPE_APPLICATION_JSON = 'application/json'
var CONTENT_TYPE_HEADER = 'Content-Type'
var ERROR_CODE_EPIPE = 'EPIPE'
var ERROR_CODE_ERR_STREAM_DESTROYED = 'ERR_STREAM_DESTROYED'
var GREEN = '\x1b[92m'
var HTTP_METHOD_POST = 'POST'
var IMAGE_EXTENSIONS = ['.bmp', '.gif', '.jpeg', '.jpg', '.png', '.svg', '.webp']
var LOG_FILE_EXT = '.log'
var LOG_FILE_MAX_BYTES = 1024 * 1024 * 1024
var LOG_FILE_RETENTION_DAYS = 100
var LOG_FOLDER_NAME = 'Logs'
var LOG_LEVEL_DEBUG = 'DEBUG'
var LOG_LEVEL_ERROR = 'ERROR'
var LOG_LEVEL_INFO = 'INFO'
var LOG_LEVEL_WARNING = 'WARNING'
var LOG_READY_API_PATH = '/api/electron/ready'
var LOG_SOURCE_RENDERER = 'Renderer'
var LOG_SOURCE_WINDOW = 'Window'
var LOG_TAG_CONSOLE = 'Console'
var MACOS = 'darwin'
var MESSAGE_HEADER_LENGTH = 4
var PARAMS_ARG_NAME = '--params'
var PDF_EXTENSION = '.pdf'
var PLATFORM_WIN32 = 'win32'
var PREVIEW_URL_HTTP_PREFIX = 'http://'
var PREVIEW_URL_HTTPS_PREFIX = 'https://'
var PURPLE = '\x1b[95m'
var READY_SIGNAL_VALUE = 'ok'
var RED = '\x1b[31m'
var RESET = '\x1b[0m'
var TIMESTAMP_FORMAT_PAD_LENGTH = 2
var WEBASSEMBLY_EXE_BASENAME = 'webassembly'
var WEBSOCKET_BINARY_TYPE_ARRAYBUFFER = 'arraybuffer'
var WEBSOCKET_DATA_TYPE_EVAL = "eval"
var WEBSOCKET_DATA_TYPE_EXECUTE_JAVASCRIPT = "executeJavaScript"
var WEBSOCKET_DATA_TYPE_TOKEN = "token"
var WEBSOCKET_RECONNECT_DELAY_MAX_MS = 5000
var WEBSOCKET_RECONNECT_DELAY_MIN_MS = 3000
var WEBSOCKET_URL_PREFIX = 'ws://127.0.0.1:'
var WINDOW_DEFAULT_HEIGHT = 600
var WINDOW_DEFAULT_WIDTH = 800
var WINDOW_SIZE_INVALID = -1
var WINDOW_SIZE_MINIMUM = 1
var WINDOW_STATE_CHECK_DELAY_MS = 5000
var authenticationToken = ''
var authenticationTokenKey = ''
var consoleOutputDisabled = false
var currentLogDate = ''
var lastLogCleanupDate = ''
var logFileSizeLimitReached = false
var logFileWriteDisabled = false
var logFileWriteFailureReported = false
var homePage = ''
var homePageLoaded = false
var isDebugMode = false
var isWindowMaximize = false
var isPreviewMode = false
var mainWindow = null
var params = null
var previewUrl = ''
var resizable = true
var backendPid = 0
var websocket = null
var websocketBuffer = new Uint8Array(0)
var websocketOpened = false
var websocketPort = 0
var websocketReconnectAttempts = 0
var websocketReconnectTimer = null
var webassemblyPath = ''
var windowSize = null

logger.info = (...args) => logger('INFO', ...args)
logger.error = (...args) => logger('ERROR', ...args)

function createWebSocket() {
    websocketOpened = false;
    logger.info('Creating WebSocket connection on port:', websocketPort);
    websocket = new WebSocket(WEBSOCKET_URL_PREFIX + websocketPort)
    websocket.binaryType = WEBSOCKET_BINARY_TYPE_ARRAYBUFFER
    websocket.onopen = () => {
        websocketOpened = true;
        websocketReconnectAttempts = 0;
        logger.info('WebSocket connected');
    }
    websocket.onmessage = (event) => {
        try {
            var newData = new Uint8Array(event.data);
            var combined = new Uint8Array(websocketBuffer.length + newData.length);
            combined.set(websocketBuffer);
            combined.set(newData, websocketBuffer.length);
            websocketBuffer = combined;
            while (websocketBuffer.length >= MESSAGE_HEADER_LENGTH) {
                var length = (websocketBuffer[0] << 24) | (websocketBuffer[1] << 16) | (websocketBuffer[2] << 8) | websocketBuffer[3]
                if (websocketBuffer.length >= MESSAGE_HEADER_LENGTH + length) {
                    var messageData = websocketBuffer.slice(MESSAGE_HEADER_LENGTH, MESSAGE_HEADER_LENGTH + length);
                    websocketBuffer = websocketBuffer.slice(MESSAGE_HEADER_LENGTH + length);
                    var message = new TextDecoder().decode(messageData);
                    try {
                        var json = JSON.parse(message);
                        if (typeof json.data === 'string') {
                            var rawData = atob(json.data);
                            var bytes = new Uint8Array(rawData.length);
                            for (var i = 0; i < rawData.length; i++) {
                                bytes[i] = rawData.charCodeAt(i);
                            }
                            var utf8Js = new TextDecoder('utf-8').decode(bytes);
                            if (json.type === WEBSOCKET_DATA_TYPE_EVAL && utf8Js !== '') {
                                if (json.silent !== true) {
                                    logger.info('Electron eval: ' + utf8Js);
                                }
                                eval(utf8Js);
                            } else if (json.type === WEBSOCKET_DATA_TYPE_EXECUTE_JAVASCRIPT && utf8Js !== '') {
                                if (mainWindow) {
                                    if (json.silent !== true) {
                                        logger.info('Electron executeJavaScript: ' + utf8Js);
                                    }
                                    mainWindow.webContents.executeJavaScript(utf8Js);
                                } else {
                                    logger.error('Main window not available for executeJavaScript');
                                }
                            } else if (json.type === WEBSOCKET_DATA_TYPE_TOKEN && rawData !== '') {
                                logger.info('Electron Token: ******');
                                var tokenObj = JSON.parse(rawData);
                                authenticationTokenKey = tokenObj.key;
                                authenticationToken = tokenObj.value;
                                if (homePageLoaded) {
                                    logger.info('Home page already loaded, refresh token only, skip reloading on reconnect');
                                } else {
                                    loadHomePage();
                                }
                            }
                        } else {
                            logger.error('Invalid data format: expected base64 string');
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
    websocket.onclose = (event) => {
        logger.info('WebSocket closed:', {
            code: event.code,
            reason: event.reason,
            wasClean: event.wasClean,
            opened: websocketOpened
        });
        scheduleWebSocketReconnect('close');
    }
    websocket.onerror = (error) => {
        logger.error('WebSocket connection error:', error);
        scheduleWebSocketReconnect('error');
    }
}

function createWindow() {
    logger.info('Creating main window');
    var APP_ICON_FILENAME = 'app.png';
    var PRELOAD_SCRIPT_FILENAME = 'preload.js';
    var windowOptions = {
        width: WINDOW_DEFAULT_WIDTH,
        height: WINDOW_DEFAULT_HEIGHT,
        icon: path.join(__dirname, APP_ICON_FILENAME),
        resizable: true,
        webPreferences: {
            preload: path.join(__dirname, PRELOAD_SCRIPT_FILENAME),
            nodeIntegration: true,
            contextIsolation: false,
        }
    };
    if (resizable === false) {
        windowOptions.resizable = resizable;
    }
    logger.info('Window resizable:', windowOptions.resizable);
    if (windowSize) {
        windowOptions.width = windowSize.width;
        windowOptions.height = windowSize.height;
        windowOptions.resizable = false;
    }
    windowOptions.show = false
    logger.info('Window options:', windowOptions);
    mainWindow = new BrowserWindow(windowOptions);
    logger.info('Main window created');
    if (isWindowMaximize) {
        logger.info('Maximizing main window');
        mainWindow.maximize();
        logger.info('Main window maximized');
    }
    Menu.setApplicationMenu(null)
    if (isDebugMode) {
        mainWindow.webContents.openDevTools();
    }
    hookLogger(mainWindow);
    if (isPreviewMode && previewUrl) {
        if (previewUrl.startsWith(PREVIEW_URL_HTTP_PREFIX) || previewUrl.startsWith(PREVIEW_URL_HTTPS_PREFIX)) {
            logger.info('Loading preview URL:', previewUrl);
            mainWindow.loadURL(previewUrl);
        } else {
            logger.info('Loading preview file:', previewUrl);
            mainWindow.loadFile(previewUrl);
        }
    }
    if (homePage) {
        try {
            var url = new URL(homePage);
            var baseUrl = url.origin;
            var readyUrl = baseUrl + LOG_READY_API_PATH;
            fetch(readyUrl, {
                method: HTTP_METHOD_POST,
                headers: {
                    [CONTENT_TYPE_HEADER]: CONTENT_TYPE_APPLICATION_JSON
                },
                body: JSON.stringify({ ready: READY_SIGNAL_VALUE })
            }).catch((error) => {
                logger.error('Failed to send ready signal:', error);
            });
        } catch (e) {
            logger.error('Error preparing ready signal:', e);
        }
    }
    mainWindow.show();
    logger.info('Main window shown');
}

function formatTimestamp(date) {
    var result = date.getFullYear() + '/' +
        padNumber(date.getMonth() + 1) + '/' +
        padNumber(date.getDate()) + ' ' +
        padNumber(date.getHours()) + ':' +
        padNumber(date.getMinutes()) + ':' +
        padNumber(date.getSeconds());
    return result;
}

function getLogDirectory() {
    var result = __dirname;
    if (app && app.isPackaged) {
        result = path.dirname(process.execPath);
    }
    result = path.join(result, LOG_FOLDER_NAME);
    return result;
}

function getLogFileName(date) {
    return String(date.getFullYear()) + padNumber(date.getMonth() + 1) + padNumber(date.getDate()) + LOG_FILE_EXT;
}

function hookLogger(win) {
    var wc = win.webContents;
    var levelNames = [LOG_LEVEL_DEBUG, LOG_LEVEL_INFO, LOG_LEVEL_WARNING, LOG_LEVEL_ERROR];
    wc.on('console-message', (event, level, message, line, sourceId) => {
        var levelName = levelNames[level] || LOG_LEVEL_INFO;
        var location = sourceId ? ' (' + sourceId + ':' + line + ')' : '';
        if (levelName === LOG_LEVEL_ERROR) {
            logger.error('[' + LOG_SOURCE_RENDERER + '][' + LOG_TAG_CONSOLE + '][' + levelName + '] ' + message + location);
        } else {
            logger.info('[' + LOG_SOURCE_RENDERER + '][' + LOG_TAG_CONSOLE + '][' + levelName + '] ' + message + location);
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

function isValidPreviewUrl(url) {
    var result = false;
    if (typeof url === 'string' && url !== '') {
        if (url.startsWith(PREVIEW_URL_HTTP_PREFIX) || url.startsWith(PREVIEW_URL_HTTPS_PREFIX)) {
            result = true;
        } else {
            try {
                if (fs.existsSync(url)) {
                    var ext = path.extname(url).toLowerCase();
                    if (IMAGE_EXTENSIONS.includes(ext) || ext === PDF_EXTENSION) {
                        result = true;
                    }
                }
            } catch (e) {
                logger.error('Error checking file existence:', e);
            }
        }
    }
    return result;
}

function loadHomePage() {
    logger.info('Loading home page:', homePage);
    if (mainWindow && authenticationToken && authenticationTokenKey && homePage) {
        homePageLoaded = true;
        var cookie = {
            url: homePage,
            name: authenticationTokenKey,
            value: authenticationToken,
            httpOnly: true,
            secure: false
        };
        mainWindow.webContents.session.cookies.set(cookie).then(() => {
            mainWindow.loadURL(homePage);
        }).catch((error) => {
            logger.error('Failed to set cookie:', error);
            mainWindow.loadURL(homePage);
        });
    }
}

function formatLogMessage(args) {
    var result = args.map(arg => {
        var mapped = arg;
        if (arg instanceof Error) {
            mapped = arg.stack || arg.message;
        } else if (typeof arg === 'object') {
            try {
                mapped = JSON.stringify(arg);
            } catch (e) {
                mapped = String(arg);
            }
        }
        return mapped;
    }).join(' ');
    return result;
}

function isBrokenPipeError(error) {
    return error && (error.code === ERROR_CODE_EPIPE || error.code === ERROR_CODE_ERR_STREAM_DESTROYED);
}

function safeConsoleLog(message) {
    if (!consoleOutputDisabled) {
        try {
            console.log(message);
        } catch (error) {
            if (isBrokenPipeError(error)) {
                consoleOutputDisabled = true;
            }
        }
    }
}

function safeConsoleError(message, error) {
    if (!consoleOutputDisabled) {
        try {
            console.error(message, error);
        } catch (e) {
            if (isBrokenPipeError(e)) {
                consoleOutputDisabled = true;
            }
        }
    }
}

function logger(level, ...args) {
    var levelColor = level === LOG_LEVEL_INFO ? GREEN : RED
    var message = formatLogMessage(args)
    var consoleText = `[${PURPLE}ELECTRON${RESET}][${levelColor}${level}${RESET}]${message}`;
    safeConsoleLog(consoleText);
    writeLogFile(level, message);
}

function logProcessInfo() {
    logger.info('Process information:', {
        pid: process.pid,
        ppid: process.ppid,
        platform: process.platform,
        arch: process.arch,
        execPath: process.execPath,
        cwd: process.cwd(),
        dirname: __dirname,
        packaged: app.isPackaged,
        argv: process.argv
    });
    logger.info('Log directory:', getLogDirectory());
}

function padNumber(value) {
    return String(value).padStart(TIMESTAMP_FORMAT_PAD_LENGTH, '0');
}

function parseCommandLine() {
    logger.info('Parsing command line arguments');
    var result = false;
    try {
        for (var i = 0; i < process.argv.length; i++) {
            if (process.argv[i] === PARAMS_ARG_NAME && i + 1 < process.argv.length) {
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
            if (params.preview && typeof params.preview === 'string') {
                if (isValidPreviewUrl(params.preview)) {
                    isPreviewMode = true;
                    previewUrl = params.preview;
                    logger.info('Preview mode enabled, URL:', previewUrl);
                }
            }
            if (params.backendPid) {
                backendPid = parseInt(params.backendPid, 10);
                logger.info('Backend PID:', backendPid);
            }
            if (params.websocketPort) {
                websocketPort = parseInt(params.websocketPort, 10);
            }
            if (params.homePage) {
                homePage = params.homePage;
            }
            if (params.isDebugMode) {
                isDebugMode = true;
            }
            if (params.resizable === false) {
                resizable = false;
            }
            if (params.isWindowMaximize === true) {
                isWindowMaximize = true;
                logger.info('Window maximize enabled');
            }
            if (params.windowSize) {
                var width = Number(params.windowSize.width);
                var height = Number(params.windowSize.height);
                if (!isNaN(width) && !isNaN(height)) {
                    if (width === WINDOW_SIZE_INVALID || height === WINDOW_SIZE_INVALID) {
                        windowSize = null;
                    } else if (width >= WINDOW_SIZE_MINIMUM && height >= WINDOW_SIZE_MINIMUM) {
                        windowSize = { width: width, height: height };
                        logger.info('Window size:', width + 'x' + height);
                    }
                }
            }
        }
        var exeSuffix = '';
        if (process.platform === PLATFORM_WIN32) {
            exeSuffix = '.exe';
        }
        webassemblyPath = path.join(path.dirname(process.execPath), WEBASSEMBLY_EXE_BASENAME + exeSuffix);
        logger.info('WebAssembly path:', webassemblyPath);
        logger.info('Startup parameters:', {
            backendPid: backendPid,
            isPreviewMode: isPreviewMode,
            previewUrl: previewUrl,
            websocketPort: websocketPort,
            homePage: homePage,
            isDebugMode: isDebugMode,
            resizable: resizable,
            isWindowMaximize: isWindowMaximize,
            windowSize: windowSize
        });
        if (isPreviewMode) {
            result = true;
        } else if (websocketPort === 0 || homePage === '') {
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
    logger.info('Quitting Electron application');
    app.quit();
}

function checkBackendProcessAlive() {
    var result = new Promise((resolve, reject) => {
        var alive = false;
        if (backendPid <= 0) {
            logger.warn('Backend PID not set, skipping process check');
            alive = true;
        } else if (!fs.existsSync(webassemblyPath)) {
            logger.error('WebAssembly executable not found:', webassemblyPath);
            alive = false;
        } else {
            execFile(webassemblyPath, [String(backendPid)], (error, stdout, stderr) => {
                if (error) {
                    logger.error('Failed to execute webassembly:', error);
                    resolve(false);
                } else {
                    try {
                        var checkResult = JSON.parse(stdout.trim());
                        logger.info('Backend process check result:', checkResult);
                        resolve(checkResult.alive);
                    } catch (e) {
                        logger.error('Failed to parse webassembly output:', stdout, e);
                        resolve(false);
                    }
                }
            });
            alive = null;
        }
        if (alive !== null) {
            resolve(alive);
        }
    });
    return result;
}

function getRandomReconnectDelay() {
    var minDelay = WEBSOCKET_RECONNECT_DELAY_MIN_MS;
    var maxDelay = WEBSOCKET_RECONNECT_DELAY_MAX_MS;
    return Math.floor(Math.random() * (maxDelay - minDelay + 1)) + minDelay;
}

function scheduleWebSocketReconnect(reason) {
    var shouldReconnect = true;
    if (websocketReconnectTimer) {
        logger.info('WebSocket reconnect already scheduled, reason:', reason);
        shouldReconnect = false;
    }
    if (shouldReconnect) {
        websocketReconnectAttempts++;
        logger.info('WebSocket reconnect attempt:', websocketReconnectAttempts, 'reason:', reason);
        checkBackendProcessAlive().then((alive) => {
            if (alive) {
                var delay = getRandomReconnectDelay();
                logger.info('Backend process alive, scheduling reconnect in', delay, 'ms');
                websocketReconnectTimer = setTimeout(() => {
                    websocketReconnectTimer = null;
                    createWebSocket();
                }, delay);
            } else {
                logger.error('Backend process not alive or webassembly missing, quitting application');
                quitApp();
            }
        }).catch((err) => {
            logger.error('Error checking backend process:', err);
            quitApp();
        });
    }
}

function stripANSIColorCodes(message) {
    return String(message).replace(ANSI_COLOR_PATTERN, '');
}

function logProcessError(prefix, error) {
    if (isBrokenPipeError(error)) {
        consoleOutputDisabled = true;
    }
    try {
        logger.error(prefix + (error && error.stack ? error.stack : error));
    } catch (e) {
        consoleOutputDisabled = true;
        logFileWriteDisabled = true;
    }
}

function getLogFileDateText(fileName) {
    var result = '';
    if (path.extname(fileName) === LOG_FILE_EXT) {
        var matched = path.basename(fileName).match(/^(\d{8})(?:\.|$)/);
        if (matched) {
            result = matched[1];
        }
    }
    return result;
}

function cleanupExpiredLogFiles(logDirectory, today) {
    if (lastLogCleanupDate !== today) {
        lastLogCleanupDate = today;
        var entries = fs.readdirSync(logDirectory, { withFileTypes: true });
        var logDateMap = {};
        entries.forEach(entry => {
            if (!entry.isDirectory()) {
                var logDateText = getLogFileDateText(entry.name);
                if (logDateText !== '') {
                    logDateMap[logDateText] = true;
                }
            }
        });
        var logDates = Object.keys(logDateMap).sort().reverse();
        var retainedDateMap = {};
        logDates.slice(0, LOG_FILE_RETENTION_DAYS).forEach(logDateText => retainedDateMap[logDateText] = true);
        entries.forEach(entry => {
            if (!entry.isDirectory()) {
                var logDateText = getLogFileDateText(entry.name);
                if (logDateText !== '' && !retainedDateMap[logDateText]) {
                    fs.unlinkSync(path.join(logDirectory, entry.name));
                }
            }
        });
    }
}

function prepareDailyLogFile(logDirectory, logFilePath, today) {
    if (currentLogDate !== today) {
        currentLogDate = today;
        logFileSizeLimitReached = false;
        cleanupExpiredLogFiles(logDirectory, today);
    }
    if (!logFileSizeLimitReached && fs.existsSync(logFilePath)) {
        var stat = fs.statSync(logFilePath);
        if (stat.size >= LOG_FILE_MAX_BYTES) {
            logFileSizeLimitReached = true;
        }
    }
}

function writeLogFile(level, message) {
    if (!logFileWriteDisabled) {
        try {
            var now = new Date();
            var logDirectory = getLogDirectory();
            fs.mkdirSync(logDirectory, { recursive: true });
            var today = now.getFullYear() + padNumber(now.getMonth() + 1) + padNumber(now.getDate());
            var logFilePath = path.join(logDirectory, getLogFileName(now));
            prepareDailyLogFile(logDirectory, logFilePath, today);
            if (!logFileSizeLimitReached) {
                var strippedMessage = stripANSIColorCodes(message);
                var fileText = '[' + formatTimestamp(now) + '][main.js::logger][PID:' + process.pid + '][' + level + ']> ' + strippedMessage + '\n';
                if (fs.existsSync(logFilePath)) {
                    var stat = fs.statSync(logFilePath);
                    logFileSizeLimitReached = stat.size + Buffer.byteLength(fileText, 'utf8') > LOG_FILE_MAX_BYTES;
                }
                if (!logFileSizeLimitReached) {
                    fs.appendFileSync(logFilePath, fileText, 'utf8');
                }
            }
        } catch (e) {
            if (!logFileWriteFailureReported) {
                logFileWriteFailureReported = true;
                safeConsoleError(`[${PURPLE}ELECTRON${RESET}][${RED}ERROR${RESET}]Failed to write log file:`, e);
            }
            if (isBrokenPipeError(e)) {
                logFileWriteDisabled = true;
            }
        }
    }
}

logger.info('[Main][ModuleLoad] main.js loaded, execPath=' + process.execPath + ', exists=' + fs.existsSync(process.execPath));
logger.info('[Main][ModuleLoad] app.isReady=' + app.isReady() + ', app.isPackaged=' + app.isPackaged);

if (process.platform === MACOS) {
    app.dock.hide()
}

process.on('uncaughtException', (error) => {
    logProcessError('[Main][UncaughtException] ', error);
});
process.on('unhandledRejection', (reason, promise) => {
    logProcessError('[Main][UnhandledRejection] ', reason);
});
process.on('warning', (warning) => {
    logProcessError('[Main][Warning] ', warning);
});
process.on('exit', (code) => {
    logger.info('[Main][ProcessExit] code=' + code);
});
app.on('ready', () => {
    logger.info('[App][ReadyEvent]');
});
app.on('before-quit', () => {
    logger.info('[App][BeforeQuit]');
});
app.on('will-quit', () => {
    logger.info('[App][WillQuit]');
});
app.on('quit', (event, exitCode) => {
    logger.info('[App][Quit] exitCode=' + exitCode);
});
app.on('render-process-gone', (event, webContents, details) => {
    logger.error('[App][RenderProcessGone] reason=' + details.reason + ', exitCode=' + details.exitCode);
});
app.on('child-process-gone', (event, details) => {
    logger.error('[App][ChildProcessGone] type=' + details.type + ', reason=' + details.reason + ', exitCode=' + details.exitCode + ', serviceName=' + details.serviceName + ', name=' + details.name);
});
setTimeout(() => {
    logger.info('[App][ReadyProbe] app.isReady=' + app.isReady() + ', windows=' + BrowserWindow.getAllWindows().length);
}, WINDOW_STATE_CHECK_DELAY_MS);
logger.info('[App][WhenReadyRegister] registering app.whenReady callback');
app.whenReady().then(() => {
    logger.info('[App][WhenReadyResolved]');
    logger.info('Electron application starting');
    logProcessInfo();
    if (parseCommandLine()) {
        logger.info('Command line parsed successfully');
        createWindow();
        if (!isPreviewMode) {
            logger.info('Starting WebSocket mode');
            createWebSocket();
        } else {
            logger.info('Starting preview mode');
        }
        app.on('activate', () => {
            logger.info('[App][Activate]');
            if (BrowserWindow.getAllWindows().length === 0) {
                createWindow();
            }
        })
        app.on('window-all-closed', () => {
            logger.info('[App][WindowAllClosed]');
            quitApp();
        })
    } else {
        logger.error('Command line parsed failed');
        quitApp();
    }
}).catch((error) => {
    logger.error('[App][WhenReadyRejected] ' + (error && error.stack ? error.stack : error));
    quitApp();
})
