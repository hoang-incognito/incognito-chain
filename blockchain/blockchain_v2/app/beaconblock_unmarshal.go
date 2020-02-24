package app

import (
	"encoding/json"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/beaconblockv1"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/beaconblockv2"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
)

// this convert beaconblock v1 to v2

func UnmarshalBeaconBlock(data []byte) (blockinterface.BeaconBlockInterface, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(string(data)), &jsonMap)
	if err != nil {
		panic(err)
	}
	if v, ok := jsonMap["Header"]; ok {
		header := v.(map[string]interface{})
		if version, ok := header["Version"]; ok {
			switch int(version.(float64)) {
			case blockchain.BEACON_BLOCK_VERSION:
				beaconBlk := &beaconblockv1.BeaconBlock{}
				err := json.Unmarshal(data, &beaconBlk)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return beaconBlk, nil
			case blockchain.BEACON_BLOCK_VERSION2:
				beaconBlk := &beaconblockv2.BeaconBlock{}
				err := json.Unmarshal(data, &beaconBlk)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return beaconBlk, nil
			}
		}
	}

	return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonBeaconBlockError, nil)
}
