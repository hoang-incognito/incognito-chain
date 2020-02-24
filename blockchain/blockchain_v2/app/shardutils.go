package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/incognitochain/incognito-chain/blockchain"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/consensusheader"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/shardblockv1"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/shardblockv2"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/consensus_v2/blsbftv2"
	"github.com/incognitochain/incognito-chain/incognitokey"
	"github.com/incognitochain/incognito-chain/metadata"
	"github.com/incognitochain/incognito-chain/transaction"
	"github.com/pkg/errors"
)

const GENESIS_TIMESTAMP = "2006-01-02T15:04:05.000Z"

func CreateShardGenesisBlock(
	version int,
	net uint16,
	genesisBlockTime string,
	initTx []string,
) blockinterface.ShardBlockInterface {
	if version == 1 {
		body := shardblockv1.ShardBody{}
		genesisTime, err := time.Parse(GENESIS_TIMESTAMP, genesisBlockTime)
		if err != nil {
			fmt.Println(err)
		}
		header := shardblockv1.ShardHeader{
			Timestamp:         genesisTime.Unix(),
			Version:           version,
			BeaconHeight:      1,
			Epoch:             1,
			Round:             1,
			Height:            1,
			PreviousBlockHash: common.Hash{},
		}

		for _, tx := range initTx {
			testSalaryTX := transaction.Tx{}
			testSalaryTX.UnmarshalJSON([]byte(tx))
			body.Transactions = append(body.Transactions, &testSalaryTX)
		}

		block := &shardblockv1.ShardBlock{
			Body:   body,
			Header: header,
		}

		return block
	} else {
		body := shardblockv2.ShardBody{}
		genesisTime, err := time.Parse(GENESIS_TIMESTAMP, genesisBlockTime)
		if err != nil {
			fmt.Println(err)
		}
		header := shardblockv2.ShardHeader{
			Timestamp:         genesisTime.Unix(),
			Version:           version,
			BeaconHeight:      1,
			Epoch:             1,
			Round:             1,
			Height:            1,
			PreviousBlockHash: common.Hash{},
		}

		for _, tx := range initTx {
			testSalaryTX := transaction.Tx{}
			testSalaryTX.UnmarshalJSON([]byte(tx))
			body.Transactions = append(body.Transactions, &testSalaryTX)
		}

		block := &shardblockv2.ShardBlock{
			Body:   body,
			Header: header,
			ConsensusHeader: consensusheader.ConsensusHeader{
				TimeSlot: common.GetTimeSlot(genesisTime.Unix(), time.Now().Unix(), blsbftv2.TIMESLOT),
				Proposer: "",
			},
		}

		return block
	}
	return nil
}

