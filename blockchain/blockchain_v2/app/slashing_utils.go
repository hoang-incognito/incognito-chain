package app

import (
	"sort"

	"github.com/incognitochain/incognito-chain/blockchain"
)

func getUpdatedProducersBlackList(
	isBeacon bool,
	shardID int,
	committee []string,
	beaconHeight uint64,
	blockchain BlockChain,
) (map[string]uint8, error) {
	db := blockchain.GetDatabase()
	producersBlackList, err := db.GetProducersBlackList(beaconHeight)
	if err != nil {
		return nil, err
	}
	if isBeacon {
		punishedProducersFinished := []string{}
		for producer, punishedEpoches := range producersBlackList {
			if punishedEpoches == 1 {
				punishedProducersFinished = append(punishedProducersFinished, producer)
			}
		}
		for _, producer := range punishedProducersFinished {
			delete(producersBlackList, producer)
		}
	}

	badProducersWithPunishment := buildBadProducersWithPunishment(isBeacon, shardID, committee, blockchain)
	for producer, punishedEpoches := range badProducersWithPunishment {
		epoches, found := producersBlackList[producer]
		if !found || epoches < punishedEpoches {
			producersBlackList[producer] = punishedEpoches
		}
	}
	return sortMapStringUint8Keys(producersBlackList), nil
}

func buildBadProducersWithPunishment(
	isBeacon bool,
	shardID int,
	committee []string,
	bc BlockChain,
) map[string]uint8 {
	slashLevels := bc.GetChainParams().SlashLevels
	numOfBlocksByProducers := map[string]uint64{}
	//TODO: refactor this
	// if isBeacon {
	// 	numOfBlocksByProducers = bc.BestState.Beacon.NumOfBlocksByProducers
	// } else {
	// 	numOfBlocksByProducers = bc.BestState.Shard[byte(shardID)].NumOfBlocksByProducers
	// }
	// numBlkPerEpoch := blockchain.config.ChainParams.Epoch
	numBlkPerEpoch := uint64(0)
	for _, numBlk := range numOfBlocksByProducers {
		numBlkPerEpoch += numBlk
	}
	badProducersWithPunishment := make(map[string]uint8)
	committeeLen := len(committee)
	if committeeLen == 0 {
		return badProducersWithPunishment
	}
	expectedNumBlkByEachProducer := numBlkPerEpoch / uint64(committeeLen)

	if expectedNumBlkByEachProducer == 0 {
		return badProducersWithPunishment
	}
	// for producer, numBlk := range numOfBlocksByProducers {
	for _, producer := range committee {
		numBlk, found := numOfBlocksByProducers[producer]
		if !found {
			numBlk = 0
		}
		if numBlk >= expectedNumBlkByEachProducer {
			continue
		}
		missingPercent := uint8((-(numBlk - expectedNumBlkByEachProducer) * 100) / expectedNumBlkByEachProducer)
		var selectedSlLev *blockchain.SlashLevel
		for _, slLev := range slashLevels {
			if missingPercent >= slLev.MinRange {
				selectedSlLev = &slLev
			}
		}
		if selectedSlLev != nil {
			badProducersWithPunishment[producer] = selectedSlLev.PunishedEpoches
		}
	}
	return sortMapStringUint8Keys(badProducersWithPunishment)
}

func sortMapStringUint8Keys(m map[string]uint8) map[string]uint8 {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedMap := make(map[string]uint8)
	for _, k := range keys {
		sortedMap[k] = m[k]
	}
	return sortedMap
}
