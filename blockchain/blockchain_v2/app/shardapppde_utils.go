package app

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/incognitochain/incognito-chain/database"
	"github.com/incognitochain/incognito-chain/database/lvdb"
)

type CurrentPDEState struct {
	WaitingPDEContributions map[string]*lvdb.PDEContribution
	PDEPoolPairs            map[string]*lvdb.PDEPoolForPair
	PDEShares               map[string]uint64
}

type DeductingAmountsByWithdrawal struct {
	Token1IDStr string
	PoolValue1  uint64
	Token2IDStr string
	PoolValue2  uint64
	Shares      uint64
}

func replaceNewBCHeightInKeyStr(key string, newBeaconHeight uint64) string {
	parts := strings.Split(key, "-")
	if len(parts) <= 1 {
		return key
	}
	parts[1] = fmt.Sprintf("%d", newBeaconHeight)
	newKey := ""
	for idx, part := range parts {
		if idx == len(parts)-1 {
			newKey += part
			continue
		}
		newKey += (part + "-")
	}
	return newKey
}

func getWaitingPDEContributions(
	db database.DatabaseInterface,
	beaconHeight uint64,
) (map[string]*lvdb.PDEContribution, error) {
	waitingPDEContributions := make(map[string]*lvdb.PDEContribution)
	waitingContribKeysBytes, waitingContribValuesBytes, err := db.GetAllRecordsByPrefix(beaconHeight, lvdb.WaitingPDEContributionPrefix)
	if err != nil {
		return nil, err
	}
	for idx, waitingContribKeyBytes := range waitingContribKeysBytes {
		var waitingContrib lvdb.PDEContribution
		err = json.Unmarshal(waitingContribValuesBytes[idx], &waitingContrib)
		if err != nil {
			return nil, err
		}
		waitingPDEContributions[string(waitingContribKeyBytes)] = &waitingContrib
	}
	return waitingPDEContributions, nil
}

func getPDEPoolPair(
	db database.DatabaseInterface,
	beaconHeight uint64,
) (map[string]*lvdb.PDEPoolForPair, error) {
	pdePoolPairs := make(map[string]*lvdb.PDEPoolForPair)
	poolPairsKeysBytes, poolPairsValuesBytes, err := db.GetAllRecordsByPrefix(beaconHeight, lvdb.PDEPoolPrefix)
	if err != nil {
		return nil, err
	}
	for idx, poolPairsKeyBytes := range poolPairsKeysBytes {
		var padePoolPair lvdb.PDEPoolForPair
		err = json.Unmarshal(poolPairsValuesBytes[idx], &padePoolPair)
		if err != nil {
			return nil, err
		}
		pdePoolPairs[string(poolPairsKeyBytes)] = &padePoolPair
	}
	return pdePoolPairs, nil
}

func getPDEShares(
	db database.DatabaseInterface,
	beaconHeight uint64,
) (map[string]uint64, error) {
	pdeShares := make(map[string]uint64)
	sharesKeysBytes, sharesValuesBytes, err := db.GetAllRecordsByPrefix(beaconHeight, lvdb.PDESharePrefix)
	if err != nil {
		return nil, err
	}
	for idx, sharesKeyBytes := range sharesKeysBytes {
		shareAmt := uint64(binary.LittleEndian.Uint64(sharesValuesBytes[idx]))
		pdeShares[string(sharesKeyBytes)] = shareAmt
	}
	return pdeShares, nil
}

func InitCurrentPDEStateFromDB(
	db database.DatabaseInterface,
	beaconHeight uint64,
) (*CurrentPDEState, error) {
	waitingPDEContributions, err := getWaitingPDEContributions(db, beaconHeight)
	if err != nil {
		return nil, err
	}
	pdePoolPairs, err := getPDEPoolPair(db, beaconHeight)
	if err != nil {
		return nil, err
	}
	pdeShares, err := getPDEShares(db, beaconHeight)
	if err != nil {
		return nil, err
	}
	return &CurrentPDEState{
		WaitingPDEContributions: waitingPDEContributions,
		PDEPoolPairs:            pdePoolPairs,
		PDEShares:               pdeShares,
	}, nil
}

func addShareAmountUpV2(
	beaconHeight uint64,
	token1IDStr string,
	token2IDStr string,
	contributedTokenIDStr string,
	contributorAddrStr string,
	amt uint64,
	currentPDEState *CurrentPDEState,
) {
	pdeShareOnTokenPrefix := string(lvdb.BuildPDESharesKeyV2(beaconHeight, token1IDStr, token2IDStr, ""))
	totalSharesOnToken := uint64(0)
	for key, value := range currentPDEState.PDEShares {
		if strings.Contains(key, pdeShareOnTokenPrefix) {
			totalSharesOnToken += value
		}
	}
	pdeShareKey := string(lvdb.BuildPDESharesKeyV2(beaconHeight, token1IDStr, token2IDStr, contributorAddrStr))
	if totalSharesOnToken == 0 {
		currentPDEState.PDEShares[pdeShareKey] = amt
		return
	}
	poolPairKey := string(lvdb.BuildPDEPoolForPairKey(beaconHeight, token1IDStr, token2IDStr))
	poolPair, found := currentPDEState.PDEPoolPairs[poolPairKey]
	if !found || poolPair == nil {
		currentPDEState.PDEShares[pdeShareKey] = amt
		return
	}
	poolValue := poolPair.Token1PoolValue
	if poolPair.Token2IDStr == contributedTokenIDStr {
		poolValue = poolPair.Token2PoolValue
	}
	if poolValue == 0 {
		currentPDEState.PDEShares[pdeShareKey] = amt
	}
	increasingAmt := big.NewInt(0)
	increasingAmt.Mul(big.NewInt(int64(totalSharesOnToken)), big.NewInt(int64(amt)))
	increasingAmt.Div(increasingAmt, big.NewInt(int64(poolValue)))
	currentShare, found := currentPDEState.PDEShares[pdeShareKey]
	addedUpAmt := increasingAmt.Uint64()
	if found {
		addedUpAmt += currentShare
	}
	currentPDEState.PDEShares[pdeShareKey] = addedUpAmt
}
