/*
 * Copyright (c) 2017-2018, NVIDIA CORPORATION. All rights reserved.
 */

#ifndef HEADER_COMMON_H
#define HEADER_COMMON_H

#include <inttypes.h>

#define PROC_PID                  "/proc/%"PRId32
#define PROC_SELF                 "/proc/self"
#define PROC_MOUNTS_PATH(proc)    proc "/mountinfo"
#define PROC_CGROUP_PATH(proc)    proc "/cgroup"
#define PROC_ROOT_PATH(proc)      proc "/root/" /* XXX Leading slash is required */
#define PROC_NS_PATH(proc)        proc "/ns/%s"
#define PROC_SETGROUPS_PATH(proc) proc "/setgroups"
#define PROC_UID_MAP_PATH(proc)   proc "/uid_map"
#define PROC_GID_MAP_PATH(proc)   proc "/gid_map"
#define PROC_LAST_CAP_PATH        "/proc/sys/kernel/cap_last_cap"
#define PROC_OVERFLOW_UID         "/proc/sys/kernel/overflowuid"
#define PROC_OVERFLOW_GID         "/proc/sys/kernel/overflowgid"

#define LDCACHE_PATH              "/etc/ld.so.cache"
#define LDCONFIG_PATH             "/sbin/ldconfig"
#define LDCONFIG_ALT_PATH         "/sbin/ldconfig.real"

#define LIB_DIR                   "/lib64"
#define USR_BIN_DIR               "/usr/bin"
#define USR_LIB_DIR               "/usr/lib64"
#define USR_LIB32_DIR             "/usr/lib32"
#define USR_LIB32_ALT_DIR         "/usr/lib"

#if defined(__x86_64__)
# define LIB_ARCH                 LD_X8664_LIB64
# define LIB32_ARCH               LD_I386_LIB32
# define USR_LIB_MULTIARCH_DIR    "/usr/lib/x86_64-linux-gnu"
# define USR_LIB32_MULTIARCH_DIR  "/usr/lib/i386-linux-gnu"
# if !defined(__NR_execveat)
#  define __NR_execveat 322
# endif /* !defined(__NR_execveat) */
#elif defined(__powerpc64__)
# define LIB_ARCH                 LD_POWERPC_LIB64
# define LIB32_ARCH               LD_UNKNOWN
# define USR_LIB_MULTIARCH_DIR    "/usr/lib/powerpc64le-linux-gnu"
# define USR_LIB32_MULTIARCH_DIR  "/var/empty"
# if !defined(__NR_execveat)
#  define __NR_execveat 362
# endif /* !defined(__NR_execveat) */
#elif defined(__aarch64__)
# define LIB_ARCH                 LD_AARCH64_LIB64
# define LIB32_ARCH               LD_ARM_LIBHF
# define USR_LIB_MULTIARCH_DIR    "/usr/lib/aarch64-linux-gnu/"
# define USR_LIB32_MULTIARCH_DIR  "/var/empty"
# if !defined(__NR_execveat)
#  define __NR_execveat 387
# endif /* !defined(__NR_execveat) */
#else
# error "unsupported architecture"
#endif /* defined(__x86_64__) */

#endif /* HEADER_COMMON_H */
