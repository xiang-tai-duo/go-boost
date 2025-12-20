//go:build darwin
/*
 * File:        hardware_darwin.c
 * Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/hardware/hardware_darwin.c
 * Author:      TRAE.AI
 * Created:     2025/12/20 12:31:58
 * Description: macOS (darwin) native implementations for hardware information,
 *              retrieving platform UUID/serial via IOKit and physical memory via sysctl.
 * --------------------------------------------------------------------------------
 */
#include "hardware_darwin.h"
#include <stddef.h>
#include <stdlib.h>
#include <string.h>
#include <sys/sysctl.h>
#include <sys/types.h>
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>

#define HARDWARE_IOKIT_PLATFORM_NAME "IOPlatformExpertDevice"
#define HARDWARE_SYSCTL_MEMSIZE      "hw.memsize"

static char *copy_registry_string(CFStringRef key);
static unsigned long long fnv1a_hash_string(const char *value);

static char *copy_registry_string(CFStringRef key) {
    char *result = NULL;
    io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching(HARDWARE_IOKIT_PLATFORM_NAME));
    if (service) {
        CFTypeRef value = IORegistryEntryCreateCFProperty(service, key, kCFAllocatorDefault, 0);
        if (value) {
            if (CFGetTypeID(value) == CFStringGetTypeID()) {
                CFIndex length = CFStringGetMaximumSizeForEncoding(CFStringGetLength((CFStringRef)value), kCFStringEncodingUTF8) + 1;
                char *buffer = (char *)calloc((size_t)length, 1);
                if (buffer) {
                    if (CFStringGetCString((CFStringRef)value, buffer, length, kCFStringEncodingUTF8)) {
                        result = buffer;
                    } else {
                        free(buffer);
                    }
                }
            }
            CFRelease(value);
        }
        IOObjectRelease(service);
    }
    return result;
}

static unsigned long long fnv1a_hash_string(const char *value) {
    unsigned long long result = 1469598103934665603ULL;
    if (value) {
        while (*value) {
            result ^= (unsigned char)*value;
            result *= 1099511628211ULL;
            value++;
        }
    }
    return result;
}

unsigned long long hardware_fingerprint_seed(void) {
    unsigned long long result = 1469598103934665603ULL;
    char *uuid = copy_registry_string(CFSTR(kIOPlatformUUIDKey));
    char *serial = copy_registry_string(CFSTR(kIOPlatformSerialNumberKey));
    if (uuid) {
        result = fnv1a_hash_string(uuid);
    }
    if (serial) {
        const char *cursor = serial;
        while (*cursor) {
            result ^= (unsigned char)*cursor;
            result *= 1099511628211ULL;
            cursor++;
        }
    }
    free(serial);
    free(uuid);
    return result;
}

char *hardware_platform_uuid(void) {
    return copy_registry_string(CFSTR(kIOPlatformUUIDKey));
}

char *hardware_platform_serial(void) {
    return copy_registry_string(CFSTR(kIOPlatformSerialNumberKey));
}

unsigned long long hardware_physical_memory(void) {
    unsigned long long result = 0;
    size_t length = sizeof(result);
    if (sysctlbyname(HARDWARE_SYSCTL_MEMSIZE, &result, &length, NULL, 0) != 0) {
        result = 0;
    }
    return result;
}
