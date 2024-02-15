# Kiwi-Chain
 A Go-based blockchain implementation featuring block creation, transaction handling, wallet handling, interoperability implementations, smart contracts, and consensus algorithms. Includes automated deployment with Terraform and Ansible, and performance monitoring with Prometheus and Grafana.


# ⚠️ Disclaimer

This code is incomplete and experimental. It is not fully tested and should not be used in a production environment. Use at your own risk.

# Blockchain Functionalities

The blockchain implementation includes the following functionalities:

- **Block Creation**: Each block in the blockchain contains a set of transactions, a timestamp, a reference to the previous block (PrevHash), and a unique identifier (Hash). Blocks are created and added to the blockchain as part of the consensus algorithm.

- **Transaction Handling**: Transactions represent the actions performed in the blockchain network. This could be the transfer of tokens from one participant to another or the execution of a smart contract. Each transaction is validated and processed by the network nodes.

- **Smart Contracts**: Smart contracts are self-executing contracts with the terms of the agreement directly written into code. They can be deployed to the blockchain and executed as part of a transaction. The smart contract code is run by every node in the blockchain, and the results of the execution are recorded on the blockchain.

- **Consensus Algorithm**: The consensus algorithm is used to agree on the validity of transactions and the order of blocks in the blockchain. The consensus algorithm used is Proof of Activity (POA). This implementation allows for pluggable consensus algorithms like Proof of Work (PoW) or Proof of Stake (PoS).

- **IBC (Inter-Blockchain Communication)**: The blockchain supports inter-blockchain communication, allowing it to interact with other blockchains. This is achieved through a series of transactions and smart contracts that enable the secure transfer of tokens between different blockchains.

## Project Structure

The project is structured into several parts:

```

/kiwi-chain
|-- /cmd
|   |-- main.go
|
|-- /pkg
|   |-- /block
|   |   |-- block.go
|   |
|   |-- /transaction
|   |   |-- transaction.go
|   |
|   |-- /wallet
|   |   |-- wallet.go
|   |
|   |-- /node
|   |   |-- node.go
|   |
|   |-- /ibc
|   |   |-- ibc.go
|   |
|   |-- /proposal
|   |   |-- proposal.go
|
|-- /deploy
|   |-- /terraform
|   |   |-- main.tf
|   |
|   |-- /ansible
|   |   |-- playbook.yml
|
|-- /monitoring
|   |-- /prometheus
|   |   |-- prometheus.yml
|   |
|   |-- /grafana
|   |   |-- dashboard.json
|
|-- /scripts
|   |-- start.sh
|   |-- stop.sh
|
|-- Dockerfile
|-- .gitignore
|-- README.md

```

# Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

# Prerequisites

- Go 1.16 or later
- Docker (optional)
- Terraform (for deployment)
- Ansible (for deployment)

# Installing

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/blockchain-project.git

```
# Navigate to the project directory:

```
cd blockchain-project

```

# Build the project:

```
go build -o main .

```

# Running the Application
You can start the application with:
```
./main

```

# Running in Docker
Build the Docker image:

```
docker build -t kiwi-chain .

```

Run the Docker image:
```
docker run -p 8080:8080 kiwi-chain

```

