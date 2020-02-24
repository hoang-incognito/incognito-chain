package app

import (
	"time"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/btc"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/database"
	"github.com/incognitochain/incognito-chain/incognitokey"
	"github.com/incognitochain/incognito-chain/mempool"
	"github.com/incognitochain/incognito-chain/metadata"
	"github.com/incognitochain/incognito-chain/privacy"
)

type FakeRandomClient struct{}

func (FakeRandomClient) GetNonceByTimestamp(startTime time.Time, maxTime time.Duration, timestamp int64) (int, int64, int64, error) {
	panic("implement me")
}

func (FakeRandomClient) VerifyNonceWithTimestamp(startTime time.Time, maxTime time.Duration, timestamp int64, nonce int64) (bool, error) {
	panic("implement me")
}

func (FakeRandomClient) GetCurrentChainTimeStamp() (int64, error) {
	return time.Now().Unix(), nil
}

func (FakeRandomClient) GetTimeStampAndNonceByBlockHeight(blockHeight int) (int64, int64, error) {
	panic("implement me")
}

type FakeBC struct {
	ShardToBeaconPool blockchain.ShardToBeaconPool
	CrossShardPool    map[byte]blockchain.CrossShardPool
	ShardPool         map[byte]blockchain.ShardPool
}

func (FakeBC) GetRandomClient() btc.RandomClient {
	return FakeRandomClient{}
}

func (FakeBC) GetBeaconHeightBreakPointBurnAddr() uint64 {
	panic("implement me")
}

func (FakeBC) GetBurningAddress(blockHeight uint64) string {
	panic("implement me")
}

func (FakeBC) InitTxSalaryByCoinID(payToAddress *privacy.PaymentAddress, amount uint64, payByPrivateKey *privacy.PrivateKey, meta metadata.Metadata, coinID common.Hash, shardID byte) (metadata.Transaction, error) {
	return nil, nil
}
func (FakeBC) GetCommitteeReward(committeeAddress []byte, tokenID common.Hash) (uint64, error) {
	return 0, nil
}

func (FakeBC) GetStakingAmountShard() uint64 {
	panic("implement me")
}

func (FakeBC) GetTxChainHeight(tx metadata.Transaction) (uint64, error) {
	panic("implement me")
}

func (FakeBC) GetChainHeight(byte) uint64 {
	panic("implement me")
}

func (FakeBC) GetBeaconHeight() uint64 {
	panic("implement me")
}

func (FakeBC) GetCurrentBeaconBlockHeight(byte) uint64 {
	panic("implement me")
}

func (FakeBC) GetAllCommitteeValidatorCandidate() (map[byte][]incognitokey.CommitteePublicKey, map[byte][]incognitokey.CommitteePublicKey, []incognitokey.CommitteePublicKey, []incognitokey.CommitteePublicKey, []incognitokey.CommitteePublicKey, []incognitokey.CommitteePublicKey, []incognitokey.CommitteePublicKey, []incognitokey.CommitteePublicKey, error) {
	panic("implement me")
}

func (FakeBC) GetAllCommitteeValidatorCandidateFlattenListFromDatabase() ([]string, error) {
	panic("implement me")
}

func (FakeBC) GetStakingTx(byte) map[string]string {
	panic("implement me")
}

func (FakeBC) GetAutoStakingList() map[string]bool {
	panic("implement me")
}

func (FakeBC) GetDatabase() database.DatabaseInterface {
	panic("implement me")
}

func (FakeBC) GetTxValue(txid string) (uint64, error) {
	panic("implement me")
}

func (FakeBC) GetShardIDFromTx(txid string) (byte, error) {
	panic("implement me")
}

func (FakeBC) GetCentralizedWebsitePaymentAddress() string {
	panic("implement me")
}

func (FakeBC) GetAllCoinID() ([]common.Hash, error) {
	panic("implement me")
}

func (FakeBC) GetValidBeaconBlockFromPool() []blockinterface.BeaconBlockInterface {
	panic("implement me")
}

func (FakeBC) GetShardPendingCommittee(shardID byte) []incognitokey.CommitteePublicKey {
	return []incognitokey.CommitteePublicKey{}
}

//look at beacon chain, get the final view and get height
func (FakeBC) GetCurrentBeaconHeight() (uint64, error) {
	return 1, nil
}

func (FakeBC) GetEpoch() (uint64, error) {
	panic("implement me")
}

func (FakeBC) GetChainParams() blockchain.Params {
	return blockchain.Params{
		Epoch:                  10,
		MinShardBlockInterval:  blockchain.TestNetMinShardBlkInterval,
		MaxShardBlockCreation:  blockchain.TestNetMaxShardBlkCreation,
		MinBeaconBlockInterval: blockchain.TestNetMinBeaconBlkInterval,
		MaxBeaconBlockCreation: blockchain.TestNetMaxBeaconBlkCreation,
	}
}

func (FakeBC) GetCrossShardPool(shardID byte) blockchain.CrossShardPool {
	return mempool.GetCrossShardPool(shardID)
}

func (FakeBC) GetLatestCrossShard(from byte, to byte) uint64 {
	return 1
}

func (FakeBC) GetNextCrossShard(from byte, to byte, startHeight uint64) uint64 {
	return 1
}

func (FakeBC) GetAllValidCrossShardBlockFromPool(toShard byte) map[byte][]blockinterface.CrossShardBlockInterface {
	return nil
}

func (FakeBC) ValidateCrossShardBlock(block blockinterface.CrossShardBlockInterface) error {
	return nil
}

func (FakeBC) GetPendingTransaction(shardID byte) (txsToAdd []metadata.Transaction, txToRemove []metadata.Transaction, totalFee uint64) {
	return []metadata.Transaction{}, []metadata.Transaction{}, uint64(0)
}

func (FakeBC) GetTransactionByHash(hash common.Hash) (byte, common.Hash, int, metadata.Transaction, error) {
	return 0, common.Hash{}, 0, nil, nil
}

func (FakeBC) FetchAutoStakingByHeight(uint64) (map[string]bool, error) {
	return map[string]bool{}, nil
}

type RawDB interface {
}

type StateDB interface {
}

type StateObject interface {
}

type DatabaseAccessWarper interface {
}
