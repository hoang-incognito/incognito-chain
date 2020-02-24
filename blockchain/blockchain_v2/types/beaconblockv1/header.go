package beaconblockv1

import (
	"fmt"

	"github.com/incognitochain/incognito-chain/common"
)

type BeaconHeader struct {
	Version           int         `json:"Version"`
	Height            uint64      `json:"Height"`
	Epoch             uint64      `json:"Epoch"`
	Round             int         `json:"Round"`
	Timestamp         int64       `json:"Timestamp"`
	PreviousBlockHash common.Hash `json:"PreviousBlockHash"`
	InstructionHash   common.Hash `json:"InstructionHash"` // hash of all parameters == hash of instruction
	ShardStateHash    common.Hash `json:"ShardStateHash"`  // each shard will have a list of blockHash, shardRoot is hash of all list
	// Merkle root of all instructions (using Keccak256 hash func) to relay to Ethreum
	// This obsoletes InstructionHash but for simplicity, we keep it for now
	InstructionMerkleRoot           common.Hash `json:"InstructionMerkleRoot"`
	BeaconCommitteeAndValidatorRoot common.Hash `json:"BeaconCommitteeAndValidatorRoot"` //Build from two list: BeaconCommittee + BeaconPendingValidator
	BeaconCandidateRoot             common.Hash `json:"BeaconCandidateRoot"`             // CandidateBeaconWaitingForCurrentRandom + CandidateBeaconWaitingForNextRandom
	ShardCandidateRoot              common.Hash `json:"ShardCandidateRoot"`              // CandidateShardWaitingForCurrentRandom + CandidateShardWaitingForNextRandom
	ShardCommitteeAndValidatorRoot  common.Hash `json:"ShardCommitteeAndValidatorRoot"`
	AutoStakingRoot                 common.Hash `json:"AutoStakingRoot"`
	ConsensusType                   string      `json:"ConsensusType"`
	Producer                        string      `json:"Producer"`
	ProducerPubKeyStr               string      `json:"ProducerPubKeyStr"`
}

func (beaconHeader *BeaconHeader) toString() string {
	res := ""
	// res += beaconHeader.ProducerAddress.String()
	res += fmt.Sprintf("%v", beaconHeader.Version)
	res += fmt.Sprintf("%v", beaconHeader.Height)
	res += fmt.Sprintf("%v", beaconHeader.Epoch)
	res += fmt.Sprintf("%v", beaconHeader.Round)
	res += fmt.Sprintf("%v", beaconHeader.Timestamp)
	res += beaconHeader.PreviousBlockHash.String()
	res += beaconHeader.BeaconCommitteeAndValidatorRoot.String()
	res += beaconHeader.BeaconCandidateRoot.String()
	res += beaconHeader.ShardCandidateRoot.String()
	res += beaconHeader.ShardCommitteeAndValidatorRoot.String()
	res += beaconHeader.AutoStakingRoot.String()
	res += beaconHeader.ShardStateHash.String()
	res += beaconHeader.InstructionHash.String()
	return res
}

func (beaconHeader BeaconHeader) GetMetaHash() common.Hash {
	return common.Keccak256([]byte(beaconHeader.toString()))
}

func (beaconHeader BeaconHeader) GetHash() *common.Hash {
	// Block header of beacon uses Keccak256 as a hash func to check on Ethereum when relaying blocks
	blkMetaHash := beaconHeader.GetMetaHash()
	blkInstHash := beaconHeader.InstructionMerkleRoot
	combined := append(blkMetaHash[:], blkInstHash[:]...)
	result := common.Keccak256(combined)
	return &result
}

func (beaconHeader BeaconHeader) GetTimestamp() int64 {
	return beaconHeader.Timestamp
}
func (beaconHeader BeaconHeader) GetVersion() int {
	return beaconHeader.Version
}
func (beaconHeader BeaconHeader) GetHeight() uint64 {
	return beaconHeader.Height
}
func (beaconHeader BeaconHeader) GetEpoch() uint64 {
	return beaconHeader.Epoch
}
func (beaconHeader BeaconHeader) GetConsensusType() string {
	return beaconHeader.ConsensusType
}
func (beaconHeader BeaconHeader) GetProducer() string {
	return beaconHeader.Producer
}
func (beaconHeader BeaconHeader) GetPreviousBlockHash() common.Hash {
	return beaconHeader.PreviousBlockHash
}
func (beaconHeader BeaconHeader) GetRound() int {
	return beaconHeader.Round
}
func (beaconHeader BeaconHeader) GetInstructionHash() common.Hash {
	return beaconHeader.InstructionHash
}
func (beaconHeader BeaconHeader) GetShardStateHash() common.Hash {
	return beaconHeader.ShardStateHash
}
func (beaconHeader BeaconHeader) GetInstructionMerkleRoot() common.Hash {
	return beaconHeader.InstructionMerkleRoot
}
func (beaconHeader BeaconHeader) GetBeaconCommitteeAndValidatorRoot() common.Hash {
	return beaconHeader.BeaconCommitteeAndValidatorRoot
}
func (beaconHeader BeaconHeader) GetBeaconCandidateRoot() common.Hash {
	return beaconHeader.BeaconCandidateRoot
}
func (beaconHeader BeaconHeader) GetShardCandidateRoot() common.Hash {
	return beaconHeader.ShardCandidateRoot
}
func (beaconHeader BeaconHeader) GetShardCommitteeAndValidatorRoot() common.Hash {
	return beaconHeader.ShardCommitteeAndValidatorRoot
}
func (beaconHeader BeaconHeader) GetAutoStakingRoot() common.Hash {
	return beaconHeader.AutoStakingRoot
}

func (beaconHeader BeaconHeader) GetBlockType() string {
	return "beacon"
}

func (beaconHeader BeaconHeader) GetTimeslot() uint64 {
	return 0
}
