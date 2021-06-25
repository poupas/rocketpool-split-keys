package main

import (
	"fmt"
	"os"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/poupas/rocketpool-split-keys/threshold"
)

type MinipoolStatus int

const (
	// This value will change as members enter and leave the ODAO
	ODAOMembersCount = 10
	// Minimum number of ODAO members required for signing
	ODAOKeySharesThreshold = 3

	// Minipool statuses
	Initialized = iota
	Prelaunch
	Staking
)

// "Deployed" minipool contracts
var minipoolContracts map[string]*MinipoolContract

type Minipool struct {
	// Smart contract
	contract MinipoolContract
	// Validator secret key
	validatorSecretKey bls.SecretKey
	// How many shares to split the master key into
	keySharesCount uint64
	// The minimum required shares for a signature to be valid
	keySharesThreshold uint64
}

type MinipoolContract struct {
	address         string
	validatorPubKey string
	status          MinipoolStatus
}

// Mock acceptValidatorKey method. This method would be called by the ODAO to signal that this
// minipool key shares have been properly distributed and verified
func (mp *MinipoolContract) acceptValidatorKey() {
	// Check if the caller is the ODAO
	// Assuming that all other checks have passed. E.g. node has deposited
	mp.status = Prelaunch
}

// Mock stake method. Starts staking IFF the ODAO has has validated the key shares
func (mp *MinipoolContract) Stake() bool {
	if mp.status != Prelaunch {
		fmt.Println("Can only start staking if in prelaunch!")
		return false
	}
	return true
}

type ODAOMember struct {
	// ID
	id uint64
	// Validator key share for each minipool
	keyShares map[string]*bls.SecretKey

	// TODO: encrypt key shares sent to validators with a modern encryption algorithm. e.g. X25519
	// pubKey X25519PubKey
	// secretKey X25519PrivKey
}

func (m *ODAOMember) setKeyShare(minipool string, share *bls.SecretKey) {
	m.keyShares[minipool] = share
}

type ODAO struct {
	members []ODAOMember
}

// Distribute the encrypted minipool validator key shares to ODAO members
// In a real setting, this distribution would be done over the network
func (odao *ODAO) distributeKeyShares(minipoolAddress string, shares map[uint64]*bls.SecretKey) {
	for _, member := range odao.members {
		fmt.Printf("Sending minipool '%s' share to ODAO member '%d': %s\n",
			minipoolAddress, member.id, shares[member.id].GetHexString())
		member.setKeyShare(minipoolAddress, shares[member.id])
	}
}

// Check if the validator key shares for a given minipool match the key set in the minipool contract
// In a real setting, ODAO members would coordinate among themselves to perform the verification
// Members would create sub-groups of ODAOKeySharesThreshold elements to ensure
// that the aggregated public key matches the key specified in the minipool contract
// For the purpose of this PoC, just check if using the first ODAOKeySharesThreshold keys
func (odao *ODAO) verifyKeyShares(minipoolAddress string) {
	// Gather the required key shares to recover the public key
	keyShares := make(map[uint64]*bls.SecretKey)
	for i, member := range odao.members {
		keyShares[member.id] = member.keyShares[minipoolAddress]
		if i == ODAOKeySharesThreshold-1 {
			break
		}
	}
	aggregatePubKey, err := threshold.ReconstructPublicKey(keyShares)
	if err != nil {
		fmt.Printf("Could not reconstruct public key: %s\n", err)
		return
	}

	// Make sure that the recovered public key matches the one in the minipool contract
	mp := minipoolContracts[minipoolAddress]
	if aggregatePubKey.GetHexString() != mp.validatorPubKey {
		fmt.Printf("Unexpected validator public key. Recovered key: %s, Contract key: %s\n",
			aggregatePubKey.GetHexString(), mp.validatorPubKey)
		return
	}

	// Everything checks out. Allow the minipool to start staking
	fmt.Printf("Successfully verified public key.\nRecovered public key: %s\nContract key: %s\n",
		aggregatePubKey.GetHexString(), mp.validatorPubKey)
	mp.acceptValidatorKey()
}

func newMinipool(address string) *Minipool {
	// Generate a random validator secret key
	// In an actual deployment, this key would be derived from the master seed phrase
	validatorSecretKey := bls.SecretKey{}
	validatorSecretKey.SetByCSPRNG()

	// Create the minipool
	contract := MinipoolContract{
		address:         address,
		validatorPubKey: validatorSecretKey.GetPublicKey().GetHexString(),
	}
	return &Minipool{
		validatorSecretKey: validatorSecretKey,
		keySharesCount:     ODAOMembersCount,
		keySharesThreshold: ODAOKeySharesThreshold,
		contract:           contract,
	}
}

func (mp *Minipool) SplitValidatorKey() (map[uint64]*bls.SecretKey, error) {
	return threshold.Create(mp.validatorSecretKey.Serialize(), mp.keySharesThreshold, mp.keySharesCount)
}

func main() {
	// Initialize BLS library
	threshold.Init()

	// Initialize the minipool contract data container
	minipoolContracts = make(map[string]*MinipoolContract)

	// Create the ODAO
	odao := ODAO{}
	// Add its initial members
	for i := uint64(0); i < ODAOMembersCount; i++ {
		member := ODAOMember{
			id:        i + 1,
			keyShares: make(map[string]*bls.SecretKey),
		}
		odao.members = append(odao.members, member)
	}

	// Create a minipool
	mp := newMinipool("0xdeadbeef")
	// "Deploy" the minipool smart contract
	minipoolContracts[mp.contract.address] = &mp.contract
	fmt.Printf("Created minipool. Address: %s, Validator pubkey: %s",
		mp.contract.address, mp.contract.validatorPubKey)

	// Share the minipool validator key among the ODAO members
	keyShares, err := mp.SplitValidatorKey()
	if err != nil {
		fmt.Printf("Could not split minipool validator key: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Sending key shares to the ODAO...")
	odao.distributeKeyShares(mp.contract.address, keyShares)

	// Check that the minipool sent the correct key shares to the ODAO
	odao.verifyKeyShares(mp.contract.address)

	// Start staking
	if mp.contract.Stake() {
		fmt.Printf("Sucessfully started staking on minipool '%s'...\n", mp.contract.address)
	}
}
