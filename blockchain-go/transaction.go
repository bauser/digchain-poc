package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"falcon-go/falcon"
	"math/big"
	"strings"
	"time"

	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey []byte, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	fmt.Printf("Signing tx\n")
	stime := time.Now()
	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKeyHash = prevTx.Vout[vin.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)
		fmt.Printf("[Sign] dataToSign: %s\n\n", txCopy)

		signature, err := falcon.Sign([]byte(dataToSign), privKey)
		if err != nil {
			log.Panic(err)
		}

		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil

		fmt.Printf("[Sign] signature: %x\n\n", signature)
	}
	etime := time.Now()

	duration := etime.Sub(stime)
	fmt.Printf("Signed tx %x\n\n", tx.Hash())
	fmt.Printf("Time taken to sign %v\n\n", duration)
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKeyHash:    %x", input.PubKeyHash))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, vin.PubKeyHash, vin.PubKey})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Verify verifies signatures of Transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction, bc *Blockchain) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	fmt.Printf("Verifying tx %x\n", tx.Hash())
	stime := time.Now()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		vout := prevTx.Vout[vin.Vout]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKeyHash = vout.PubKeyHash
		txCopy.Vin[inID].PubKey = vin.PubKey

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		pubKey := make([]byte, 1793)
		pubKeySet := PubKeySet{bc}
		pubKeyFound, _ := pubKeySet.FindPubKeyOfAddr(vout.PubKeyHash)

		if vin.PubKey == nil {
			if pubKeyFound == nil {
				log.Println("Tx pubkey not found")
				return false
			} else {
				copy(pubKey, pubKeyFound)
			}
		} else if pubKeyFound != nil {
			if bytes.Compare(pubKeyFound, vin.PubKey) == 0 {
				copy(pubKey, vin.PubKey)
			} else {
				log.Println("Tx pubkey not the correct one")
				return false
			}
		} else {
			copy(pubKey, vin.PubKey)
		}

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		if !falcon.Verify([]byte(dataToVerify), vin.Signature, pubKey) {
			log.Println("Falcon verification failed")
			return false
		}
		txCopy.Vin[inID].PubKey = nil
	}
	etime := time.Now()

	duration := etime.Sub(stime)
	fmt.Printf("Verified tx %x\n\n", tx.Hash())
	fmt.Printf("Time taken to verify %v\n\n", duration)

	now := time.Now()
	fmt.Println("Time: ", now)

	return true
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{[]byte{}, -1, nil, HashPubKey([]byte(data)), []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(wallet *Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		pubKeySet := PubKeySet{UTXOSet.Blockchain}
		for _, out := range outs {
			pubKey, _ := pubKeySet.FindPubKeyOfAddr(wallet.GetAddress())
			var input TXInput
			if pubKey == nil {
				input = TXInput{txID, out, nil, HashPubKey(wallet.PublicKey), wallet.PublicKey}
			} else {
				input = TXInput{txID, out, nil, HashPubKey(wallet.PublicKey), nil}
			}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", wallet.GetAddress())
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

// DeserializeTransaction deserializes a transaction
func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}
