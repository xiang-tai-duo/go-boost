/* --------------------------------------------------------------------------------
 * File:        preload.js
 * Author:      TRAE AI
 * Created:     2025/12/20 12:31:58
 * Description: Preload script for Electron application, provides browser window with access to Node.js APIs.
 * --------------------------------------------------------------------------------
 */
window.addEventListener('DOMContentLoaded', () => {
  const replaceText = (selector, text) => {
    const element = document.getElementById(selector)
    if (element) element.innerText = text
  }

  for (const dependency of ['chrome', 'node', 'electron']) {
    replaceText(`${dependency}-version`, process.versions[dependency])
  }
})