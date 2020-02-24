package beaconblockv2

import (
	"fmt"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/consensusheader"
	"github.com/incognitochain/incognito-chain/common"
)

type BeaconBlock struct {
	ConsensusHeader consensusheader.ConsensusHeader
	Body            BeaconBody
	Header          BeaconHeader
}

func (beaconBlock BeaconBlock) GetBlockType() string {
	return "beacon"
}

func (beaconBlock BeaconBlock) GetBeaconHeight() uint64 {
	return beaconBlock.Header.Height
}

func (beaconBlock BeaconBlock) GetPreviousBlockHash() common.Hash {
	return beaconBlock.Header.PreviousBlockHash
}

func (beaconBlock BeaconBlock) GetTimeslot() uint64 {
	return beaconBlock.ConsensusHeader.TimeSlot
}

func (beaconBlock BeaconBlock) GetCreateTimeslot() uint64 {
	return beaconBlock.Header.TimeSlot
}

func (beaconBlock BeaconBlock) GetBlockTimestamp() int64 {
	return beaconBlock.Header.Timestamp
}

func (beaconBlock BeaconBlock) GetHash() *common.Hash {
	return beaconBlock.Header.GetHash()
}

func (beaconBlock BeaconBlock) GetEpoch() uint64 {
	return beaconBlock.Header.Epoch
}

func (beaconBlock BeaconBlock) GetHeight() uint64 {
	return beaconBlock.Header.Height
}

func (beaconBlock *BeaconBlock) AddValidationField(validationData string) error {
	beaconBlock.ConsensusHeader.ValidationData = validationData
	return nil
}
func (beaconBlock BeaconBlock) GetValidationField() string {
	return beaconBlock.ConsensusHeader.ValidationData
}

func (beaconBlock BeaconBlock) GetRound() int {
	return beaconBlock.Header.Round
}
func (beaconBlock BeaconBlock) GetRoundKey() string {
	return fmt.Sprint(beaconBlock.Header.Height, "_", beaconBlock.Header.Round)
}

func (beaconBlock BeaconBlock) GetInstructions() [][]string {
	return beaconBlock.Body.Instructions
}

func (beaconBlock BeaconBlock) GetProducer() string {
	return beaconBlock.Header.Producer
}

func (beaconBlock BeaconBlock) GetConsensusType() string {
	return beaconBlock.Header.ConsensusType
}

func (beaconBlock BeaconBlock) GetHeader() blockinterface.BlockHeaderInterface {
	return beaconBlock.Header
}
func (beaconBlock BeaconBlock) GetBody() blockinterface.BlockBodyInterface {
	return beaconBlock.Body
}

func (beaconBlock BeaconBlock) GetBeaconHeader() blockinterface.BeaconHeaderInterface {
	return beaconBlock.Header
}
func (beaconBlock BeaconBlock) GetBeaconBody() blockinterface.BeaconBodyInterface {
	return beaconBlock.Body
}

func (beaconBlock BeaconBlock) GetVersion() int {
	return beaconBlock.Header.Version
}
