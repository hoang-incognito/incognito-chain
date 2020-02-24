package beaconblockv2

import (
	"encoding/json"
	"strconv"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/shardstate"
	"github.com/incognitochain/incognito-chain/common"
)

type BeaconBody struct {
	// Shard State extract from shard to beacon block
	// Store all shard state == store content of all shard to beacon block
	ShardState   map[byte][]shardstate.ShardState
	Instructions [][]string
}

func (beaconBlock *BeaconBody) toString() string {
	res := ""
	for _, l := range beaconBlock.ShardState {
		for _, r := range l {
			res += strconv.Itoa(int(r.Height))
			res += r.Hash.String()
			crossShard, _ := json.Marshal(r.CrossShard)
			res += string(crossShard)

		}
	}
	for _, l := range beaconBlock.Instructions {
		for _, r := range l {
			res += r
		}
	}
	return res
}

func (beaconBody BeaconBody) Hash() common.Hash {
	return common.HashH([]byte(beaconBody.toString()))
}

func (beaconBody BeaconBody) GetShardState() map[byte][]shardstate.ShardState {
	return beaconBody.ShardState
}
func (beaconBody BeaconBody) GetInstructions() [][]string {
	return beaconBody.Instructions
}
