//go:build windows

// File:        mpdefender.c
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/defender/mpdefender.c
// Author:      TRAE.AI
// Created:     2026/06/24 00:00:00
// Description: Windows Defender exclusion management implementation
// --------------------------------------------------------------------------------

#include "mpdefender.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#pragma comment(lib, "advapi32.lib")

static const wchar_t* DEFENDER_EXCLUSION_PATH = L"SOFTWARE\\Policies\\Microsoft\\Windows Defender\\Exclusions\\Paths";

HRESULT MpDefenderAddExclusion(const wchar_t* path) {
    HRESULT hr = S_OK;
    HKEY hKey = NULL;
    DWORD dwDisposition = 0;
    LONG lResult = 0;
    wchar_t expandedPath[MAX_PATH];
    
    ExpandEnvironmentStringsW(path, expandedPath, MAX_PATH);
    
    lResult = RegCreateKeyExW(
        HKEY_LOCAL_MACHINE,
        DEFENDER_EXCLUSION_PATH,
        0,
        NULL,
        0,
        KEY_SET_VALUE,
        NULL,
        &hKey,
        &dwDisposition
    );
    
    if (lResult == ERROR_SUCCESS) {
        lResult = RegSetValueExW(
            hKey,
            expandedPath,
            0,
            REG_DWORD,
            (const BYTE*)&dwDisposition,
            sizeof(DWORD)
        );
        
        if (lResult == ERROR_SUCCESS) {
            RegCloseKey(hKey);
        } else {
            hr = HRESULT_FROM_WIN32(lResult);
        }
    } else {
        hr = HRESULT_FROM_WIN32(lResult);
    }
    
    return hr;
}

HRESULT MpDefenderRemoveExclusion(const wchar_t* path) {
    HRESULT hr = S_OK;
    HKEY hKey = NULL;
    LONG lResult = 0;
    wchar_t expandedPath[MAX_PATH];
    
    ExpandEnvironmentStringsW(path, expandedPath, MAX_PATH);
    
    lResult = RegOpenKeyExW(
        HKEY_LOCAL_MACHINE,
        DEFENDER_EXCLUSION_PATH,
        0,
        KEY_SET_VALUE,
        &hKey
    );
    
    if (lResult == ERROR_SUCCESS) {
        lResult = RegDeleteValueW(hKey, expandedPath);
        
        if (lResult == ERROR_SUCCESS || lResult == ERROR_FILE_NOT_FOUND) {
            RegCloseKey(hKey);
        } else {
            hr = HRESULT_FROM_WIN32(lResult);
        }
    } else {
        hr = HRESULT_FROM_WIN32(lResult);
    }
    
    return hr;
}
