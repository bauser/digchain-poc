#include <stdlib.h>
#include <stdio.h>
#include <time.h>
#include "falcon.h"
#include "falcon_dig.h"

#define SEED_SIZE 16 // Define the size of the seed in bytes

void generate_random_seed(unsigned char *seed, size_t size) {
    // Initialize the random number generator
    srand((unsigned int)time(NULL));

    // Fill the seed array with random bytes
    for (size_t i = 0; i < size; ++i) {
        seed[i] = rand() % 256; // Generate a random byte (0-255)
    }
}

int dig_falcon_keygen(
	void *privkey,
	void *pubkey
) {
	unsigned logn = 10;
    size_t privkey_len, pubkey_len;
    size_t tmp_len = FALCON_TMPSIZE_SIGNDYN(logn);


	pubkey_len = FALCON_PUBKEY_SIZE(logn);
	privkey_len = FALCON_PRIVKEY_SIZE(logn);

    uint8_t seed[SEED_SIZE];
    generate_random_seed(seed, SEED_SIZE);

    shake256_context sc;
    shake256_init(&sc);
    shake256_inject(&sc, seed, SEED_SIZE);
    shake256_flip(&sc);

    unsigned char *tmp = (unsigned char *)calloc(tmp_len, sizeof(unsigned char));

    privkey_len = FALCON_PRIVKEY_SIZE(logn);
    privkey_len = FALCON_PRIVKEY_SIZE(logn);

    int ret = falcon_keygen_make(&sc, logn, privkey, privkey_len, pubkey, pubkey_len, tmp, tmp_len);
    free(tmp);

    return ret;
}

int dig_falcon_sign(
    void *sig, size_t *sig_len,
    const void *privkey, size_t privkey_len,
    const void *data, size_t data_len
) {
    int logn = 10;
    size_t tmp_len = FALCON_TMPSIZE_SIGNDYN(logn);
    unsigned char *tmp = (unsigned char *)calloc(tmp_len, sizeof(unsigned char));
    
    shake256_context sc;
    shake256_init(&sc);
    shake256_inject(&sc, (const uint8_t *)"digshakesource", 14);
    shake256_flip(&sc);

    int ret = falcon_sign_dyn(&sc, sig, sig_len, FALCON_SIG_PADDED,
        privkey, privkey_len, data, data_len, tmp, tmp_len);
    free(tmp);

    return ret;
}

int dig_falcon_verify(
    void *sig, size_t sig_len,
    const void *pubkey, size_t pubkey_len,
    const void *data, size_t data_len
) {
    int logn = 10;
    size_t tmp_len = FALCON_TMPSIZE_SIGNDYN(logn);
    unsigned char *tmp = (unsigned char *)calloc(tmp_len, sizeof(unsigned char));

    int ret = falcon_verify(sig, sig_len, FALCON_SIG_PADDED,
        pubkey, pubkey_len, data, data_len, tmp, tmp_len);
    free(tmp);

    return ret;
}

int dig_falcon_make_public(
    void *pubkey, size_t pubkey_len,
    const void *privkey, size_t privkey_len
) {
    int logn = 10;
    size_t tmp_len = FALCON_TMPSIZE_MAKEPUB(logn);
    unsigned char *tmp = (unsigned char *)calloc(tmp_len, sizeof(unsigned char));

    int ret = falcon_make_public(pubkey, pubkey_len,
        privkey, privkey_len,
        tmp, tmp_len);
    free(tmp);

    return ret;
}

// int main() {
//     int logn = 10;

//     size_t privkey_len = FALCON_PRIVKEY_SIZE(logn);
//     size_t pubkey_len = FALCON_PUBKEY_SIZE(logn);
//     void *privkey = malloc(privkey_len);
//     void *pubkey = malloc(pubkey_len);

//     size_t *sig_len = malloc(1);
//     *sig_len = FALCON_SIG_PADDED_SIZE(logn);
//     void *sig = malloc(*sig_len);

//     printf("\nsig len: %d\n", (int)(*sig_len));
//     printf("\npk len: %d\n", (int)privkey_len);
//     printf("\npb len: %d\n", (int)pubkey_len);

//     printf("\n\nkeygen ret: %d\n\n", dig_falcon_keygen(privkey, pubkey));
//     printf("\n\nsign ret: %d\n\n", dig_falcon_sign(sig, sig_len, privkey, privkey_len, "Hello, Falcon!!!", 16));
//     printf("\n\nverify ret: %d\n\n", dig_falcon_verify(sig, *sig_len, pubkey, pubkey_len, "Hello, Falcon!!!", 16));

//     free(privkey);
//     free(pubkey);
//     free(sig);
//     free(sig_len);
//     // printf("privkey: ");
//     // for (int i = 0; i < privkey_len; i++) {
//     //     printf("%x", ((uint8_t*)privkey)[i]);
//     // }
//     // printf("\n\n");
//     // printf("pubkey: ");
//     // for (int i = 0; i < pubkey_len; i++) {
//     //     printf("%x", ((uint8_t*)pubkey)[i]);
//     // }
//     // printf("\n\n");

//     return 0;
// }
