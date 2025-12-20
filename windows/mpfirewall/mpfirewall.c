//go:build windows

// File:        mpfirewall.c
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/mpfirewall/mpfirewall.c
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Windows Firewall COM API implementation
// --------------------------------------------------------------------------------

#define INITGUID
#include "mpfirewall.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <wchar.h>

static BSTR utf8_to_bstr(const char *utf8) {
    BSTR result = NULL;
    int len = 0;
    int wlen = 0;
    wchar_t *wbuf = NULL;
    if (utf8 == NULL) {
        result = NULL;
    } else {
        len = (int)strlen(utf8);
        if (len == 0) {
            result = SysAllocString(L"");
        } else {
            wlen = MultiByteToWideChar(CP_UTF8, 0, utf8, -1, NULL, 0);
            if (wlen <= 0) {
                result = NULL;
            } else {
                wbuf = (wchar_t *)malloc(wlen * sizeof(wchar_t));
                if (wbuf == NULL) {
                    result = NULL;
                } else {
                    MultiByteToWideChar(CP_UTF8, 0, utf8, -1, wbuf, wlen);
                    result = SysAllocString(wbuf);
                    free(wbuf);
                }
            }
        }
    }
    return result;
}

static char *bstr_to_utf8(BSTR bstr) {
    char *result = NULL;
    int len = 0;
    int ulen = 0;
    if (bstr == NULL) {
        result = NULL;
    } else {
        len = SysStringLen(bstr);
        if (len == 0) {
            result = (char *)malloc(1);
            if (result != NULL) {
                result[0] = '\0';
            }
        } else {
            ulen = WideCharToMultiByte(CP_UTF8, 0, bstr, len, NULL, 0, NULL, NULL);
            if (ulen <= 0) {
                result = NULL;
            } else {
                result = (char *)malloc(ulen + 1);
                if (result != NULL) {
                    WideCharToMultiByte(CP_UTF8, 0, bstr, len, result, ulen, NULL, NULL);
                    result[ulen] = '\0';
                }
            }
        }
    }
    return result;
}

static int wcsieq(const wchar_t *a, const wchar_t *b) {
    int result = 0;
    if (a == NULL && b == NULL) {
        result = 1;
    } else if (a != NULL && b != NULL) {
        result = _wcsicmp(a, b) == 0;
    } else {
        result = 0;
    }
    return result;
}

HRESULT mpfirewall_init(void **ppPolicy2) {
    HRESULT hr = CoInitializeEx(0, COINIT_APARTMENTTHREADED);
    INetFwPolicy2 *pPolicy2 = NULL;
    if (FAILED(hr) && hr != RPC_E_CHANGED_MODE) {
        *ppPolicy2 = NULL;
    } else {
        hr = CoCreateInstance(
            &CLSID_NetFwPolicy2,
            NULL,
            CLSCTX_INPROC_SERVER,
            &IID_INetFwPolicy2,
            (void **)&pPolicy2
        );
        if (SUCCEEDED(hr)) {
            *ppPolicy2 = (void *)pPolicy2;
        } else {
            *ppPolicy2 = NULL;
            CoUninitialize();
        }
    }
    return hr;
}

