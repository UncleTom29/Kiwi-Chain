package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"fmt"
	"math/rand"
)

type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
	Nonce        string
	Reward       int
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + fmt.Sprintf("%v", block.Transactions) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func createBlock(oldBlock Block, Transactions []Transaction, validator string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transactions = Transactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	newBlock.Reward = blockReward

	// The validator gets a reward
	balances[validator] += blockReward

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	for _, tx := range newBlock.Transactions {
		if !isTransactionValid(tx) {
			return false
		}
	}

	return true
}

func proofOfWork(block Block) Block {
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		block.Nonce = hex
		if !isHashValid(calculateHash(block), Difficulty) {
			fmt.Println(calculateHash(block), " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(block), " work done!")
			block.Hash = calculateHash(block)
			break
		}

	}
	return block
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func proofOfStake(block Block, stakeholder map[string]int) Block {
	rand.Seed(time.Now().UnixNano())
	stakeholders := make([]string, 0, len(stakeholder))
	for s := range stakeholder {
		stakeholders = append(stakeholders, s)
	}

	for {
		chosen := stakeholders[rand.Intn(len(stakeholders))]
		if rand.Intn(totalSupply) < stakeholder[chosen] {
			block.Data = chosen + " " + block.Data
			block.Hash = calculateHash(block)
			break
		}
		time.Sleep(time.Second)
	}
	return block
}