/*
	Action Generate From Transaction:
	- Stake
		+ ["stake", "pubkey1,pubkey2,..." "shard" "txStake1,txStake2,..." "rewardReceiver1,rewardReceiver2,..." "flag1,flag2,..."]
		+ ["stake", "pubkey1,pubkey2,..." "beacon" "txStake1,txStake2,..." "rewardReceiver1,rewardReceiver2,..." "flag1,flag2,..."]
	- Stop Auto Staking
		+ ["stopautostaking" "pubkey1,pubkey2,..."]
*/
func CreateShardInstructionsFromTransactionAndInstruction(transactions []metadata.Transaction, bc metadata.BlockchainRetriever, shardID byte) (instructions [][]string, err error) {
	// Generate stake action
	stakeShardPublicKey := []string{}
	stakeBeaconPublicKey := []string{}
	stakeShardTxID := []string{}
	stakeBeaconTxID := []string{}
	stakeShardRewardReceiver := []string{}
	stakeBeaconRewardReceiver := []string{}
	stakeShardAutoStaking := []string{}
	stakeBeaconAutoStaking := []string{}
	stopAutoStaking := []string{}
	// @Notice: move build action from metadata into one loop
	//instructions, err = buildActionsFromMetadata(transactions, bc, shardID)
	//if err != nil {
	//	return nil, err
	//}
	for _, tx := range transactions {
		metadataValue := tx.GetMetadata()
		if metadataValue != nil {
			actionPairs, err := metadataValue.BuildReqActions(tx, bc, shardID)
			if err == nil {
				instructions = append(instructions, actionPairs...)
			}
		}
		switch tx.GetMetadataType() {
		case metadata.ShardStakingMeta:
			stakingMetadata, ok := tx.GetMetadata().(*metadata.StakingMetadata)
			if !ok {
				return nil, fmt.Errorf("Expect metadata type to be *metadata.StakingMetadata but get %+v", reflect.TypeOf(tx.GetMetadata()))
			}
			stakeShardPublicKey = append(stakeShardPublicKey, stakingMetadata.CommitteePublicKey)
			stakeShardTxID = append(stakeShardTxID, tx.Hash().String())
			stakeShardRewardReceiver = append(stakeShardRewardReceiver, stakingMetadata.RewardReceiverPaymentAddress)
			if stakingMetadata.AutoReStaking {
				stakeShardAutoStaking = append(stakeShardAutoStaking, "true")
			} else {
				stakeShardAutoStaking = append(stakeShardAutoStaking, "false")
			}
		case metadata.BeaconStakingMeta:
			stakingMetadata, ok := tx.GetMetadata().(*metadata.StakingMetadata)
			if !ok {
				return nil, fmt.Errorf("Expect metadata type to be *metadata.StakingMetadata but get %+v", reflect.TypeOf(tx.GetMetadata()))
			}
			stakeBeaconPublicKey = append(stakeBeaconPublicKey, stakingMetadata.CommitteePublicKey)
			stakeBeaconTxID = append(stakeBeaconTxID, tx.Hash().String())
			stakeBeaconRewardReceiver = append(stakeBeaconRewardReceiver, stakingMetadata.RewardReceiverPaymentAddress)
			if stakingMetadata.AutoReStaking {
				stakeBeaconAutoStaking = append(stakeBeaconAutoStaking, "true")
			} else {
				stakeBeaconAutoStaking = append(stakeBeaconAutoStaking, "false")
			}
		case metadata.StopAutoStakingMeta:
			{
				stopAutoStakingMetadata, ok := tx.GetMetadata().(*metadata.StopAutoStakingMetadata)
				if !ok {
					return nil, fmt.Errorf("Expect metadata type to be *metadata.StopAutoStakingMetadata but get %+v", reflect.TypeOf(tx.GetMetadata()))
				}
				stopAutoStaking = append(stopAutoStaking, stopAutoStakingMetadata.CommitteePublicKey)
			}
		}
	}
	if !reflect.DeepEqual(stakeShardPublicKey, []string{}) {
		if len(stakeShardPublicKey) != len(stakeShardTxID) && len(stakeShardTxID) != len(stakeShardRewardReceiver) && len(stakeShardRewardReceiver) != len(stakeShardAutoStaking) {
			return nil, blockchain.NewBlockChainError(blockchain.StakeInstructionError, fmt.Errorf("Expect public key list (length %+v) and reward receiver list (length %+v), auto restaking (length %+v) to be equal", len(stakeShardPublicKey), len(stakeShardRewardReceiver), len(stakeShardAutoStaking)))
		}
		stakeShardPublicKey, err = incognitokey.ConvertToBase58ShortFormat(stakeShardPublicKey)
		if err != nil {
			return nil, fmt.Errorf("Failed To Convert Stake Shard Public Key to Base58 Short Form")
		}
		// ["stake", "pubkey1,pubkey2,..." "shard" "txStake1,txStake2,..." "rewardReceiver1,rewardReceiver2,..." "flag1,flag2,..."]
		instruction := []string{blockchain.StakeAction, strings.Join(stakeShardPublicKey, ","), "shard", strings.Join(stakeShardTxID, ","), strings.Join(stakeShardRewardReceiver, ","), strings.Join(stakeShardAutoStaking, ",")}
		instructions = append(instructions, instruction)
	}
	if !reflect.DeepEqual(stakeBeaconPublicKey, []string{}) {
		if len(stakeBeaconPublicKey) != len(stakeBeaconTxID) && len(stakeBeaconTxID) != len(stakeBeaconRewardReceiver) && len(stakeBeaconRewardReceiver) != len(stakeBeaconAutoStaking) {
			return nil, blockchain.NewBlockChainError(blockchain.StakeInstructionError, fmt.Errorf("Expect public key list (length %+v) and reward receiver list (length %+v), auto restaking (length %+v) to be equal", len(stakeBeaconPublicKey), len(stakeBeaconRewardReceiver), len(stakeBeaconAutoStaking)))
		}
		stakeBeaconPublicKey, err = incognitokey.ConvertToBase58ShortFormat(stakeBeaconPublicKey)
		if err != nil {
			return nil, fmt.Errorf("Failed To Convert Stake Beacon Public Key to Base58 Short Form")
		}
		// ["stake", "pubkey1,pubkey2,..." "beacon" "txStake1,txStake2,..." "rewardReceiver1,rewardReceiver2,..." "flag1,flag2,..."]
		instruction := []string{blockchain.StakeAction, strings.Join(stakeBeaconPublicKey, ","), "beacon", strings.Join(stakeBeaconTxID, ","), strings.Join(stakeBeaconRewardReceiver, ","), strings.Join(stakeBeaconAutoStaking, ",")}
		instructions = append(instructions, instruction)
	}
	if !reflect.DeepEqual(stopAutoStaking, []string{}) {
		// ["stopautostaking" "pubkey1,pubkey2,..."]
		instruction := []string{blockchain.StopAutoStake, strings.Join(stopAutoStaking, ",")}
		instructions = append(instructions, instruction)
	}
	return instructions, nil
}

