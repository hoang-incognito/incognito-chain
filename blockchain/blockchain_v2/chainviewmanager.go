package blockchain_v2

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/common"
	consensus "github.com/incognitochain/incognito-chain/consensus_v2"
	"github.com/incognitochain/incognito-chain/incognitokey"
)

/*
	How to load Chain Manager when bootup
	Logic to store chainview
	How to manage Chain View
*/

type ChainViewManager struct {
	manager *ViewGraph
	name    string
	chainID int
	lock    *sync.RWMutex
}

func (s *ChainViewManager) CommitView(view consensus.ChainViewInterface) error {
	return view.StoreDatabase(context.Background())
}

func (s *ChainViewManager) ConnectBlockAndAddView(block blockinterface.BlockInterface) error {
	preBlkHash := block.GetHeader().GetPreviousBlockHash()
	view, err := s.GetViewByHash(preBlkHash)
	if err != nil {
		panic(err)
		return err
	}
	newView, err := view.ValidateBlockAndCreateNewView(context.Background(), block, false)
	if err != nil {
		panic(err)
		return err
	}

	s.manager.AddView(newView)
	s.manager.Print()

	//persist view, block, related block data to disk
	if err := s.CommitView(newView); err != nil {
		panic(err)
	}
	return nil
}

//to get byte array of this chain data, for loading later
func (s *ChainViewManager) GetChainData() (res []byte) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

//create new chain with root manager
func InitNewChainViewManager(name string, chainID int, rootView consensus.ChainViewInterface) *ChainViewManager {
	cm := &ChainViewManager{
		name: name,
		lock: new(sync.RWMutex),
	}
	cm.manager = NewViewGraph(name, rootView, cm.lock)
	return cm
}

func LoadChain(data []byte) (res *ChainViewManager) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(res)
	if err != nil {
		panic(err)
	}
	return res
}

func (s *ChainViewManager) AddView(view consensus.ChainViewInterface) error {
	s.manager.AddView(view)
	return nil
}

func (ChainViewManager) GetChainName() string {
	panic("implement me")
}

func (ChainViewManager) IsReady() bool {
	panic("implement me")
}

func (cManager ChainViewManager) GetActiveShardNumber() int {
	return 0
}

func (ChainViewManager) GetPubkeyRole(pubkey string, round int) (string, byte) {
	panic("implement me")
}

func (ChainViewManager) GetConsensusType() string {
	panic("implement me")
}

func (ChainViewManager) GetTimeStamp() int64 {
	panic("implement me")
}

func (ChainViewManager) GetMinBlkInterval() time.Duration {
	panic("implement me")
}

func (ChainViewManager) GetMaxBlkCreateTime() time.Duration {
	panic("implement me")
}

func (ChainViewManager) GetHeight() uint64 {
	panic("implement me")
}

func (ChainViewManager) GetCommitteeSize() int {
	panic("implement me")
}

func (ChainViewManager) GetCommittee() []incognitokey.CommitteePublicKey {
	panic("implement me")
}

func (ChainViewManager) GetPubKeyCommitteeIndex(string) int {
	panic("implement me")
}

func (ChainViewManager) GetLastProposerIndex() int {
	panic("implement me")
}

func (s *ChainViewManager) UnmarshalBlock(blockString []byte) (blockinterface.BlockInterface, error) {
	return s.GetBestView().UnmarshalBlock(blockString)
}

func (ChainViewManager) ValidateBlockSignatures(block blockinterface.BlockInterface) error {
	panic("implement me")
}

func (s *ChainViewManager) ValidateBlockProposer(block blockinterface.BlockInterface) error {
	return nil
}

func (s *ChainViewManager) GetShardID() int {
	return s.chainID
}

func (s *ChainViewManager) GetBestView() consensus.ChainViewInterface {
	return s.manager.GetBestView()
}

func (s ChainViewManager) GetFinalView() consensus.ChainViewInterface {
	return s.manager.GetFinalView()
}

func (s *ChainViewManager) GetAllViews() map[string]consensus.ChainViewInterface {
	fmt.Println("All views:", len(s.manager.node))
	return nil
}

func (s *ChainViewManager) GetViewByRange(from string, to string) (res []consensus.ChainViewInterface) {
	for _, viewNode := range s.manager.node {
		if to == viewNode.view.Hash().String() {
			res = append(res, viewNode.view)
			for {
				if viewNode.prev != nil {
					res = append([]consensus.ChainViewInterface{viewNode.prev.view}, res...)
					if viewNode.prev.view.Hash().String() == from {
						break
					}
					viewNode = viewNode.prev
				} else {
					break
				}
			}
			break
		}
	}
	return res
}

func (s *ChainViewManager) GetViewByHash(h common.Hash) (consensus.ChainViewInterface, error) {
	viewnode := s.manager.GetNodeByHash(h)
	if viewnode == nil {
		return nil, errors.New("view not exist")
	} else {
		return viewnode.view, nil
	}
}

func (s *ChainViewManager) GetGenesisTime() int64 {
	return s.GetBestView().GetGenesisTime()
}

func (s *ChainViewManager) GetAllTipBlocksHash() []*common.Hash {
	var result []*common.Hash
	for _, node := range s.manager.node {
		result = append(result, node.view.GetBlock().GetHeader().GetHash())
	}
	return result
}
