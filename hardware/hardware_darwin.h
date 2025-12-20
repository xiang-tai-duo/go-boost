/*
 * File:        hardware_darwin.h
 * Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/hardware/hardware_darwin.h
 * Author:      TRAE.AI
 * Created:     2025/12/20 12:31:58
 * Description: Header for macOS (darwin) native hardware information helpers.
 * --------------------------------------------------------------------------------
 */
#ifndef GO_BOOST_HARDWARE_DARWIN_H
#define GO_BOOST_HARDWARE_DARWIN_H

unsigned long long hardware_fingerprint_seed(void);
char *hardware_platform_serial(void);
char *hardware_platform_uuid(void);
unsigned long long hardware_physical_memory(void);

#endif