void mpfirewall_release(void *pPolicy2) {
    if (pPolicy2 != NULL) {
        INetFwPolicy2 *p = (INetFwPolicy2 *)pPolicy2;
        p->lpVtbl->Release(p);
    }
    CoUninitialize();
}

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
) {
    HRESULT hr = S_OK;
    INetFwPolicy2 *p = NULL;
    INetFwRules *pRules = NULL;
    INetFwRule *pRule = NULL;
    BSTR bstrName = NULL;
    BSTR bstrDesc = NULL;
    BSTR bstrApp = NULL;
    BSTR bstrPorts = NULL;

    if (pPolicy2 == NULL) {
        hr = E_POINTER;
    } else {
        p = (INetFwPolicy2 *)pPolicy2;
        hr = p->lpVtbl->get_Rules(p, &pRules);
        if (FAILED(hr) || pRules == NULL) {
        } else {
            hr = CoCreateInstance(
                &CLSID_NetFwRule,
                NULL,
                CLSCTX_INPROC_SERVER,
                &IID_INetFwRule,
                (void **)&pRule
            );
            if (FAILED(hr) || pRule == NULL) {
            } else {
                bstrName = utf8_to_bstr(name);
                bstrDesc = utf8_to_bstr(description);
                bstrApp = utf8_to_bstr(application_name);
                bstrPorts = utf8_to_bstr(local_ports);

                if (bstrName != NULL) {
                    pRule->lpVtbl->put_Name(pRule, bstrName);
                    SysFreeString(bstrName);
                }
                if (bstrDesc != NULL) {
                    pRule->lpVtbl->put_Description(pRule, bstrDesc);
                    SysFreeString(bstrDesc);
                }
                if (bstrApp != NULL && application_name != NULL && application_name[0] != '\0') {
                    pRule->lpVtbl->put_ApplicationName(pRule, bstrApp);
                    SysFreeString(bstrApp);
                }

                pRule->lpVtbl->put_Protocol(pRule, (LONG)protocol);
                if (bstrPorts != NULL) {
                    pRule->lpVtbl->put_LocalPorts(pRule, bstrPorts);
                    SysFreeString(bstrPorts);
                }

                pRule->lpVtbl->put_Direction(pRule, (NET_FW_RULE_DIRECTION)direction);
                pRule->lpVtbl->put_Action(pRule, (NET_FW_ACTION)action);
                pRule->lpVtbl->put_Profiles(pRule, (long)profiles);
                pRule->lpVtbl->put_Enabled(pRule, enabled ? VARIANT_TRUE : VARIANT_FALSE);

                hr = pRules->lpVtbl->Add(pRules, pRule);
            }
        }
    }

    if (pRule == NULL) {
    } else {
        pRule->lpVtbl->Release(pRule);
    }
    if (pRules == NULL) {
    } else {
        pRules->lpVtbl->Release(pRules);
    }
    return hr;
}

HRESULT mpfirewall_delete_rule(void *pPolicy2, const char *ruleName) {
    HRESULT hr = S_OK;
    INetFwPolicy2 *p = NULL;
    INetFwRules *pRules = NULL;
    BSTR bstrName = NULL;

    if (pPolicy2 == NULL) {
        hr = E_POINTER;
    } else {
        p = (INetFwPolicy2 *)pPolicy2;
        hr = p->lpVtbl->get_Rules(p, &pRules);
        if (FAILED(hr) || pRules == NULL) {
        } else {
            bstrName = utf8_to_bstr(ruleName);
            if (bstrName == NULL) {
                hr = E_OUTOFMEMORY;
            } else {
                hr = pRules->lpVtbl->Remove(pRules, bstrName);
            }
        }
    }

    if (bstrName == NULL) {
    } else {
        SysFreeString(bstrName);
    }
    if (pRules == NULL) {
    } else {
        pRules->lpVtbl->Release(pRules);
    }
    return hr;
}

typedef void (*rule_callback_t)(INetFwRule *pRule, void *context);

static void enumerate_rules(INetFwPolicy2 *p, rule_callback_t callback, void *context) {
    HRESULT hr = S_OK;
    INetFwRules *pRules = NULL;
    IUnknown *pEnumUnk = NULL;

    if (p == NULL || callback == NULL) {
    } else {
        hr = p->lpVtbl->get_Rules(p, &pRules);
        if (FAILED(hr) || pRules == NULL) {
        } else {
            hr = pRules->lpVtbl->get__NewEnum(pRules, &pEnumUnk);
            if (SUCCEEDED(hr) && pEnumUnk != NULL) {
                IEnumVARIANT *pEnum = NULL;
                hr = pEnumUnk->lpVtbl->QueryInterface(pEnumUnk, &IID_IEnumVARIANT, (void **)&pEnum);
                if (SUCCEEDED(hr) && pEnum != NULL) {
                    VARIANT var;
                    ULONG fetched = 0;
                    while (SUCCEEDED(pEnum->lpVtbl->Next(pEnum, 1, &var, &fetched)) && fetched > 0) {
                        if (V_VT(&var) == VT_DISPATCH && V_DISPATCH(&var) != NULL) {
                            INetFwRule *pRule = NULL;
                            hr = V_DISPATCH(&var)->lpVtbl->QueryInterface(V_DISPATCH(&var), &IID_INetFwRule, (void **)&pRule);
                            if (SUCCEEDED(hr) && pRule != NULL) {
                                callback(pRule, context);
                                pRule->lpVtbl->Release(pRule);
                            }
                        }
                        VariantClear(&var);
                    }
                    pEnum->lpVtbl->Release(pEnum);
                }
                pEnumUnk->lpVtbl->Release(pEnumUnk);
            }
        }
    }

    if (pRules == NULL) {
    } else {
        pRules->lpVtbl->Release(pRules);
    }
}

