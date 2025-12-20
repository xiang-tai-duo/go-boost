//go:build windows

// File:        mpdefender.h
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/defender/mpdefender.h
// Author:      TRAE.AI
// Created:     2026/06/24 00:00:00
// Description: Windows Defender exclusion management header
// --------------------------------------------------------------------------------

#ifndef GO_BOOST_MPDEFENDER_H
#define GO_BOOST_MPDEFENDER_H

#include <windows.h>

#ifdef __cplusplus
extern "C" {
#endif

#define MPDEFENDER_S_OK           0x00000000
#define MPDEFENDER_E_ACCESSDENIED 0x80070005
#define MPDEFENDER_E_INVALIDARG    0x80070057

HRESULT MpDefenderAddExclusion(const wchar_t* path);

HRESULT MpDefenderRemoveExclusion(const wchar_t* path);

#ifdef __cplusplus
}
#endif

#endif