func checkReturnStakingTxExistence(txId string, shardBlock blockinterface.ShardBlockInterface) bool {
	for _, tx := range shardBlock.GetShardBody().GetTransactions() {
		if tx.GetMetadata() != nil {
			if tx.GetMetadata().GetType() == metadata.ReturnStakingMeta {
				if returnStakingMeta, ok := tx.GetMetadata().(*metadata.ReturnStakingMetadata); ok {
					if returnStakingMeta.TxID == txId {
						return true
					}
				}
			}
		}
	}
	return false
}

/*
	Create Swap Action
	Return param:
	#1: swap instruction
	#2: new pending validator list after swapped
	#3: new committees after swapped
	#4: error
*/
func CreateSwapAction(
	pendingValidator []string,
	commitees []string,
	maxCommitteeSize int,
	minCommitteeSize int,
	shardID byte,
	producersBlackList map[string]uint8,
	badProducersWithPunishment map[string]uint8,
	offset int,
	swapOffset int,
) ([]string, []string, []string, error) {
	newPendingValidator, newShardCommittees, shardSwapedCommittees, shardNewCommittees, err := SwapValidator(pendingValidator, commitees, maxCommitteeSize, minCommitteeSize, offset, producersBlackList, swapOffset)
	if err != nil {
		return nil, nil, nil, err
	}
	badProducersWithPunishmentBytes, err := json.Marshal(badProducersWithPunishment)
	if err != nil {
		return nil, nil, nil, err
	}
	swapInstruction := []string{"swap", strings.Join(shardNewCommittees, ","), strings.Join(shardSwapedCommittees, ","), "shard", strconv.Itoa(int(shardID)), string(badProducersWithPunishmentBytes)}
	return swapInstruction, newPendingValidator, newShardCommittees, nil
}

// pickInstructionFromBeaconBlocks extracts all instructions of a specific type
func pickInstructionFromBeaconBlocks(beaconBlocks []blockinterface.BeaconBlockInterface, instType string) [][]string {
	insts := [][]string{}
	for _, block := range beaconBlocks {
		found := pickInstructionWithType(block.GetBeaconBody().GetInstructions(), instType)
		if len(found) > 0 {
			insts = append(insts, found...)
		}
	}
	return insts
}

