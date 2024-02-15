package main

import (
	"crypto/rsa"
	"crypto/rand"
	"log"
	"encoding/hex"
	"encoding/pem"
	"crypto/x509"
)

type Wallet struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewWallet() *Wallet {
	privKey, pubKey := generateKeyPair()
	return &Wallet{PrivateKey: privKey, PublicKey: pubKey}
}

func (w *Wallet) CreateTransaction(to string, amount int) Transaction {
	tx := Transaction{
		From:   publicKeyToString(w.PublicKey),
		To:     to,
		Amount: amount,
	}

	tx.Signature = w.signTransaction(tx)

	return tx
}

func (w *Wallet) signTransaction(tx Transaction) string {
	// Create a new hash
	hash := sha256.New()

	// Write the transaction data to the hash
	_, err := hash.Write([]byte(tx.Data)) // assuming tx.Data contains the transaction data
	if err != nil {
		log.Println(err)
		return ""
	}
	hashed := hash.Sum(nil)

	// Sign the hash with the private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, w.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		log.Println(err)
		return ""
	}

	return hex.EncodeToString(signature)
}

func generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return privkey, &privkey.PublicKey
}

func publicKeyToString(pubKey *rsa.PublicKey) string {
	pubASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Println(err)
		return ""
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return string(pubBytes)
}
