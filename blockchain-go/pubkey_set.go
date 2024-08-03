package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const pubkeyBucket = "chainstate"

type PubKeySet struct {
	Blockchain *Blockchain
}

func (p PubKeySet) FindPubKeyOfAddr(addr []byte) ([]byte, error) {
	db := p.Blockchain.db

	pubKey := []byte{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pubkeyBucket))
		pubKey = b.Get(addr)

		return nil
	})
	if err != nil {
		return nil, nil
	} else {
		return pubKey, nil
	}
}

func (u PubKeySet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(pubkeyBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	addrBook := u.Blockchain.FindAddressBook()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for addr, pk := range addrBook {
			key := []byte(addr)

			err = b.Put(key, pk)
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

func (u PubKeySet) Update(block *Block) {
	bc := u.Blockchain
	db := bc.db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pubkeyBucket))

		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == true {
				continue
			}

			for _, vin := range tx.Vin {
				if vin.PubKey == nil {
					continue
				}

				addr := PubKeyHashToAddress(vin.PubKeyHash)
				pk, _ := u.FindPubKeyOfAddr([]byte(addr))
				if pk == nil {
					b.Put([]byte(addr), vin.PubKey)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
