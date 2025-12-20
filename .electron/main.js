/* --------------------------------------------------------------------------------
 * File:        main.js
 * Author:      TRAE AI
 * Created:     2025/12/20 12:31:58
 * Description: Main entry point for Electron application, handles window creation and app lifecycle.
 * --------------------------------------------------------------------------------
 */
const { app, BrowserWindow } = require('electron')
const path = require('path')

function createWindow () {
  // Create browser window
  const win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: true,
      contextIsolation: false
    }
  })

  // Load index.html file
  win.loadFile('index.html')

  // Open developer tools (optional)
  // win.webContents.openDevTools()
}

// Called when Electron is ready to create browser windows
app.whenReady().then(() => {
  createWindow()

  // On macOS, re-create a window when the dock icon is clicked and no windows are open
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })

  // Quit the app when all windows are closed
  app.on('window-all-closed', () => {
    // On macOS, keep the app running unless explicitly quit with Cmd + Q
    if (process.platform !== 'darwin') {
      app.quit()
    }
  })