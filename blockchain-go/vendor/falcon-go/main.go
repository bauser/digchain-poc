package main

import (
	"falcon-go/falcon"
	"fmt"
)

func main() {
	// Generate key pair
	prvk, pubk, err := falcon.Keygen()
	if err != nil {
		fmt.Println("Keygen error:", err)
		return
	}

	// Message to be signed
	message := []byte("Hello, Falcon!")

	// Sign the message
	sig, err := falcon.Sign(message, prvk)
	if err != nil {
		fmt.Println("Sign error:", err)
		return
	}

	fmt.Println("Sig length:", len(sig))
	fmt.Println("PrvK length:", len(prvk))
	fmt.Println("PubK length:", len(pubk))

	// Verify the signature
	if !falcon.Verify(message, sig, pubk) {
		fmt.Println("Verify failed")
	} else {
		fmt.Println("Signature verified successfully!")
	}
}