typedef struct {
    const char *exeFilePath;
    int port;
    char **rule_names;
    int count;
    int capacity;
} BLOCKED_RULES_CONTEXT_T;

static void blocked_rules_callback(INetFwRule *pRule, void *context) {
    BLOCKED_RULES_CONTEXT_T *ctx = (BLOCKED_RULES_CONTEXT_T *)context;
    NET_FW_ACTION action = NET_FW_ACTION_ALLOW;
    BSTR appName = NULL;
    BSTR localPorts = NULL;
    BSTR ruleName = NULL;
    int match = 0;
    char *nameUtf8 = NULL;
    char **newNames = NULL;

    pRule->lpVtbl->get_Action(pRule, &action);
    pRule->lpVtbl->get_ApplicationName(pRule, &appName);
    pRule->lpVtbl->get_LocalPorts(pRule, &localPorts);
    pRule->lpVtbl->get_Name(pRule, &ruleName);

    if (action == NET_FW_ACTION_BLOCK) {
        if (appName != NULL) {
            wchar_t exeW[MAX_PATH];
            MultiByteToWideChar(CP_UTF8, 0, ctx->exeFilePath, -1, exeW, MAX_PATH);
            if (wcsieq(appName, exeW)) {
                match = 1;
            }
        }
        if (!match && localPorts != NULL) {
            char portStr[16];
            snprintf(portStr, sizeof(portStr), "%d", ctx->port);
            wchar_t portW[16];
            MultiByteToWideChar(CP_UTF8, 0, portStr, -1, portW, 16);
            if (wcsstr(localPorts, portW) != NULL) {
                match = 1;
            }
        }
        if (match && ruleName != NULL) {
            nameUtf8 = bstr_to_utf8(ruleName);
            if (nameUtf8 != NULL) {
                if (ctx->count >= ctx->capacity) {
                    ctx->capacity *= 2;
                    newNames = (char **)realloc(ctx->rule_names, ctx->capacity * sizeof(char *));
                    if (newNames != NULL) {
                        ctx->rule_names = newNames;
                    } else {
                        ctx->capacity = ctx->count;
                    }
                }
                if (ctx->count < ctx->capacity) {
                    ctx->rule_names[ctx->count] = nameUtf8;
                    ctx->count++;
                } else {
                    free(nameUtf8);
                }
            }
        }
    }

    if (ruleName == NULL) {
    } else {
        SysFreeString(ruleName);
    }
    if (localPorts == NULL) {
    } else {
        SysFreeString(localPorts);
    }
    if (appName == NULL) {
    } else {
        SysFreeString(appName);
    }
}

MPFIREWALL_RULE_NAME_LIST_T *mpfirewall_get_blocked_rules(void *pPolicy2, const char *exeFilePath, int port) {
    MPFIREWALL_RULE_NAME_LIST_T *result = (MPFIREWALL_RULE_NAME_LIST_T *)calloc(1, sizeof(MPFIREWALL_RULE_NAME_LIST_T));
    BLOCKED_RULES_CONTEXT_T ctx;
    char **names = NULL;

    if (result == NULL || pPolicy2 == NULL || exeFilePath == NULL) {
    } else {
        ctx.exeFilePath = exeFilePath;
        ctx.port = port;
        ctx.rule_names = (char **)calloc(16, sizeof(char *));
        ctx.count = 0;
        ctx.capacity = 16;

        if (ctx.rule_names != NULL) {
            enumerate_rules((INetFwPolicy2 *)pPolicy2, blocked_rules_callback, &ctx);
        }

        result->rule_names = ctx.rule_names;
        result->count = ctx.count;
    }

    return result;
}

void mpfirewall_free_rule_name_list(MPFIREWALL_RULE_NAME_LIST_T *list) {
    if (list == NULL) {
    } else {
        for (int i = 0; i < list->count; i++) {
            free(list->rule_names[i]);
        }
        free(list->rule_names);
        free(list);
    }
}

