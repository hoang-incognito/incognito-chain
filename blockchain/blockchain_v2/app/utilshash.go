package app

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/shardstate"
	"github.com/incognitochain/incognito-chain/common"
)

//=========================HASH util==================================
func GenerateZeroValueHash() (common.Hash, error) {
	hash := common.Hash{}
	hash.SetBytes(make([]byte, 32))
	return hash, nil
}
func GenerateHashFromHashArray(hashes []common.Hash) (common.Hash, error) {
	// if input is empty list
	// return hash value of bytes zero
	if len(hashes) == 0 {
		return GenerateZeroValueHash()
	}
	strs := []string{}
	for _, value := range hashes {
		str := value.String()
		strs = append(strs, str)
	}
	return GenerateHashFromStringArray(strs)
}

func GenerateHashFromStringArray(strs []string) (common.Hash, error) {
	// if input is empty list
	// return hash value of bytes zero
	if len(strs) == 0 {
		return GenerateZeroValueHash()
	}
	var (
		hash common.Hash
		buf  bytes.Buffer
	)
	for _, value := range strs {
		buf.WriteString(value)
	}
	temp := common.HashB(buf.Bytes())
	if err := hash.SetBytes(temp[:]); err != nil {
		return common.Hash{}, blockchain.NewBlockChainError(blockchain.HashError, err)
	}
	return hash, nil
}

func GenerateHashFromMapByteString(maps1 map[byte][]string, maps2 map[byte][]string) (common.Hash, error) {
	var keys1 []int
	for k := range maps1 {
		keys1 = append(keys1, int(k))
	}
	sort.Ints(keys1)
	shardPendingValidator := []string{}
	// To perform the opertion you want
	for _, k := range keys1 {
		shardPendingValidator = append(shardPendingValidator, maps1[byte(k)]...)
	}

	var keys2 []int
	for k := range maps2 {
		keys2 = append(keys2, int(k))
	}
	sort.Ints(keys2)
	shardValidator := []string{}
	// To perform the opertion you want
	for _, k := range keys2 {
		shardValidator = append(shardValidator, maps2[byte(k)]...)
	}
	return GenerateHashFromStringArray(append(shardPendingValidator, shardValidator...))
}

func GenerateHashFromMapStringString(maps1 map[string]string) (common.Hash, error) {
	var keys []string
	var res []string
	for k := range maps1 {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		res = append(res, key)
		res = append(res, maps1[key])
	}
	return GenerateHashFromStringArray(res)
}
func GenerateHashFromMapStringBool(maps1 map[string]bool) (common.Hash, error) {
	var keys []string
	var res []string
	for k := range maps1 {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		res = append(res, key)
		if maps1[key] {
			res = append(res, "true")
		} else {
			res = append(res, "false")
		}
	}
	return GenerateHashFromStringArray(res)
}
func GenerateHashFromShardState(allShardState map[byte][]shardstate.ShardState) (common.Hash, error) {
	allShardStateStr := []string{}
	var keys []int
	for k := range allShardState {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, shardID := range keys {
		res := ""
		for _, shardState := range allShardState[byte(shardID)] {
			res += strconv.Itoa(int(shardState.Height))
			res += shardState.Hash.String()
			crossShard, _ := json.Marshal(shardState.CrossShard)
			res += string(crossShard)
		}
		allShardStateStr = append(allShardStateStr, res)
	}
	return GenerateHashFromStringArray(allShardStateStr)
}
func generateLastCrossShardStateHash(lastCrossShardState map[byte]map[byte]uint64) common.Hash {
	res := ""
	var fromKeys = []int{}
	for key, _ := range lastCrossShardState {
		fromKeys = append(fromKeys, int(key))
	}
	sort.Ints(fromKeys)
	for _, fromKey := range fromKeys {
		fromShardID := byte(fromKey)
		toCrossShardState := lastCrossShardState[fromShardID]
		var toKeys = []int{}
		for key, _ := range toCrossShardState {
			toKeys = append(toKeys, int(key))
		}
		sort.Ints(toKeys)
		for _, toKey := range toKeys {
			toShardID := byte(toKey)
			lastHeight := toCrossShardState[toShardID]
			res += strconv.Itoa(int(lastHeight))
		}
	}
	return common.HashH([]byte(res))
}
func VerifyHashFromHashArray(hashes []common.Hash, hash common.Hash) (common.Hash, bool) {
	strs := []string{}
	for _, value := range hashes {
		str := value.String()
		strs = append(strs, str)
	}
	return VerifyHashFromStringArray(strs, hash)
}

func VerifyHashFromStringArray(strs []string, hash common.Hash) (common.Hash, bool) {
	res, err := GenerateHashFromStringArray(strs)
	if err != nil {
		return common.Hash{}, false
	}
	return res, bytes.Equal(res.GetBytes(), hash.GetBytes())
}

func VerifyHashFromMapByteString(maps1 map[byte][]string, maps2 map[byte][]string, hash common.Hash) bool {
	res, err := GenerateHashFromMapByteString(maps1, maps2)
	if err != nil {
		return false
	}
	return bytes.Equal(res.GetBytes(), hash.GetBytes())
}

func VerifyHashFromShardState(allShardState map[byte][]shardstate.ShardState, hash common.Hash) bool {
	res, err := GenerateHashFromShardState(allShardState)
	if err != nil {
		return false
	}
	return bytes.Equal(res.GetBytes(), hash.GetBytes())
}

// NOTICE: this function is deprecate, just return empty data
func CalHashFromTxTokenDataList() (common.Hash, error) {
	hashes := []common.Hash{}
	hash, err := GenerateHashFromHashArray(hashes)
	if err != nil {
		return common.Hash{}, err
	}
	return hash, nil
}
func VerifyLastCrossShardStateHash(lastCrossShardState map[byte]map[byte]uint64, targetHash common.Hash) (common.Hash, bool) {
	hash := generateLastCrossShardStateHash(lastCrossShardState)
	return hash, hash.IsEqual(&targetHash)
}
func VerifyHashFromMapStringString(maps1 map[string]string, targetHash common.Hash) (common.Hash, bool) {
	hash, err := GenerateHashFromMapStringString(maps1)
	if err != nil {
		return hash, false
	}
	return hash, hash.IsEqual(&targetHash)
}
func VerifyHashFromMapStringBool(maps1 map[string]bool, targetHash common.Hash) (common.Hash, bool) {
	hash, err := GenerateHashFromMapStringBool(maps1)
	if err != nil {
		return hash, false
	}
	return hash, hash.IsEqual(&targetHash)
}
