package app

import (
	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/database"
)

type BeaconBridgeApp struct {
	Logger        common.Logger
	CreateState   *CreateBeaconBlockState
	ValidateState *ValidateBeaconBlockState
	StoreState    *StoreBeaconDatabaseState
	storeSuccess  bool
}

func (s *BeaconBridgeApp) preCreateBlock() error {
	return nil
}

func (s *BeaconBridgeApp) buildInstructionByEpoch() error {
	return nil
}

func (s *BeaconBridgeApp) buildInstructionFromShardAction() error {
	if len(s.CreateState.s2bBlks) == 0 {
		return nil
	}
	db := s.CreateState.bc.GetDatabase()
	newBeaconHeight := s.CreateState.curView.GetHeight() + 1

	for shardID, shardBlocks := range s.CreateState.s2bBlks {
		for _, block := range shardBlocks {
			bridgeInstructionForBlock, err := buildBridgeInstructions(
				shardID,
				block.GetInstructions(),
				newBeaconHeight,
				db,
				s.Logger,
			)
			if err != nil {
				s.Logger.Errorf("Build bridge instructions failed: %s", err.Error())
			}
			// Pick instruction with shard committee's pubkeys to save to beacon block
			confirmInsts := pickBridgeSwapConfirmInst(block)
			if len(confirmInsts) > 0 {
				bridgeInstructionForBlock = append(bridgeInstructionForBlock, confirmInsts...)
				s.Logger.Infof("Beacon block %d found bridge swap confirm inst in shard block %d: %s", newBeaconHeight, block.GetShardHeader().GetHeight(), confirmInsts)
			}

			s.CreateState.bridgeInstructions = append(s.CreateState.bridgeInstructions, bridgeInstructionForBlock...)
		}
	}
	return nil
}

func (s *BeaconBridgeApp) buildHeader() error {
	return nil
}

func (s *BeaconBridgeApp) updateNewViewFromBlock(block blockinterface.BeaconBlockInterface) error {
	return nil
}

func (s *BeaconBridgeApp) preValidate() error {
	return nil
}

//==============================Save Database Logic===========================
func (s *BeaconBridgeApp) storeDatabase() error {
	//TODO: store db?
	batchPutData := []database.BatchData{}
	err := storeBridgeInstructions(s.StoreState.block, &batchPutData, s.StoreState.bc, s.Logger)
	if err != nil {
		return blockchain.NewBlockChainError(blockchain.ProcessBridgeInstructionError, err)
	}
	return nil
}
