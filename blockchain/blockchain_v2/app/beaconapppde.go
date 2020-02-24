package app

import (
	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/database"
)

type BeaconPDEApp struct {
	Logger        common.Logger
	CreateState   *CreateBeaconBlockState
	ValidateState *ValidateBeaconBlockState
	StoreState    *StoreBeaconDatabaseState
	storeSuccess  bool
}

func (s *BeaconPDEApp) preCreateBlock() error {
	return nil
}

func (s *BeaconPDEApp) buildInstructionByEpoch() error {
	return nil
}

func (s *BeaconPDEApp) buildInstructionFromShardAction() error {

	if len(s.CreateState.s2bBlks) == 0 {
		return nil
	}
	db := s.CreateState.bc.GetDatabase()
	newBeaconHeight := s.CreateState.curView.GetHeight() + 1
	statefulActionsByShardID := map[byte][][]string{}

	for shardID, shardBlocks := range s.CreateState.s2bBlks {
		for _, block := range shardBlocks {
			// Collect stateful actions
			statefulActions := collectStatefulActions(block.GetInstructions())
			// group stateful actions by shardID
			_, found := statefulActionsByShardID[shardID]
			if !found {
				statefulActionsByShardID[shardID] = statefulActions
			} else {
				statefulActionsByShardID[shardID] = append(statefulActionsByShardID[shardID], statefulActions...)
			}
		}
	}

	// build stateful instructions
	statefulInsts := buildStatefulInstructions(
		statefulActionsByShardID,
		newBeaconHeight,
		db,
		s.CreateState.bc,
	)

	s.CreateState.statefulInstructions = append(s.CreateState.statefulInstructions, statefulInsts...)
	return nil
}

func (s *BeaconPDEApp) buildHeader() error {
	return nil
}

func (s *BeaconPDEApp) updateNewViewFromBlock(block blockinterface.BeaconBlockInterface) error {
	return nil
}

func (s *BeaconPDEApp) preValidate() error {
	return nil
}

//==============================Save Database Logic===========================
func (s *BeaconPDEApp) storeDatabase() error {
	//TODO: store db?
	batchPutData := []database.BatchData{}
	// execute, store
	err := storePDEInstructions(s.StoreState.block, &batchPutData, s.StoreState.bc)
	if err != nil {
		return blockchain.NewBlockChainError(blockchain.ProcessPDEInstructionError, err)
	}
	return nil
}
