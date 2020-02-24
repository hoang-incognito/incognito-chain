package app

import (
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
)

type ShardPDEApp struct {
	Logger        common.Logger
	CreateState   *CreateShardBlockState
	ValidateState *ValidateShardBlockState
	StoreState    *StoreShardDatabaseState
	storeSuccess  bool
}

func (s *ShardPDEApp) preCreateBlock() error {
	return nil
}
func (s *ShardPDEApp) buildTxFromCrossShard() error {
	return nil
}
func (s *ShardPDEApp) buildTxFromMemPool() error {
	return nil
}
func (s *ShardPDEApp) buildResponseTxFromTxWithMetadata() error {
	return nil
}
func (s *ShardPDEApp) processBeaconInstruction() error {
	return nil
}
func (s *ShardPDEApp) generateInstruction() error {
	return nil
}

func (s *ShardPDEApp) buildHeader() error {
	return nil
}

//crete view from block
func (s *ShardPDEApp) updateNewViewFromBlock(block blockinterface.ShardBlockInterface) error {
	return nil
}

//validate block
func (s *ShardPDEApp) preValidate() error {
	return nil
}

//store block
func (s *ShardPDEApp) storeDatabase(state *StoreShardDatabaseState) error {
	s.storeSuccess = true
	return nil
}
