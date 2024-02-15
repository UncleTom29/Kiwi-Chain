package main

import (
	"encoding/json"
)

type Proposal struct {
	ID          string
	Description string
	Changes     string
	Votes       map[string]bool // Maps node addresses to their votes
	Delegations map[string]string // Maps node addresses to the addresses of nodes they've delegated their vote to
}

type ProtocolChange struct {
	ParameterUpdates map[string]int
	NewFeatures      []string
}

type ProposalType int

func NewProposal(node Node, description string, changes string, proposalType ProposalType) Proposal {
	proposal := Proposal{
		ID:          node.ID + "-" + time.Now().String(),
		Description: description,
		Changes:     changes,
		Votes:       make(map[string]bool),
		Delegations: make(map[string]string),
	}

	// The node that created the proposal automatically votes for it
	proposal.Votes[node.Address] = true

	return proposal
}

func DelegateVote(node Node, delegate Node, proposalID string) error {
	for i, proposal := range Proposals {
		if proposal.ID == proposalID {
			_, delegated := proposal.Delegations[node.Address]
			if delegated {
				return errors.New("node has already delegated their vote on this proposal")
			}

			Proposals[i].Delegations[node.Address] = delegate.Address
			return nil
		}
	}

	return errors.New("proposal not found")
}

func Vote(node Node, proposalID string, vote bool) error {
	for i, proposal := range Proposals {
		if proposal.ID == proposalID {
			_, voted := proposal.Votes[node.Address]
			if voted {
				return errors.New("node has already voted on this proposal")
			}

			Proposals[i].Votes[node.Address] = vote
			return nil
		}
	}

	return errors.New("proposal not found")
}

func TallyVotes(proposalID string) (bool, error) {
	for _, proposal := range Proposals {
		if proposal.ID == proposalID {
			yesVotes := 0
			noVotes := 0

			for node, vote := range proposal.Votes {
				// If the node has delegated their vote, use the vote of the delegate
				if delegate, ok := proposal.Delegations[node]; ok {
					vote = proposal.Votes[delegate]
				}

				if vote {
					yesVotes++
				} else {
					noVotes++
				}
			}

			// Check if the proposal has reached quorum
			totalVotes := yesVotes + noVotes
			quorum := (totalVotes * QuorumPercentage) / 100
			if yesVotes < quorum {
				return false, errors.New("proposal did not reach quorum")
			}

			return yesVotes > noVotes, nil
		}
	}

	return false, errors.New("proposal not found")
}


func executeProposalApplicationLogic(ctx sdk.Context, k keeper.Keeper, proposal Proposal) error {
	// Parse the proposal data
	var data ProposalData
	err := json.Unmarshal(proposal.GetData(), &data)
	if err != nil {
		return err
	}

	switch data.Type {
	case Transfer:
		// Handle token transfer
		err = handleTokenTransferProposal(ctx, k, data, proposal)
		if err != nil {
			return err
		}
	case Contract:
		// Handle smart contract execution
		err = handleSmartContractExecutionProposal(ctx, k, data, proposal)
		if err != nil {
			return err
		}
	case ProtocolChange:
		// Handle protocol change
		err = handleProtocolChangeProposal(ctx, k, data, proposal)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown proposal type")
	}

	return nil
}

func handleTokenTransferProposal(ctx sdk.Context, k keeper.Keeper, data ProposalData, proposal Proposal) error {
	// Extract the sender, receiver and amount from the proposal data
	sender := data.Sender
	receiver := data.Receiver
	amount := data.Amount

	// Deduct the tokens from the sender's account
	balances[sender] -= amount

	// Add the tokens to the receiver's account
	balances[receiver] += amount

	return nil
}

func handleSmartContractExecutionProposal(ctx sdk.Context, k keeper.Keeper, data ProposalData, proposal Proposal) error {
	// Extract the contract and parameters from the proposal data
	contract := data.Contract
	parameters := data.Parameters

	// Create a new instance of the smart contract
	sc := SmartContract{
		Code: contract,
		Data: parameters,
	}

	// Execute the smart contract
	vm, err := exec.NewVirtualMachine(sc.Code, exec.VMConfig{}, nil, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = vm.Run(sc.Code)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func handleProtocolChangeProposal(ctx sdk.Context, k keeper.Keeper, data ProposalData, proposal Proposal) error {
	// Extract the proposed changes from the proposal data
	var changes ProtocolChange
	err := json.Unmarshal([]byte(data.Changes), &changes)
	if err != nil {
		return err
	}

	// Apply the parameter updates
	for parameter, value := range changes.ParameterUpdates {
		// The specific logic will depend on your protocol and the nature of the parameters
		// For example, you might update a global variable, modify a configuration file, etc.
		updateParameter(parameter, value)
	}

	// Add the new features
	for _, feature := range changes.NewFeatures {
		// The specific logic will depend on your protocol and the nature of the features
		// For example, you might register a new transaction type, add a new API endpoint, etc.
		addFeature(feature)
	}

	return nil
}


func updateParameter(parameter string, value int) error {
	// Validate the parameter
	if _, ok := GlobalConfig[parameter]; !ok {
		return fmt.Errorf("invalid parameter: %s", parameter)
	}

	// Validate the value
	if value < 0 {
		return errors.New("value must be non-negative")
	}

	// Update the parameter in the global configuration
	GlobalConfig[parameter] = value
	fmt.Printf("Updated parameter %s to %d\n", parameter, value)

	return nil
}

func addFeature(feature string) error {
	// Validate the feature
	if _, ok := TransactionTypes[feature]; ok {
		return fmt.Errorf("feature already exists: %s", feature)
	}

	// Add the feature to the supported transaction types
	switch feature {
	case "custom":
		TransactionTypes[feature] = CustomTransaction
	default:
		return fmt.Errorf("unknown feature: %s", feature)
	}

	fmt.Printf("Added feature %s\n", feature)

	return nil
}

func CustomTransaction(tx Transaction) error {
	// Extract the operation type from the transaction data
	operation, ok := tx.Data["operation"].(string)
	if !ok {
		return errors.New("invalid operation type")
	}

	// Extract the amount from the transaction data
	amount, ok := tx.Data["amount"].(int)
	if !ok {
		return errors.New("invalid amount")
	}

	// Handle the operation
	switch operation {
	case "mint":
		// Mint tokens to the sender's account
		balances[tx.From] += amount
		fmt.Printf("Minted %d tokens to %s\n", amount, tx.From)
	case "burn":
		// Check if the sender has enough balance to burn
		if balances[tx.From] < amount {
			return errors.New("insufficient balance to burn")
		}

		// Burn tokens from the sender's account
		balances[tx.From] -= amount
		fmt.Printf("Burned %d tokens from %s\n", amount, tx.From)
	default:
		return fmt.Errorf("unknown operation: %s", operation)
	}

	return nil
}


func handleProtocolChangeProposal(ctx sdk.Context, k keeper.Keeper, data ProposalData, proposal Proposal) error {
	// Extract the proposed changes from the proposal data
	var changes ProtocolChange
	err := json.Unmarshal([]byte(data.Changes), &changes)
	if err != nil {
		return err
	}

	// Start a database transaction to ensure atomic changes
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Apply the parameter updates
	for parameter, value := range changes.ParameterUpdates {
		err = updateParameter(parameter, value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Add the new features
	for _, feature := range changes.NewFeatures {
		err = addFeature(feature)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the database transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
