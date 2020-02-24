package app

import (
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
)

type ShardBridgeApp struct {
	Logger        common.Logger
	CreateState   *CreateShardBlockState
	ValidateState *ValidateShardBlockState
	StoreState    *StoreShardDatabaseState
	storeSuccess  bool
}

func (s *ShardBridgeApp) preCreateBlock() error {
	return nil
}
func (s *ShardBridgeApp) buildTxFromCrossShard() error {
	return nil
}
func (s *ShardBridgeApp) buildTxFromMemPool() error {
	return nil
}
func (s *ShardBridgeApp) buildResponseTxFromTxWithMetadata() error {
	return nil
}
func (s *ShardBridgeApp) processBeaconInstruction() error {
	return nil
}
func (s *ShardBridgeApp) generateInstruction() error {
	return nil
}

func (s *ShardBridgeApp) buildHeader() error {
	return nil
}

//crete view from block
func (s *ShardBridgeApp) updateNewViewFromBlock(block blockinterface.ShardBlockInterface) error {
	return nil
}

//validate block
func (s *ShardBridgeApp) preValidate() error {
	return nil
}

//store block
func (s *ShardBridgeApp) storeDatabase(state *StoreShardDatabaseState) error {
	return nil
}
