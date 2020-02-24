package app

import (
	"encoding/json"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/shardblockv1"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/shardblockv2"
)

func UnmarshalShardBlock(data []byte) (blockinterface.ShardBlockInterface, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(string(data)), &jsonMap)
	if err != nil {
		panic(err)
	}
	if v, ok := jsonMap["Header"]; ok {
		header := v.(map[string]interface{})
		if version, ok := header["Version"]; ok {
			switch int(version.(float64)) {
			case blockchain.SHARD_BLOCK_VERSION:
				shardBlk := &shardblockv1.ShardBlock{}
				err := json.Unmarshal(data, &shardBlk)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return shardBlk, nil
			case blockchain.SHARD_BLOCK_VERSION2:
				shardBlk := &shardblockv2.ShardBlock{}
				err := json.Unmarshal(data, &shardBlk)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return shardBlk, nil
			}
		}
	}
	return nil, nil
}

func UnmarshalShardToBeaconBlock(data []byte) (blockinterface.ShardToBeaconBlockInterface, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(string(data)), &jsonMap)
	if err != nil {
		panic(err)
	}
	if v, ok := jsonMap["Header"]; ok {
		header := v.(map[string]interface{})
		if version, ok := header["Version"]; ok {
			switch int(version.(float64)) {
			case blockchain.SHARD_BLOCK_VERSION:
				s2bBlock := &shardblockv1.ShardToBeaconBlock{}
				err := json.Unmarshal(data, &s2bBlock)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return s2bBlock, nil
			case blockchain.SHARD_BLOCK_VERSION2:
				s2bBlock := &shardblockv2.ShardToBeaconBlock{}
				err := json.Unmarshal(data, &s2bBlock)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return s2bBlock, nil
			}
		}
	}
	return nil, nil
}

func UnmarshalCrossShardBlock(data []byte) (blockinterface.CrossShardBlockInterface, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(string(data)), &jsonMap)
	if err != nil {
		panic(err)
	}
	if v, ok := jsonMap["Header"]; ok {
		header := v.(map[string]interface{})
		if version, ok := header["Version"]; ok {
			switch int(version.(float64)) {
			case blockchain.SHARD_BLOCK_VERSION:
				shardBlk := &shardblockv1.CrossShardBlock{}
				err := json.Unmarshal(data, &shardBlk)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return shardBlk, nil
			case blockchain.SHARD_BLOCK_VERSION2:
				shardBlk := &shardblockv2.CrossShardBlock{}
				err := json.Unmarshal(data, &shardBlk)
				if err != nil {
					return nil, blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
				}
				return shardBlk, nil
			}
		}
	}
	return nil, nil
}
