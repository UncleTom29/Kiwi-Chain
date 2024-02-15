package main

import (
	"crypto/sha256"
	"crypto/rsa"
	"crypto/rand"
	"log"
	"encoding/hex"
	"encoding/pem"
	"crypto/x509"
)

type Transaction struct {
	From      string
	To        string
	Amount    int
	Fee       int
	Signature string
	Contract  *SmartContract
}

type SmartContract struct {
	Code []byte
	Data map[string]interface{}
}

func isTransactionValid(tx Transaction) bool {
	// Check if the signature is valid
	if !isValidSignature(tx) {
		return false
	}

	// Check if the sender has enough balance for the transaction
	if !hasEnoughBalance(tx) {
		return false
	}

	// Check if the transaction has not been processed before
	if isDoubleSpending(tx) {
		return false
	}

	// Execute the smart contract (if any)
	if tx.Contract != nil {
		vm, err := exec.NewVirtualMachine(tx.Contract.Code, exec.VMConfig{}, nil, nil)
		if err != nil {
			log.Println(err)
			return false
		}

		_, err = vm.Run(tx.Contract.Code)
		if err != nil {
			log.Println(err)
			return false
		}
	}

	return true
}

func isValidSignature(tx Transaction) bool {
	// Convert the public key from string to *rsa.PublicKey
	pubKey := convertPublicKey(tx.From)

	// Create a new hash
	hash := sha256.New()

	// Write the transaction data to the hash
	_, err := hash.Write([]byte(tx.Data)) // assuming tx.Data contains the transaction data
	if err != nil {
		log.Println(err)
		return false
	}
	hashed := hash.Sum(nil)

	// Decode the signature
	signature, err := hex.DecodeString(tx.Signature)
	if err != nil {
		log.Println(err)
		return false
	}

	// Verify the signature
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func hasEnoughBalance(tx Transaction) bool {
	// Check the sender's balance
	balance := balances[tx.From]
	// Check if the sender's balance is less than the amount
	if balance < tx.Amount {
		return false
	}
	return true
}

func isDoubleSpending(tx Transaction) bool {
	for _, block := range Blockchain {
		for _, transaction := range block.Transactions {
			// Compare transaction details to check for identical transactions
			if tx.From == transaction.From && tx.To == transaction.To && tx.Amount == transaction.Amount {
				return true
			}
		}
	}
	return false
}


func convertPublicKey(pubKeyStr string) *rsa.PublicKey {
	// Decode the public key from PEM format
	block, _ := pem.Decode([]byte(pubKeyStr))
	if block == nil {
		return nil
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Type assert to rsa.PublicKey
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil
	}

	return rsaPub
}
