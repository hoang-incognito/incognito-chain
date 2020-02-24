package app

import (
	"encoding/json"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
)

func storeForSlashing(block blockinterface.BeaconBlockInterface, bc BlockChain) error {
	var err error
	db := bc.GetDatabase()
	beaconHeight := block.GetHeader().GetHeight()
	producersBlackList, err := db.GetProducersBlackList(beaconHeight - 1)
	if err != nil {
		return err
	}

	chainParamEpoch := bc.GetChainParams().Epoch
	newBeaconHeight := block.GetHeader().GetHeight()
	if newBeaconHeight%uint64(chainParamEpoch) == 0 { // end of epoch
		punishedProducersFinished := []string{}
		for producer := range producersBlackList {
			producersBlackList[producer]--
			if producersBlackList[producer] == 0 {
				punishedProducersFinished = append(punishedProducersFinished, producer)
			}
		}
		for _, producer := range punishedProducersFinished {
			delete(producersBlackList, producer)
		}
	}

	for _, inst := range block.GetBody().GetInstructions() {
		if len(inst) == 0 {
			continue
		}
		if inst[0] != blockchain.SwapAction {
			continue
		}
		badProducersWithPunishmentBytes := []byte{}
		if len(inst) == 6 && inst[3] == "shard" {
			badProducersWithPunishmentBytes = []byte(inst[5])
		}
		if len(inst) == 5 && inst[3] == "beacon" {
			badProducersWithPunishmentBytes = []byte(inst[4])
		}
		if len(badProducersWithPunishmentBytes) == 0 {
			continue
		}

		var badProducersWithPunishment map[string]uint8
		err = json.Unmarshal(badProducersWithPunishmentBytes, &badProducersWithPunishment)
		if err != nil {
			return err
		}
		for producer, punishedEpoches := range badProducersWithPunishment {
			epoches, found := producersBlackList[producer]
			if !found || epoches < punishedEpoches {
				producersBlackList[producer] = punishedEpoches
			}
		}
	}
	err = db.StoreProducersBlackList(beaconHeight, producersBlackList)
	return err
}
