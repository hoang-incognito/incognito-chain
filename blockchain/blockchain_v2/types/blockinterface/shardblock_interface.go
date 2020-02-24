package blockinterface

import (
	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/metadata"
	"github.com/incognitochain/incognito-chain/privacy"
)

type ShardBlockInterface interface {
	BlockInterface

	GetShardHeader() ShardHeaderInterface
	GetShardBody() ShardBodyInterface
}

type ShardHeaderInterface interface {
	BlockHeaderInterface

	GetShardID() byte
	GetCrossShardBitMap() []byte
	GetBeaconHeight() uint64
	GetBeaconHash() common.Hash
	GetTotalTxsFee() map[common.Hash]uint64
	GetTxRoot() common.Hash
	GetShardTxRoot() common.Hash
	GetCrossTransactionRoot() common.Hash
	GetInstructionsRoot() common.Hash
	GetCommitteeRoot() common.Hash
	GetPendingValidatorRoot() common.Hash
	GetStakingTxRoot() common.Hash
	GetInstructionMerkleRoot() common.Hash
}

type ShardHeaderV2Interface interface {
	GetTimeSlot() uint64
}

type ShardBodyInterface interface {
	BlockBodyInterface

	GetTransactions() []metadata.Transaction
	GetCrossTransactions() map[byte][]blockchain.CrossTransaction
}

type ShardToBeaconBlockInterface interface {
	GetValidationField() string
	GetShardHeader() ShardHeaderInterface
	GetInstructions() [][]string
}

type CrossShardBlockInterface interface {
	GetValidationField() string
	GetShardHeader() ShardHeaderInterface

	GetToShardID() byte
	GetMerklePathShard() []common.Hash
	// Cross Shard data for PRV
	GetCrossOutputCoin() []privacy.OutputCoin
	// Cross Shard For Custom token privacy
	GetCrossTxTokenPrivacyData() []blockchain.ContentCrossShardTokenPrivacyData
}
