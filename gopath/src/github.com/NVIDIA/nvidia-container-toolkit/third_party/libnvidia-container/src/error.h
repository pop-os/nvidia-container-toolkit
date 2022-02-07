/*
 * Copyright (c) 2017-2018, NVIDIA CORPORATION. All rights reserved.
 */

#ifndef HEADER_ERROR_H
#define HEADER_ERROR_H

#include <assert.h>
#include <stdalign.h>

#include "error_generic.h"

int error_set_elf(struct error *, const char *, ...)
    __attribute__((format(printf, 2, 3), nonnull(2)));
int error_set_nvml(struct error *, void *, int, const char *, ...)
    __attribute__((format(printf, 4, 5), nonnull(4)));
int error_set_rpc(struct error *, int, const char *, ...)
    __attribute__((format(printf, 3, 4), nonnull(3)));
int error_set_dl(struct error *, const char *, ...)
    __attribute__((format(printf, 2, 3), nonnull(2)));

#include "xfuncs.h"

#define error_from_xdr(err, xdr) do {                                             \
        struct error *xdrerr_ = (struct error *)xdr;                              \
        char *msg_;                                                               \
                                                                                  \
        static_assert(alignof(*err) == alignof(*xdr), "incompatible alignment");  \
        if (xdrerr_->code != 0 && (msg_ = xstrdup(err, xdrerr_->msg)) != NULL) {  \
                (err)->code = xdrerr_->code;                                      \
                (err)->msg = msg_;                                                \
        }                                                                         \
} while (0)

#define error_to_xdr(err, xdr) do {                                               \
        struct error *xdrerr_ = (struct error *)xdr;                              \
        char *msg_;                                                               \
                                                                                  \
        static_assert(alignof(*err) == alignof(*xdr), "incompatible alignment");  \
        if ((err)->code != 0 && (msg_ = xstrdup(xdrerr_, (err)->msg)) != NULL) {  \
                xdrerr_->code = (err)->code;                                      \
                xdrerr_->msg = msg_;                                              \
        }                                                                         \
} while (0)

#endif /* HEADER_ERROR_H */
