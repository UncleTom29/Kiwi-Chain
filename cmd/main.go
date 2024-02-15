package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"fmt"
	"net"
	"log"
	"io"
	"strings"
	"errors"
	"crypto/rsa"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"
	"math/rand"
	"sync"
	"encoding/json"
	"bufio"
	"github.com/perlin-network/life/exec"
	"github.com/cosmos/ibc-go/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/modules/core/keeper"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

var Blockchain []Block
var tempBlocks []Block
var candidateBlocks = make(chan Block)
var announcements = make(chan string)

var validators = make(map[string]int)
var balances = make(map[string]int)
var mutex = &sync.Mutex{}

var Difficulty = 4
const blockReward = 50
const totalSupply = 1000000 // Total supply of tokens
const numShards = 10 // Number of shards

const (
	Transfer ProposalType = iota
	Contract
	ProtocolChange
)

var Proposals []Proposal
var QuorumPercentage = 51 // The percentage of votes needed for a proposal to pass
// Global configuration for the blockchain protocol
var GlobalConfig = map[string]int{
	"blockSize": 10,
	"difficulty": 5,
	"reward":    50,
	"supply":   1000000,
	"shards":   10,
	// Add other parameters as needed
}

// Supported transaction types for the blockchain protocol
var TransactionTypes = map[string]func(Transaction) error{
	"transfer": TransferTransaction,
	"contract": ContractTransaction,
	"custom":   CustomTransaction, // New transaction type
	// Add other transaction types as needed
}

func main() {
	var circuit Circuit

	r1cs, err := frontend.Compile(groth16.SNARK, &circuit)
	if err != nil {
		panic(err)
	}

	{
		// Proving phase
		var witness Circuit
		witness.Sender.Assign(42)
		witness.Receiver.Assign(41)
		witness.Amount.Assign(1)
		witness.SecretKey.Assign(123456789) // This is the private input

		proof, err := groth16.Prove(r1cs, &witness)
		if err != nil {
			panic(err)
		}

		// Verification phase
		var publicWitness Circuit
		publicWitness.Sender.Assign(42)
		publicWitness.Receiver.Assign(41)
		publicWitness.Amount.Assign(1)

		err = groth16.Verify(proof, r1cs, &publicWitness)
		if err != nil {
			panic(err)
		}
	}
}
