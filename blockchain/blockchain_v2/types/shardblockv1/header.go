package shardblockv1

import (
	"fmt"
	"sort"

	"github.com/incognitochain/incognito-chain/common"
)

/*
	-TxRoot and MerkleRootShard: make from transaction
	-Validator Root is root hash of current committee in beststate
	-PendingValidator Root is root hash of pending validator in beststate
*/
type ShardHeader struct {
	Producer              string                 `json:"Producer"`
	ProducerPubKeyStr     string                 `json:"ProducerPubKeyStr"`
	ShardID               byte                   `json:"ShardID"`               // shard ID which block belong to
	Version               int                    `json:"Version"`               // version of block structure
	PreviousBlockHash     common.Hash            `json:"PreviousBlockHash"`     // previous block hash or Parent block hash
	Height                uint64                 `json:"Height"`                // block height
	Round                 int                    `json:"Round"`                 // bpft consensus round
	Epoch                 uint64                 `json:"Epoch"`                 // epoch of block (according to current beacon height)
	CrossShardBitMap      []byte                 `json:"CrossShardBitMap"`      // crossShards bitmap for beacon
	BeaconHeight          uint64                 `json:"BeaconHeight"`          // beacon check point height
	BeaconHash            common.Hash            `json:"BeaconHash"`            // beacon check point hash
	TotalTxsFee           map[common.Hash]uint64 `json:"TotalTxsFee"`           // fee of all txs in block
	ConsensusType         string                 `json:"ConsensusType"`         // consensus type, by which this block is produced
	Timestamp             int64                  `json:"Timestamp"`             // timestamp of block
	TxRoot                common.Hash            `json:"TxRoot"`                // Transaction root created from transaction in shard
	ShardTxRoot           common.Hash            `json:"ShardTxRoot"`           // output root created for other shard
	CrossTransactionRoot  common.Hash            `json:"CrossTransactionRoot"`  // transaction root created from transaction of micro shard to shard block (from other shard)
	InstructionsRoot      common.Hash            `json:"InstructionsRoot"`      // actions root created from Instructions and Metadata of transaction
	CommitteeRoot         common.Hash            `json:"CommitteeRoot"`         // hash from public key list of all committees designated to create this block
	PendingValidatorRoot  common.Hash            `json:"PendingValidatorRoot"`  // hash from public key list of all pending validators designated to this ShardID
	StakingTxRoot         common.Hash            `json:"StakingTxRoot"`         // hash from staking transaction map in shard best state
	InstructionMerkleRoot common.Hash            `json:"InstructionMerkleRoot"` // Merkle root of all instructions (using Keccak256 hash func) to relay to Ethreum
}

func (shardHeader ShardHeader) String() string {
	res := common.EmptyString
	// res += shardHeader.ProducerAddress.String()
	res += string(shardHeader.ShardID)
	res += fmt.Sprintf("%v", shardHeader.Version)
	res += shardHeader.PreviousBlockHash.String()
	res += fmt.Sprintf("%v", shardHeader.Height)
	res += fmt.Sprintf("%v", shardHeader.Round)
	res += fmt.Sprintf("%v", shardHeader.Epoch)
	res += fmt.Sprintf("%v", shardHeader.Timestamp)
	res += shardHeader.TxRoot.String()
	res += shardHeader.ShardTxRoot.String()
	res += shardHeader.CrossTransactionRoot.String()
	res += shardHeader.InstructionsRoot.String()
	res += shardHeader.CommitteeRoot.String()
	res += shardHeader.PendingValidatorRoot.String()
	res += shardHeader.BeaconHash.String()
	res += shardHeader.StakingTxRoot.String()
	res += fmt.Sprintf("%v", shardHeader.BeaconHeight)
	tokenIDs := make([]common.Hash, 0)
	for tokenID, _ := range shardHeader.TotalTxsFee {
		tokenIDs = append(tokenIDs, tokenID)
	}
	sort.Slice(tokenIDs, func(i int, j int) bool {
		res, _ := tokenIDs[i].Cmp(&tokenIDs[j])
		return res == -1
	})

	for _, tokenID := range tokenIDs {
		res += fmt.Sprintf("%v~%v", tokenID.String(), shardHeader.TotalTxsFee[tokenID])
	}
	for _, value := range shardHeader.CrossShardBitMap {
		res += string(value)
	}
	return res
}
func (shardHeader ShardHeader) GetMetaHash() common.Hash {
	return common.Keccak256([]byte(shardHeader.String()))
}

func (shardHeader ShardHeader) GetHash() *common.Hash {
	// Block header of bridge uses Keccak256 as a hash func to check on Ethereum when relaying blocks
	blkMetaHash := shardHeader.GetMetaHash()
	blkInstHash := shardHeader.InstructionMerkleRoot
	combined := append(blkMetaHash[:], blkInstHash[:]...)
	result := common.Keccak256(combined)
	return &result
}

func (shardHeader ShardHeader) GetProducer() string { return shardHeader.Producer }
func (shardHeader ShardHeader) GetShardID() byte    { return shardHeader.ShardID }
func (shardHeader ShardHeader) GetVersion() int     { return shardHeader.Version }
func (shardHeader ShardHeader) GetPreviousBlockHash() common.Hash {
	return shardHeader.PreviousBlockHash
}
func (shardHeader ShardHeader) GetHeight() uint64           { return shardHeader.Height }
func (shardHeader ShardHeader) GetRound() int               { return shardHeader.Round }
func (shardHeader ShardHeader) GetEpoch() uint64            { return shardHeader.Epoch }
func (shardHeader ShardHeader) GetCrossShardBitMap() []byte { return shardHeader.CrossShardBitMap }
func (shardHeader ShardHeader) GetBeaconHeight() uint64     { return shardHeader.BeaconHeight }
func (shardHeader ShardHeader) GetBeaconHash() common.Hash  { return shardHeader.BeaconHash }
func (shardHeader ShardHeader) GetTotalTxsFee() map[common.Hash]uint64 {
	return shardHeader.TotalTxsFee
}
func (shardHeader ShardHeader) GetConsensusType() string    { return shardHeader.ConsensusType }
func (shardHeader ShardHeader) GetTimestamp() int64         { return shardHeader.Timestamp }
func (shardHeader ShardHeader) GetTxRoot() common.Hash      { return shardHeader.TxRoot }
func (shardHeader ShardHeader) GetShardTxRoot() common.Hash { return shardHeader.ShardTxRoot }
func (shardHeader ShardHeader) GetCrossTransactionRoot() common.Hash {
	return shardHeader.CrossTransactionRoot
}
func (shardHeader ShardHeader) GetInstructionsRoot() common.Hash { return shardHeader.InstructionsRoot }
func (shardHeader ShardHeader) GetCommitteeRoot() common.Hash    { return shardHeader.CommitteeRoot }
func (shardHeader ShardHeader) GetPendingValidatorRoot() common.Hash {
	return shardHeader.PendingValidatorRoot
}
func (shardHeader ShardHeader) GetStakingTxRoot() common.Hash { return shardHeader.StakingTxRoot }
func (shardHeader ShardHeader) GetInstructionMerkleRoot() common.Hash {
	return shardHeader.InstructionMerkleRoot
}