func CreateShardToBeaconBlock(shardBlock blockinterface.ShardBlockInterface, bc *blockchain.BlockChain) (blockinterface.ShardToBeaconBlockInterface, error) {
	blockVersion := shardBlock.GetHeader().GetVersion()

	instructions, err := CreateShardInstructionsFromTransactionAndInstruction(shardBlock.GetShardBody().GetTransactions(), bc, shardBlock.GetShardHeader().GetShardID())
	if err != nil {
		return nil, err
	}

	switch blockVersion {
	case blockchain.SHARD_BLOCK_VERSION:
		s2bBlock := &shardblockv1.ShardToBeaconBlock{}
		s2bBlock.ValidationData = shardBlock.GetValidationField()
		s2bBlock.Header = shardBlock.GetShardHeader().(shardblockv1.ShardHeader)
		blockInstructions := shardBlock.GetShardBody().GetInstructions()
		s2bBlock.Instructions = append(instructions, blockInstructions...)
		return s2bBlock, nil
	case blockchain.SHARD_BLOCK_VERSION2:
		s2bBlock := &shardblockv2.ShardToBeaconBlock{}
		return s2bBlock, nil
	default:
		return nil, blockchain.NewBlockChainError(blockchain.CreateShardToBeaconBlockError, nil)
	}
}

func CreateAllCrossShardBlock(activeShards int, shardBlock blockinterface.ShardBlockInterface) (map[byte]blockinterface.CrossShardBlockInterface, error) {
	allCrossShard := make(map[byte]blockinterface.CrossShardBlockInterface)
	if activeShards == 1 {
		return allCrossShard, nil
	}
	for i := 0; i < activeShards; i++ {
		shardID := common.GetShardIDFromLastByte(byte(i))
		if shardID != shardBlock.GetShardHeader().GetShardID() {
			crossShard, err := CreateCrossShardBlock(shardBlock, shardID)
			if crossShard != nil {
				Logger.log.Infof("Create CrossShardBlock from Shard %+v to Shard %+v: %+v \n", shardBlock.GetShardHeader().GetShardID(), shardID, crossShard)
			}
			if err != nil {
				return nil, err
			}
			if crossShard != nil && err == nil {
				allCrossShard[byte(i)] = crossShard
			}
		}
	}
	return allCrossShard, nil
}

func CreateCrossShardBlock(shardBlock blockinterface.ShardBlockInterface, shardID byte) (blockinterface.CrossShardBlockInterface, error) {
	blockVersion := shardBlock.GetHeader().GetVersion()
	crossOutputCoin, crossCustomTokenPrivacyData := blockchain.GetCrossShardData(shardBlock.GetShardBody().GetTransactions(), shardID)
	// Return nothing if nothing to cross
	if len(crossOutputCoin) == 0 && len(crossCustomTokenPrivacyData) == 0 {
		return nil, blockchain.NewBlockChainError(blockchain.CreateCrossShardBlockError, errors.New("No cross Outputcoin, Cross Custom Token, Cross Custom Token Privacy"))
	}
	merklePathShard, merkleShardRoot := blockchain.GetMerklePathCrossShard2(shardBlock.GetShardBody().GetTransactions(), shardID)
	if merkleShardRoot != shardBlock.GetShardHeader().GetShardTxRoot() {
		return nil, blockchain.NewBlockChainError(blockchain.VerifyCrossShardBlockShardTxRootError, fmt.Errorf("Expect Shard Tx Root To be %+v but get %+v", shardBlock.GetShardHeader().GetShardTxRoot(), merkleShardRoot))
	}

	switch blockVersion {
	case blockchain.SHARD_BLOCK_VERSION:
		crossShard := &shardblockv1.CrossShardBlock{}
		crossShard.ValidationData = shardBlock.GetValidationField()
		crossShard.Header = shardBlock.GetHeader().(shardblockv1.ShardHeader)
		crossShard.MerklePathShard = merklePathShard
		crossShard.CrossOutputCoin = crossOutputCoin
		crossShard.CrossTxTokenPrivacyData = crossCustomTokenPrivacyData
		crossShard.ToShardID = shardID
		return crossShard, nil
	case blockchain.SHARD_BLOCK_VERSION2:
		crossShard := &shardblockv2.CrossShardBlock{}

		return crossShard, nil
	default:
		return nil, blockchain.NewBlockChainError(blockchain.CreateCrossShardBlockError, nil)
	}
}
