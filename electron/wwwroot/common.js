// File:        common.js
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/electron/wwwroot/common.js
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Common helper functions for popup message boxes (error / warning / info).
// --------------------------------------------------------------------------------

var DOCUMENT_TAG_DIVISION = 'div';
var DOCUMENT_EVENT_CLICK = 'click';
var DOCUMENT_EVENT_KEYDOWN = 'keydown';
var CSS_CLASS_HIDDEN = 'hidden';
var JAVASCRIPT_TYPE_UNDEFINED = 'undefined';
var TEXT_CONFIRM_BUTTON_OK = 'OK';
var KEY_ESCAPE = 'Escape';
var KEY_ENTER = 'Enter';

var popupOverlayMap = {};

function createOverlayCore(overlayId, headerClass, bodyClass, messageClass, footerClass, iconClass, defaultTitle) {
    var overlay = undefined;
    var panel = undefined;
    var header = undefined;
    var icon = undefined;
    var titleSpan = undefined;
    var body = undefined;
    var messageDiv = undefined;
    var footer = undefined;
    var btn = undefined;
    overlay = document.createElement(DOCUMENT_TAG_DIVISION);
    overlay.className = 'popup-overlay hidden';
    overlay.id = overlayId;
    panel = document.createElement(DOCUMENT_TAG_DIVISION);
    panel.className = 'popup-panel';
    header = document.createElement(DOCUMENT_TAG_DIVISION);
    header.className = headerClass;
    icon = document.createElement('i');
    icon.className = iconClass;
    titleSpan = document.createElement('span');
    titleSpan.className = messageClass === 'error2-message' ? 'error2-title' : (messageClass === 'warning2-message' ? 'warning2-title' : 'info2-title');
    header.appendChild(icon);
    header.appendChild(titleSpan);
    body = document.createElement(DOCUMENT_TAG_DIVISION);
    body.className = bodyClass;
    messageDiv = document.createElement(DOCUMENT_TAG_DIVISION);
    messageDiv.className = messageClass;
    body.appendChild(messageDiv);
    footer = document.createElement(DOCUMENT_TAG_DIVISION);
    footer.className = footerClass;
    btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'btn-primary';
    btn.innerHTML = '<i class="fas fa-check" style="margin-right:6px;"></i><span>' + TEXT_CONFIRM_BUTTON_OK + '</span>';
    footer.appendChild(btn);
    panel.appendChild(header);
    panel.appendChild(body);
    panel.appendChild(footer);
    overlay.appendChild(panel);
    document.body.appendChild(overlay);
    popupOverlayMap[overlayId] = overlay;
    btn.addEventListener(DOCUMENT_EVENT_CLICK, function () {
        hideError(overlayId);
        if (overlay._callback) {
            overlay._callback();
        }
    });
    return overlay;
}

function findTopmostVisibleOverlay() {
    var result = undefined;
    var overlayId = undefined;
    var overlay = undefined;
    for (overlayId in popupOverlayMap) {
        if (popupOverlayMap.hasOwnProperty(overlayId)) {
            overlay = popupOverlayMap[overlayId];
            if (overlay && !overlay.classList.contains(CSS_CLASS_HIDDEN)) {
                result = overlay;
            }
        }
    }
    return result;
}

function handlePopupOverlayKeydown(e) {
    var overlay = undefined;
    var primaryButton = undefined;
    if (e.key === KEY_ESCAPE || e.key === KEY_ENTER) {
        overlay = findTopmostVisibleOverlay();
        if (overlay) {
            if (overlay._blocking) {
                e.preventDefault();
            } else {
                primaryButton = overlay.querySelector('.btn-primary');
                if (e.key === KEY_ENTER) {
                    if (primaryButton) {
                        e.preventDefault();
                        primaryButton.click();
                    }
                } else {
                    e.preventDefault();
                    if (primaryButton) {
                        primaryButton.click();
                    } else {
                        hideError(overlay.id);
                    }
                }
            }
        }
    }
}

function hideError(overlayId) {
    var overlay = popupOverlayMap[overlayId];
    if (!overlay) {
        overlay = document.getElementById(overlayId);
    }
    if (overlay) {
        overlay.classList.add(CSS_CLASS_HIDDEN);
        if (typeof coordinateOverlays !== JAVASCRIPT_TYPE_UNDEFINED) {
            coordinateOverlays();
        }
    }
}

function show(overlayId, headerClass, bodyClass, messageClass, footerClass, iconClass, title, message, callback, defaultTitle) {
    var isActive = true;
    var overlay = undefined;
    var existing = undefined;
    var titleEl = undefined;
    var messageEl = undefined;
    if (typeof isInactiveMode !== JAVASCRIPT_TYPE_UNDEFINED && isInactiveMode) {
        isActive = false;
    }
    if (isActive) {
        existing = popupOverlayMap[overlayId];
        if (existing) {
            overlay = existing;
        } else {
            overlay = createOverlayCore(overlayId, headerClass, bodyClass, messageClass, footerClass, iconClass, defaultTitle);
        }
        overlay._callback = callback;
        titleEl = overlay.querySelector(messageClass === 'error2-message' ? '.error2-title' : (messageClass === 'warning2-message' ? '.warning2-title' : '.info2-title'));
        messageEl = overlay.querySelector('.' + messageClass);
        if (titleEl) {
            titleEl.textContent = title || defaultTitle;
        }
        if (messageEl) {
            messageEl.innerHTML = message || '';
        }
        overlay.classList.remove(CSS_CLASS_HIDDEN);
        if (typeof coordinateOverlays !== JAVASCRIPT_TYPE_UNDEFINED) {
            coordinateOverlays();
        }
    }
}

function showError(title, message, callback) {
    return show(
        'divError2Overlay',
        'error2-header',
        'error2-body',
        'error2-message',
        'error2-footer',
        'fas fa-exclamation-circle',
        title,
        message,
        callback,
        'Error'
    );
}

function showBlockingError(title, message) {
    var overlay = undefined;
    var footer = undefined;
    show(
        'divBlockingError2Overlay',
        'error2-header',
        'error2-body',
        'error2-message',
        'error2-footer',
        'fas fa-exclamation-circle',
        title,
        message,
        undefined,
        'Error'
    );
    overlay = popupOverlayMap['divBlockingError2Overlay'];
    if (overlay) {
        overlay._blocking = true;
        footer = overlay.querySelector('.error2-footer');
        if (footer) {
            footer.style.display = 'none';
        }
    }
}

function showInfo(title, message, callback) {
    return show(
        'divInfo2Overlay',
        'info2-header',
        'info2-body',
        'info2-message',
        'info2-footer',
        'fas fa-info-circle',
        title,
        message,
        callback,
        'Info'
    );
}

function showWarning(title, message, callback) {
    return show(
        'divWarning2Overlay',
        'warning2-header',
        'warning2-body',
        'warning2-message',
        'warning2-footer',
        'fas fa-exclamation-triangle',
        title,
        message,
        callback,
        'Warning'
    );
}

document.addEventListener(DOCUMENT_EVENT_KEYDOWN, handlePopupOverlayKeydown);
