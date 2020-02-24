package shardstate

import "github.com/incognitochain/incognito-chain/common"

type ShardState struct {
	Height     uint64
	Hash       common.Hash
	CrossShard []byte //In this state, shard i send cross shard tx to which shard
}
