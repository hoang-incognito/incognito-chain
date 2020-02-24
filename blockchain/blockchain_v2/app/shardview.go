package app

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/shardblockv2"
	"github.com/incognitochain/incognito-chain/common"
	consensus "github.com/incognitochain/incognito-chain/consensus_v2"
	"github.com/incognitochain/incognito-chain/consensus_v2/blsbftv2"
	"github.com/incognitochain/incognito-chain/incognitokey"
)

type ShardView struct {
	//field that copy manualy
	BC     BlockChain
	DB     DB `json:"-"`
	Lock   *sync.RWMutex
	Logger common.Logger

	//field that copy automatically and need to update
	ShardID               byte
	Block                 blockinterface.ShardBlockInterface
	ShardCommittee        []incognitokey.CommitteePublicKey
	ShardPendingValidator []incognitokey.CommitteePublicKey

	StakingTx map[string]string

	//================================ StateDB Method
	// block height => root hash
	// consensusStateDB   *statedb.StateDB
	// transactionStateDB *statedb.StateDB
	// featureStateDB     *statedb.StateDB
	// rewardStateDB      *statedb.StateDB
	// slashStateDB       *statedb.StateDB

	consensusStateDB   StateDB
	transactionStateDB StateDB
	featureStateDB     StateDB
	rewardStateDB      StateDB
	slashStateDB       StateDB
}

func (s *ShardView) CreateBlockFromOldBlockData(block blockinterface.BlockInterface) blockinterface.BlockInterface {
	block1 := block.(*shardblockv2.ShardBlock)
	block1.ConsensusHeader.TimeSlot = common.GetTimeSlot(s.GetGenesisTime(), time.Now().Unix(), blsbftv2.TIMESLOT)
	return block1
}

func (s *ShardView) GetBlock() blockinterface.BlockInterface {
	return s.Block
}

func (s *ShardView) CreateNewViewFromBlock(block blockinterface.BlockInterface) (consensus.ChainViewInterface, error) {
	panic("implement me")
}

func (s *ShardView) UnmarshalBlock(b []byte) (blockinterface.BlockInterface, error) {
	block, err := UnmarshalShardBlock(b)
	if err != nil {
		return nil, err
	}
	return block.(blockinterface.BlockInterface), nil
}

func (s *ShardView) GetGenesisTime() int64 {
	return s.DB.GetGenesisBlock().GetHeader().GetTimestamp()
}

func (s *ShardView) GetConsensusConfig() string {
	panic("implement me")
}

func (s *ShardView) GetConsensusType() string {
	return "bls"
}

func (s *ShardView) GetBlkMinInterval() time.Duration {
	panic("implement me")
}

func (s *ShardView) GetBlkMaxCreateTime() time.Duration {
	panic("implement me")
}

func (s *ShardView) GetPubkeyRole(pubkey string, round int) (string, byte) {
	panic("implement me")
}

func (s *ShardView) GetCommittee() []incognitokey.CommitteePublicKey {
	return s.ShardCommittee
}

func (s *ShardView) GetCommitteeHash() common.Hash {
	panic("implement me")
}

func (s *ShardView) GetCommitteeIndex(string) int {
	panic("implement me")
}

func (s *ShardView) GetHeight() uint64 {
	return s.Block.GetHeader().GetHeight()
}

// func (s *ShardView) GetRound() int {
// 	return s.Block.Header.Round
// }

func (s *ShardView) GetTimeStamp() int64 {
	return s.Block.GetHeader().GetTimestamp()
}

func (s *ShardView) GetTimeslot() uint64 {
	if s.Block.GetHeader().GetVersion() == 1 {
		return 1
	}
	return s.Block.GetHeader().(blockinterface.BlockHeaderV2Interface).GetTimeslot()
}

func (s *ShardView) GetEpoch() uint64 {
	return s.Block.GetHeader().GetEpoch()
}

func (s *ShardView) Hash() common.Hash {
	return *s.Block.GetHeader().GetHash()
}

func (s *ShardView) GetPreviousViewHash() common.Hash {
	prevHash := s.Block.GetHeader().GetPreviousBlockHash()
	return prevHash
}

func (s *ShardView) GetActiveShardNumber() int {
	panic("implement me")
}

func (s *ShardView) GetNextProposer(timeSlot uint64) string {
	committee := s.GetCommittee()
	idx := int(timeSlot) % len(committee)
	return committee[idx].GetMiningKeyBase58(common.BlsConsensus)
}

func (s *ShardView) CloneNewView() consensus.ChainViewInterface {
	b, _ := s.MarshalJSON()
	var newView *ShardView
	err := json.Unmarshal(b, &newView)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	newView.DB = s.DB
	newView.Logger = s.Logger
	newView.Lock = &sync.RWMutex{}
	newView.BC = s.BC
	return newView
}

func (s *ShardView) MarshalJSON() ([]byte, error) {
	type Alias ShardView
	b, err := json.Marshal(&struct {
		*Alias
		DB     interface{}
		Lock   interface{}
		Logger interface{}
		BC     interface{}
		Block  interface{}
	}{
		(*Alias)(s),
		nil,
		nil,
		nil,
		nil,
	})
	if err != nil {
		Logger.log.Error(err)
	}
	return b, err
}

func (s *ShardView) GetRootTimeSlot() uint64 {
	if s.DB.GetGenesisBlock().GetHeader().GetVersion() == 1 {
		return 1
	}
	return s.DB.GetGenesisBlock().GetHeader().(blockinterface.BlockHeaderV2Interface).GetTimeslot()
}

func (s *ShardView) InitStateRootHash(bc *BlockChain) error {
	panic("implement me")
}
