#ifndef DIG_FALCON_H__
#define DIG_FALCON_H__

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

int dig_falcon_keygen(
	void *privkey,
	void *pubkey
);

int dig_falcon_sign(
    void *sig, size_t *sig_len,
    const void *privkey, size_t privkey_len,
    const void *data, size_t data_len
);

int dig_falcon_verify(
    void *sig, size_t sig_len,
    const void *pubkey, size_t pubkey_len,
    const void *data, size_t data_len
);

int dig_falcon_make_public(
    void *pubkey, size_t pubkey_len,
    const void *privkey, size_t privkey_len
);

#ifdef __cplusplus
}
#endif

#endif
