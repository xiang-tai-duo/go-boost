//go:build windows

// File:        mpfirewall.h
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/mpfirewall/mpfirewall.h
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Windows Firewall COM API header
// --------------------------------------------------------------------------------

#ifndef GO_BOOST_MPFIREWALL_H
#define GO_BOOST_MPFIREWALL_H

#include <windows.h>
#include <objbase.h>
#include <netfw.h>

#ifdef __cplusplus
extern "C" {
#endif

#define MPFIREWALL_S_OK               0x00000000
#define MPFIREWALL_E_ACCESSDENIED     0x80070005
#define MPFIREWALL_E_NO_SUCH_RULE     0x80070002

#define MPFIREWALL_RULE_DIR_INBOUND   1
#define MPFIREWALL_RULE_DIR_OUTBOUND  2

#define MPFIREWALL_ACTION_ALLOW       1
#define MPFIREWALL_ACTION_BLOCK       0

#define MPFIREWALL_PROFILE_ALL       0x7FFFFFFF
#define MPFIREWALL_PROFILE_DOMAIN    0x00000001
#define MPFIREWALL_PROFILE_PRIVATE   0x00000002
#define MPFIREWALL_PROFILE_PUBLIC    0x00000004

#define MPFIREWALL_PROTOCOL_TCP     6
#define MPFIREWALL_PROTOCOL_UDP     17
#define MPFIREWALL_PROTOCOL_ANY     256

typedef struct {
    char *name;
    char *description;
    char *application_name;
    char *local_ports;
    int direction;
    int action;
    int protocol;
    int profiles;
    int enabled;
} MPFIREWALL_RULE_T;

typedef struct {
    char **rule_names;
    int count;
} MPFIREWALL_RULE_NAME_LIST_T;

typedef struct {
    char *name;
    char *local_port;
    int found;
} MPFIREWALL_RULE_INFO_T;

HRESULT mpfirewall_init(void **ppPolicy2);

void mpfirewall_release(void *pPolicy2);

HRESULT mpfirewall_add_rule(
    void *pPolicy2,
    const char *name,
    const char *description,
    const char *application_name,
    const char *local_ports,
    int direction,
    int action,
    int protocol,
    int profiles,
    int enabled
);

HRESULT mpfirewall_delete_rule(void *pPolicy2, const char *ruleName);

MPFIREWALL_RULE_NAME_LIST_T *mpfirewall_get_blocked_rules(void *pPolicy2, const char *exeFilePath, int port);

void mpfirewall_free_rule_name_list(MPFIREWALL_RULE_NAME_LIST_T *list);

HRESULT mpfirewall_get_rule_local_port(void *pPolicy2, const char *ruleName, char **outLocalPort);

MPFIREWALL_RULE_INFO_T *mpfirewall_is_rule_exists(void *pPolicy2, const char *exeFilePath, int port);

void mpfirewall_free_rule_info(MPFIREWALL_RULE_INFO_T *info);

#ifdef __cplusplus
}
#endif

#endif
