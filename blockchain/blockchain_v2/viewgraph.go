package blockchain_v2

import (
	"fmt"
	"os"
	"sync"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"

	"github.com/incognitochain/incognito-chain/common"
	consensus "github.com/incognitochain/incognito-chain/consensus_v2"
)

type ViewNode struct {
	view consensus.ChainViewInterface
	next map[common.Hash]*ViewNode
	prev *ViewNode
}
type ViewGraph struct {
	name        string
	root        *ViewNode
	node        map[common.Hash]*ViewNode
	leaf        map[common.Hash]*ViewNode
	edgeStr     string
	bestView    *ViewNode
	confirmView *ViewNode
	lock        *sync.RWMutex
}

func NewViewGraph(name string, rootView consensus.ChainViewInterface, lock *sync.RWMutex) *ViewGraph {
	s := &ViewGraph{name: name, lock: lock}
	s.leaf = make(map[common.Hash]*ViewNode)
	s.node = make(map[common.Hash]*ViewNode)
	s.root = &ViewNode{
		view: rootView,
		next: make(map[common.Hash]*ViewNode),
		prev: nil,
	}
	s.leaf[*rootView.GetBlock().GetHeader().GetHash()] = s.root
	s.node[*rootView.GetBlock().GetHeader().GetHash()] = s.root
	s.confirmView = s.root
	s.update()
	return s
}

func (s *ViewGraph) GetNodeByHash(h common.Hash) *ViewNode {
	return s.node[h]
}

func (s *ViewGraph) AddView(b consensus.ChainViewInterface) {
	newBlockHash := *b.GetBlock().GetHeader().GetHash()
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.node[newBlockHash] != nil {
		return
	}

	for h, v := range s.node {
		if h == b.GetBlock().GetHeader().GetPreviousBlockHash() {
			delete(s.leaf, h)
			s.leaf[newBlockHash] = &ViewNode{
				view: b,
				next: make(map[common.Hash]*ViewNode),
				prev: v,
			}
			v.next[newBlockHash] = s.leaf[newBlockHash]
			s.node[newBlockHash] = s.leaf[newBlockHash]
		}
	}
	s.update()
}

func (s *ViewGraph) update() {
	s.traverse(s.root)
	s.updateConfirmBlock(s.bestView)
}

func (s *ViewGraph) GetBestView() consensus.ChainViewInterface {
	return s.bestView.view
}

func (s *ViewGraph) GetFinalView() consensus.ChainViewInterface {
	return s.confirmView.view
}

func (s *ViewGraph) Print() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.edgeStr = ""
	s.traverse(s.root)

	dotContent := `digraph {
node [shape=record];
//    rankdir="LR";
newrank=true;
`
	maxTimeSlot := uint64(0)
	for k, v := range s.node {
		shortK := k.String()[0:5]
		dotContent += fmt.Sprintf(`%s_%d_%s [label = "%d:%s"]`, s.name, v.view.GetBlock().GetHeader().GetHeight(), string(shortK), v.view.GetBlock().GetHeader().GetHeight(), string(shortK)) + "\n"
		dotContent += fmt.Sprintf(`{rank=same; %s_%d_%s; slot_%d;}`, s.name, v.view.GetBlock().GetHeader().GetHeight(), string(shortK), v.view.GetBlock().GetHeader().(blockinterface.BlockHeaderV2Interface).GetTimeslot()-s.root.view.GetTimeslot()) + "\n"
		if v.view.GetTimeslot() > maxTimeSlot {
			maxTimeSlot = v.view.GetTimeslot()
		}
	}

	for i := s.root.view.GetTimeslot(); i < maxTimeSlot; i++ {
		dotContent += fmt.Sprintf("slot_%d -> slot_%d;", i-s.root.view.GetTimeslot(), i+1-s.root.view.GetTimeslot()) + "\n"
	}

	dotContent += s.edgeStr
	dotContent += `}`

	fd, _ := os.OpenFile(s.name+".dot", os.O_WRONLY|os.O_CREATE, 0666)
	fd.Truncate(0)
	fd.Write([]byte(dotContent))
	fd.Close()

}

func (s *ViewGraph) traverse(n *ViewNode) {
	if n.next != nil && len(n.next) != 0 {
		for h, v := range n.next {
			s.edgeStr += fmt.Sprintf("%s_%d_%s -> %s_%d_%s;\n", s.name, n.view.GetBlock().GetHeader().GetHeight(), string(n.view.GetBlock().GetHeader().GetHash().String()[0:5]), s.name, v.view.GetBlock().GetHeader().GetHeight(), string(h.String()[0:5]))
			s.traverse(v)
		}
	} else {
		if s.bestView == nil {
			s.bestView = n
		} else {
			if n.view.GetBlock().GetHeader().GetHeight() > s.bestView.view.GetBlock().GetHeader().GetHeight() {
				s.bestView = n
			}
			if (n.view.GetBlock().GetHeader().GetHeight() == s.bestView.view.GetBlock().GetHeader().GetHeight()) && n.view.GetBlock().GetHeader().GetTimestamp() < s.bestView.view.GetBlock().GetHeader().GetTimestamp() {
				s.bestView = n
			}
		}
	}
}

func (s *ViewGraph) updateConfirmBlock(node *ViewNode) {
	_1block := node.prev
	if _1block == nil {
		s.confirmView = node
		return
	}
	_2block := _1block.prev
	if _2block == nil {
		s.confirmView = _1block
		return
	}
	if _2block.view.GetTimeslot() == _1block.view.GetTimeslot()-1 && _2block.view.GetTimeslot() == node.view.GetTimeslot()-2 {
		s.confirmView = _2block
		return
	}
	s.updateConfirmBlock(_1block)
}
