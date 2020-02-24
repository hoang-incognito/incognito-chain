package app

import (
	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/btc"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/incognitokey"
	"github.com/incognitochain/incognito-chain/metadata"
	"github.com/incognitochain/incognito-chain/privacy"
)

type ShardApp interface {
	//create block
	preCreateBlock() error
	buildTxFromCrossShard() error             // build tx from crossshard
	buildTxFromMemPool() error                // build tx from mempool
	buildResponseTxFromTxWithMetadata() error // build tx from metadata tx
	processBeaconInstruction() error          // execute beacon instruction & build tx if any
	generateInstruction() error               //create block instruction

	buildHeader() error

	//crete view from block
	updateNewViewFromBlock(block blockinterface.ShardBlockInterface) error

	//validate block
	preValidate() error

	//store block
	storeDatabase(state *StoreShardDatabaseState) error
}

type BeaconApp interface {
	//create block
	preCreateBlock() error
	buildInstructionByEpoch() error
	buildInstructionFromShardAction() error

	buildHeader() error

	//crete view from block
	updateNewViewFromBlock(block blockinterface.BeaconBlockInterface) error

	//validate block
	preValidate() error

	////store block
	storeDatabase() error
}

type AppData struct {
	Logger      common.Logger
	CreateBlock struct {
		crossShardTx             map[byte][]blockchain.CrossTransaction
		txToRemove               []metadata.Transaction
		txsToAdd                 []metadata.Transaction
		txsFromMetadataTx        []metadata.Transaction
		txsFromBeaconInstruction []metadata.Transaction
		errInstruction           [][]string
		stakingTx                map[string]string
		newShardPendingValidator []string
		instruction              [][]string
	}
}

type DB interface {
	GetGenesisBlock() blockinterface.BlockInterface
	GetAllTokenIDForReward(epoch uint64) ([]common.Hash, error)
	GetRewardOfShardByEpoch(epoch uint64, shardID byte, tokenID common.Hash) (uint64, error)
	//GetBeaconBlockHashByIndex(uint64) (common.Hash, error)
	//FetchBeaconBlock(common.Hash) ([]byte, error)
}

type BlockChain interface {
	//GetDB() DB
	GetCurrentBeaconHeight() (uint64, error) //get final confirm beacon block height
	GetEpoch() (uint64, error)               //get final confirm beacon block height
	GetChainParams() blockchain.Params

	ValidateCrossShardBlock(block blockinterface.CrossShardBlockInterface) error

	GetCrossShardPool(shardID byte) blockchain.CrossShardPool
	GetLatestCrossShard(from byte, to byte) uint64
	GetNextCrossShard(from byte, to byte, startHeight uint64) uint64

	GetAllValidCrossShardBlockFromPool(toShard byte) map[byte][]blockinterface.CrossShardBlockInterface
	GetValidBeaconBlockFromPool() []blockinterface.BeaconBlockInterface
	GetPendingTransaction(shardID byte) (txsToAdd []metadata.Transaction, txToRemove []metadata.Transaction, totalFee uint64)

	GetShardPendingCommittee(shardID byte) []incognitokey.CommitteePublicKey
	//GetShardCommittee(shardID byte) []incognitokey.CommitteePublicKey
	FetchAutoStakingByHeight(uint64) (map[string]bool, error)

	GetCommitteeReward(committeeAddress []byte, tokenID common.Hash) (uint64, error)
	InitTxSalaryByCoinID(payToAddress *privacy.PaymentAddress, amount uint64, payByPrivateKey *privacy.PrivateKey, meta metadata.Metadata, coinID common.Hash, shardID byte) (metadata.Transaction, error)

	metadata.BlockchainRetriever

	GetRandomClient() btc.RandomClient
}
