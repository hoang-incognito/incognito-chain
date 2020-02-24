package app

import (
	"encoding/binary"
	"encoding/json"
	"sort"

	"github.com/incognitochain/incognito-chain/database"
	"github.com/incognitochain/incognito-chain/database/lvdb"
	"github.com/pkg/errors"
)

func storePDEShares(
	db database.DatabaseInterface,
	beaconHeight uint64,
	pdeShares map[string]uint64,
) error {
	for shareKey, shareAmt := range pdeShares {
		newKey := replaceNewBCHeightInKeyStr(shareKey, beaconHeight)
		buf := make([]byte, binary.MaxVarintLen64)
		binary.LittleEndian.PutUint64(buf, shareAmt)
		dbErr := db.Put([]byte(newKey), buf)
		if dbErr != nil {
			return database.NewDatabaseError(database.AddShareAmountUpError, errors.Wrap(dbErr, "db.lvdb.put"))
		}
	}
	return nil
}

func storeWaitingPDEContributions(
	db database.DatabaseInterface,
	beaconHeight uint64,
	waitingPDEContributions map[string]*lvdb.PDEContribution,
) error {
	for contribKey, contribution := range waitingPDEContributions {
		newKey := replaceNewBCHeightInKeyStr(contribKey, beaconHeight)
		contributionBytes, err := json.Marshal(contribution)
		if err != nil {
			return err
		}
		err = db.Put([]byte(newKey), contributionBytes)
		if err != nil {
			return database.NewDatabaseError(database.StoreWaitingPDEContributionError, errors.Wrap(err, "db.lvdb.put"))
		}
	}
	return nil
}

func storePDEPoolPairs(
	db database.DatabaseInterface,
	beaconHeight uint64,
	pdePoolPairs map[string]*lvdb.PDEPoolForPair,
) error {
	for poolPairKey, poolPair := range pdePoolPairs {
		newKey := replaceNewBCHeightInKeyStr(poolPairKey, beaconHeight)
		poolPairBytes, err := json.Marshal(poolPair)
		if err != nil {
			return err
		}
		err = db.Put([]byte(newKey), poolPairBytes)
		if err != nil {
			return database.NewDatabaseError(database.StorePDEPoolForPairError, errors.Wrap(err, "db.lvdb.put"))
		}
	}
	return nil
}

func storePDEStateToDB(
	db database.DatabaseInterface,
	beaconHeight uint64,
	currentPDEState *CurrentPDEState,
) error {
	err := storeWaitingPDEContributions(db, beaconHeight, currentPDEState.WaitingPDEContributions)
	if err != nil {
		return err
	}
	err = storePDEPoolPairs(db, beaconHeight, currentPDEState.PDEPoolPairs)
	if err != nil {
		return err
	}
	err = storePDEShares(db, beaconHeight, currentPDEState.PDEShares)
	if err != nil {
		return err
	}
	return nil
}

func updateWaitingContributionPairToPoolV2(
	beaconHeight uint64,
	waitingContribution1 *lvdb.PDEContribution,
	waitingContribution2 *lvdb.PDEContribution,
	currentPDEState *CurrentPDEState,
) {
	addShareAmountUpV2(
		beaconHeight,
		waitingContribution1.TokenIDStr,
		waitingContribution2.TokenIDStr,
		waitingContribution1.TokenIDStr,
		waitingContribution1.ContributorAddressStr,
		waitingContribution1.Amount,
		currentPDEState,
	)

	waitingContributions := []*lvdb.PDEContribution{waitingContribution1, waitingContribution2}
	sort.Slice(waitingContributions, func(i, j int) bool {
		return waitingContributions[i].TokenIDStr < waitingContributions[j].TokenIDStr
	})
	pdePoolForPairKey := string(lvdb.BuildPDEPoolForPairKey(beaconHeight, waitingContributions[0].TokenIDStr, waitingContributions[1].TokenIDStr))
	pdePoolForPair, found := currentPDEState.PDEPoolPairs[pdePoolForPairKey]
	if !found || pdePoolForPair == nil {
		storePDEPoolForPair(
			pdePoolForPairKey,
			waitingContributions[0].TokenIDStr,
			waitingContributions[0].Amount,
			waitingContributions[1].TokenIDStr,
			waitingContributions[1].Amount,
			currentPDEState,
		)
		return
	}
	storePDEPoolForPair(
		pdePoolForPairKey,
		waitingContributions[0].TokenIDStr,
		pdePoolForPair.Token1PoolValue+waitingContributions[0].Amount,
		waitingContributions[1].TokenIDStr,
		pdePoolForPair.Token2PoolValue+waitingContributions[1].Amount,
		currentPDEState,
	)
}