HRESULT mpfirewall_get_rule_local_port(void *pPolicy2, const char *ruleName, char **outLocalPort) {
    HRESULT hr = S_OK;
    INetFwPolicy2 *p = NULL;
    INetFwRules *pRules = NULL;
    BSTR bstrName = NULL;
    INetFwRule *pRule = NULL;
    BSTR localPorts = NULL;
    char *portUtf8 = NULL;

    if (outLocalPort == NULL) {
        hr = E_POINTER;
    } else {
        *outLocalPort = NULL;
        if (pPolicy2 == NULL || ruleName == NULL) {
            hr = E_POINTER;
        } else {
            p = (INetFwPolicy2 *)pPolicy2;
            hr = p->lpVtbl->get_Rules(p, &pRules);
            if (FAILED(hr) || pRules == NULL) {
            } else {
                bstrName = utf8_to_bstr(ruleName);
                if (bstrName == NULL) {
                    hr = E_OUTOFMEMORY;
                } else {
                    hr = pRules->lpVtbl->Item(pRules, bstrName, &pRule);
                    if (SUCCEEDED(hr) && pRule != NULL) {
                        hr = pRule->lpVtbl->get_LocalPorts(pRule, &localPorts);
                        if (SUCCEEDED(hr) && localPorts != NULL) {
                            portUtf8 = bstr_to_utf8(localPorts);
                            *outLocalPort = portUtf8;
                            SysFreeString(localPorts);
                        }
                    }
                }
            }
        }
    }

    if (bstrName == NULL) {
    } else {
        SysFreeString(bstrName);
    }
    if (pRule == NULL) {
    } else {
        pRule->lpVtbl->Release(pRule);
    }
    if (pRules == NULL) {
    } else {
        pRules->lpVtbl->Release(pRules);
    }
    return hr;
}

typedef struct {
    const char *exeFilePath;
    int port;
    int found;
    char *name;
    char *local_port;
} EXISTS_CONTEXT_T;

static void exists_callback(INetFwRule *pRule, void *context) {
    EXISTS_CONTEXT_T *ctx = (EXISTS_CONTEXT_T *)context;
    BSTR appName = NULL;
    BSTR localPorts = NULL;
    BSTR ruleName = NULL;
    wchar_t exeW[MAX_PATH];
    char portStr[16];
    wchar_t portW[16];

    if (ctx->found) {
    } else {
        pRule->lpVtbl->get_ApplicationName(pRule, &appName);
        pRule->lpVtbl->get_LocalPorts(pRule, &localPorts);

        if (appName == NULL) {
        } else {
            MultiByteToWideChar(CP_UTF8, 0, ctx->exeFilePath, -1, exeW, MAX_PATH);
            if (wcsieq(appName, exeW)) {
                if (localPorts == NULL) {
                } else {
                    snprintf(portStr, sizeof(portStr), "%d", ctx->port);
                    MultiByteToWideChar(CP_UTF8, 0, portStr, -1, portW, 16);
                    if (wcsstr(localPorts, portW) == NULL) {
                    } else {
                        ctx->found = 1;
                        pRule->lpVtbl->get_Name(pRule, &ruleName);
                        if (ruleName == NULL) {
                        } else {
                            ctx->name = bstr_to_utf8(ruleName);
                            SysFreeString(ruleName);
                        }
                        ctx->local_port = bstr_to_utf8(localPorts);
                    }
                }
            }
        }
    }

    if (localPorts == NULL) {
    } else {
        SysFreeString(localPorts);
    }
    if (appName == NULL) {
    } else {
        SysFreeString(appName);
    }
}

MPFIREWALL_RULE_INFO_T *mpfirewall_is_rule_exists(void *pPolicy2, const char *exeFilePath, int port) {
    MPFIREWALL_RULE_INFO_T *result = (MPFIREWALL_RULE_INFO_T *)calloc(1, sizeof(MPFIREWALL_RULE_INFO_T));
    EXISTS_CONTEXT_T ctx;

    if (result == NULL) {
    } else {
        result->found = 0;
        result->name = NULL;
        result->local_port = NULL;

        if (pPolicy2 == NULL || exeFilePath == NULL) {
        } else {
            ctx.exeFilePath = exeFilePath;
            ctx.port = port;
            ctx.found = 0;
            ctx.name = NULL;
            ctx.local_port = NULL;

            enumerate_rules((INetFwPolicy2 *)pPolicy2, exists_callback, &ctx);

            result->found = ctx.found;
            result->name = ctx.name;
            result->local_port = ctx.local_port;
        }
    }

    return result;
}

void mpfirewall_free_rule_info(MPFIREWALL_RULE_INFO_T *info) {
    if (info == NULL) {
    } else {
        free(info->name);
        free(info->local_port);
        free(info);
    }
}
