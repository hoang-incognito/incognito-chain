package app

import (
	"context"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
)

type StoreShardDatabaseState struct {
	ctx     context.Context
	block   blockinterface.ShardBlockInterface
	bc      BlockChain
	curView *ShardView
	newView *ShardView
	//app
	app []ShardApp
}

func (s *ShardView) NewStoreDBState(ctx context.Context) *StoreShardDatabaseState {
	storeDBState := &StoreShardDatabaseState{
		ctx:     ctx,
		bc:      s.BC,
		curView: s,
		newView: s.CloneNewView().(*ShardView),
		app:     []ShardApp{},
	}

	//ADD YOUR APP HERE
	storeDBState.app = append(storeDBState.app, &ShardCoreApp{Logger: s.Logger, StoreState: storeDBState})
	storeDBState.app = append(storeDBState.app, &ShardPDEApp{Logger: s.Logger, StoreState: storeDBState})
	storeDBState.app = append(storeDBState.app, &ShardBridgeApp{Logger: s.Logger, StoreState: storeDBState})

	return storeDBState
}

func (s *ShardView) StoreDatabase(ctx context.Context) error {
	state := s.NewStoreDBState(context.Background())

	for _, app := range state.app {
		err := app.storeDatabase(state)
		if err != nil {
			//TODO: revert db state if get error
			return err
		}
	}
	return nil
}
