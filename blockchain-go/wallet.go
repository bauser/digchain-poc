package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"log"

	"falcon-go/falcon"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	private, public, err := falcon.Keygen()
	if err != nil {
		log.Panic(err)
	}
	wallet := Wallet{private, public}

	return &wallet
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

func PubKeyHashToAddress(pubKeyHash []byte) string {
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return string(address)
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func newKeyPair() ([]byte, []byte) {
	privKey, pubKey, err := falcon.Keygen()
	if err != nil {
		log.Panic(err)
	}

	return privKey, pubKey
}

func (w Wallet) MarshalJSON() ([]byte, error) {
	mapStringAny := map[string]any{
		"PrivateKey": w.PrivateKey,
		"PublicKey":  w.PublicKey,
	}
	return json.Marshal(mapStringAny)
}

func (w *Wallet) UnmarshalJSON(data []byte) error {
	temp := struct {
		PrivateKey []byte
		PublicKey  []byte
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	w.PrivateKey = temp.PrivateKey
	w.PublicKey = temp.PublicKey

	return nil
}
