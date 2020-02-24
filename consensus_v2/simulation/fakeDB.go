package main

import (
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
)

type FakeDB struct {
	genesisBlock blockinterface.BlockInterface
}

func (s *FakeDB) GetRewardOfShardByEpoch(epoch uint64, shardID byte, tokenID common.Hash) (uint64, error) {
	//TODO
	return 0, nil
}

func (s *FakeDB) GetAllTokenIDForReward(epoch uint64) ([]common.Hash, error) {
	//TODO
	return []common.Hash{}, nil
}

func (s *FakeDB) GetGenesisBlock() blockinterface.BlockInterface {
	return s.genesisBlock
}
