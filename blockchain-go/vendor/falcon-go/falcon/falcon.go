package falcon

/*
#cgo CFLAGS: -I../falcon-c
#cgo LDFLAGS: -L../falcon-c/build -ldigfalcon
#include "falcon_dig.h"  // Adjust the header file name according to the Falcon C implementation
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func Keygen() ([]byte, []byte, error) {
	var prvk [2305]byte
	var pubk [1793]byte
	res := C.dig_falcon_keygen(unsafe.Pointer(&prvk[0]), unsafe.Pointer(&pubk[0]))
	if res != 0 {
		return nil, nil, fmt.Errorf("falcon_keygen failed with error code %d", res)
	}
	return prvk[:], pubk[:], nil
}

func Sign(message []byte, prvk []byte) ([]byte, error) {
	var sig [1280]byte
	// var sig [3000]byte
	var sig_len int = len(sig)
	res := C.dig_falcon_sign(unsafe.Pointer(&sig[0]), (*C.size_t)(unsafe.Pointer(&sig_len)),
		unsafe.Pointer(&prvk[0]), C.size_t(len(prvk)),
		unsafe.Pointer(&message[0]), C.size_t(len(message)))
	if res != 0 {
		return nil, fmt.Errorf("falcon_sign failed with error code %d", res)
	}
	return sig[:], nil
}

func Verify(message []byte, sig []byte, pubk []byte) bool {
	res := C.dig_falcon_verify(unsafe.Pointer(&sig[0]), C.size_t(len(sig)),
		unsafe.Pointer(&pubk[0]), C.size_t(len(pubk)),
		unsafe.Pointer(&message[0]), C.size_t(len(message)))

	if res != 0 {
		fmt.Printf("falcon_verify result is %d\n", res)
	}

	return res == 0
}

func MakePubkey(prvk []byte) bool {
	var pubk [1793]byte
	res := C.dig_falcon_make_public(unsafe.Pointer(&pubk[0]), C.size_t(len(pubk)),
		unsafe.Pointer(&prvk[0]), C.size_t(len(prvk)))

	return res == 0
}
