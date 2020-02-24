package blockinterface

import "github.com/incognitochain/incognito-chain/common"

type BlockInterface interface {
	GetHeader() BlockHeaderInterface
	GetBody() BlockBodyInterface
	GetValidationField() string
	GetBlockType() string
}

type BlockHeaderInterface interface {
	GetVersion() int
	GetConsensusType() string
	GetMetaHash() common.Hash
	GetHash() *common.Hash
	GetPreviousBlockHash() common.Hash
	GetHeight() uint64
	GetEpoch() uint64
	GetTimestamp() int64
	GetProducer() string
}

type BlockHeaderV1Interface interface {
	BlockHeaderInterface
	GetRound() int
	GetRoundKey() string
}

type BlockHeaderV2Interface interface {
	BlockHeaderInterface
	GetTimeslot() uint64
}

type BlockBodyInterface interface {
	GetInstructions() [][]string
}
