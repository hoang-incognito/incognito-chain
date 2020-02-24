package shardblockv1

import (
	"fmt"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
)

type ShardToBeaconBlock struct {
	ValidationData string `json:"ValidationData"`

	Instructions [][]string
	Header       ShardHeader
}

func (block ShardToBeaconBlock) GetConsensusType() string {
	return block.Header.ConsensusType
}

func (block ShardToBeaconBlock) GetValidationField() string {
	return block.ValidationData
}

func (block ShardToBeaconBlock) GetHeight() uint64 {
	return block.Header.Height
}

func (block ShardToBeaconBlock) GetRound() int {
	return block.Header.Round
}

func (block ShardToBeaconBlock) GetRoundKey() string {
	return fmt.Sprint(block.Header.Height, "_", block.Header.Round)
}
func (block ShardToBeaconBlock) GetInstructions() [][]string {
	return block.Instructions
}

func (block ShardToBeaconBlock) GetProducer() string {
	return block.Header.Producer
}

func (shardToBeaconBlock ShardToBeaconBlock) GetEpoch() uint64 {
	return shardToBeaconBlock.Header.Epoch
}

func (shardToBeaconBlock *ShardToBeaconBlock) Hash() *common.Hash {
	return shardToBeaconBlock.Header.GetHash()
}

func (block ShardToBeaconBlock) GetShardHeader() blockinterface.ShardHeaderInterface {
	return block.Header
}
