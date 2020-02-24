package shardblockv2

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/incognitochain/incognito-chain/blockchain"

	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/blockinterface"
	"github.com/incognitochain/incognito-chain/blockchain/blockchain_v2/types/consensusheader"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/metadata"
)

type ShardBlock struct {
	ConsensusHeader consensusheader.ConsensusHeader
	Body            ShardBody
	Header          ShardHeader
}

func NewShardBlock() *ShardBlock {
	return &ShardBlock{
		Header: ShardHeader{},
		Body: ShardBody{
			Instructions:      [][]string{},
			CrossTransactions: make(map[byte][]blockchain.CrossTransaction),
			Transactions:      make([]metadata.Transaction, 0),
		},
	}
}
func NewShardBlockWithHeader(header ShardHeader) *ShardBlock {
	return &ShardBlock{
		Header: header,
		Body: ShardBody{
			Instructions:      [][]string{},
			CrossTransactions: make(map[byte][]blockchain.CrossTransaction),
			Transactions:      make([]metadata.Transaction, 0),
		},
	}
}
func NewShardBlockWithBody(body ShardBody) *ShardBlock {
	return &ShardBlock{
		Header: ShardHeader{},
		Body:   body,
	}
}
func NewShardBlockFull(header ShardHeader, body ShardBody) *ShardBlock {
	return &ShardBlock{
		Header: header,
		Body:   body,
	}
}
func (shardBlock *ShardBlock) BuildShardBlockBody(instructions [][]string, crossTransaction map[byte][]blockchain.CrossTransaction, transactions []metadata.Transaction) {
	shardBlock.Body.Instructions = append(shardBlock.Body.Instructions, instructions...)
	shardBlock.Body.CrossTransactions = crossTransaction
	shardBlock.Body.Transactions = append(shardBlock.Body.Transactions, transactions...)
}

func (shardBlock *ShardBlock) validateSanityData() (bool, error) {
	//TODO: after simulation, remove this
	return true, nil
	//Check Header
	if shardBlock.Header.Height == 1 && len(shardBlock.ConsensusHeader.ValidationData) != 0 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, errors.New("Expect Shard Block with Height 1 to not have validationData"))
	}
	// producer address must have 66 bytes: 33-byte public key, 33-byte transmission key
	if shardBlock.Header.Height > 1 && len(shardBlock.ConsensusHeader.ValidationData) == 0 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, errors.New("Expect Shard Block to have validationData"))
	}
	if int(shardBlock.Header.ShardID) < 0 || int(shardBlock.Header.ShardID) > 256 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block ShardID in range 0 - 255 but get %+v ", shardBlock.Header.ShardID))
	}
	if shardBlock.Header.Version < blockchain.SHARD_BLOCK_VERSION {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Version greater or equal than %+v but get %+v ", blockchain.SHARD_BLOCK_VERSION, shardBlock.Header.Version))
	}
	if len(shardBlock.Header.PreviousBlockHash[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Previous Hash in the right format"))
	}
	if shardBlock.Header.Height < 1 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Height to be greater than 0"))
	}
	if shardBlock.Header.Height == 1 && !shardBlock.Header.PreviousBlockHash.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block with Height 1 (first block) have Zero Hash Value"))
	}
	if shardBlock.Header.Height > 1 && shardBlock.Header.PreviousBlockHash.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block with Height greater than 1 have Non-Zero Hash Value"))
	}
	if shardBlock.Header.Round < 1 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Round greater or equal than 1"))
	}
	if shardBlock.Header.Epoch < 1 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Epoch greater or equal than 1"))
	}
	if shardBlock.Header.Timestamp < 0 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Epoch greater or equal than 0"))
	}
	if len(shardBlock.Header.TxRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Tx Root in the right format"))
	}
	if len(shardBlock.Header.ShardTxRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Shard Tx Root in the right format"))
	}
	if len(shardBlock.Header.CrossTransactionRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Cross Transaction Root in the right format"))
	}
	if len(shardBlock.Header.InstructionsRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Instructions Root in the right format"))
	}
	if len(shardBlock.Header.CommitteeRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Committee Root in the right format"))
	}
	if shardBlock.Header.Height == 1 && !shardBlock.Header.CommitteeRoot.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block with Height 1 have Zero Hash Value"))
	}
	if shardBlock.Header.Height > 1 && shardBlock.Header.CommitteeRoot.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block with Height greater than 1 have Non-Zero Hash Value"))
	}
	if len(shardBlock.Header.PendingValidatorRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Committee Root in the right format"))
	}
	if len(shardBlock.Header.StakingTxRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Staking Tx Root in the right format"))
	}
	if len(shardBlock.Header.CrossShardBitMap) > 254 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Cross Shard Length Less Than 255"))
	}
	if shardBlock.Header.BeaconHeight < 1 {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block has Beacon Height greater or equal than 1"))
	}
	//if shardBlock.Header.BeaconHeight == 1 && !shardBlock.Header.BeaconHash.IsPointEqual(&common.Hash{}) {
	//	return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block with Beacon Height 1 have Zero Hash Value"))
	//}
	if shardBlock.Header.BeaconHeight > 1 && shardBlock.Header.BeaconHash.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block with Beacon Height greater or equal than 1 have Non-Zero Hash Value"))
	}
	if shardBlock.Header.TotalTxsFee == nil {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Total Txs Fee have nil value"))
	}
	if len(shardBlock.Header.InstructionMerkleRoot[:]) != common.HashSize {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Instruction Merkle Root in the right format"))
	}
	// body
	if shardBlock.Body.Instructions == nil {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Instruction is not nil"))
	}
	if len(shardBlock.Body.Instructions) != 0 && shardBlock.Header.InstructionMerkleRoot.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Instruction Merkle Root have Non-Zero Hash Value because Instrucstion List is not empty"))
	}
	if shardBlock.Body.CrossTransactions == nil {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Cross Transactions Map is not nil"))
	}
	if len(shardBlock.Body.CrossTransactions) != 0 && shardBlock.Header.CrossTransactionRoot.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Cross Transaction Root have Non-Zero Hash Value because Cross Transaction List is not empty"))
	}
	if shardBlock.Body.Transactions == nil {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Transactions is not nil"))
	}
	if len(shardBlock.Body.Transactions) != 0 && shardBlock.Header.TxRoot.IsEqual(&common.Hash{}) {
		return false, blockchain.NewBlockChainError(blockchain.ShardBlockSanityError, fmt.Errorf("Expect Shard Block Tx Root have Non-Zero Hash Value because Transactions List is not empty"))
	}
	return true, nil
}

func (shardBlock *ShardBlock) UnmarshalJSON(data []byte) error {
	tempShardBlock := &struct {
		ConsensusHeader consensusheader.ConsensusHeader
		Header          ShardHeader
		Body            *json.RawMessage
	}{}
	err := json.Unmarshal(data, &tempShardBlock)
	if err != nil {
		return blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
	}
	shardBlock.ConsensusHeader = tempShardBlock.ConsensusHeader
	blkBody := ShardBody{}
	err = blkBody.UnmarshalJSON(*tempShardBlock.Body)
	if err != nil {
		return blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
	}
	shardBlock.Header = tempShardBlock.Header
	if shardBlock.Body.Transactions == nil {
		shardBlock.Body.Transactions = []metadata.Transaction{}
	}
	if shardBlock.Body.Instructions == nil {
		shardBlock.Body.Instructions = [][]string{}
	}
	if shardBlock.Body.CrossTransactions == nil {
		shardBlock.Body.CrossTransactions = make(map[byte][]blockchain.CrossTransaction)
	}
	if shardBlock.Header.TotalTxsFee == nil {
		shardBlock.Header.TotalTxsFee = make(map[common.Hash]uint64)
	}
	if ok, err := shardBlock.validateSanityData(); !ok || err != nil {
		// panic(string(data) + err.Error())
		return blockchain.NewBlockChainError(blockchain.UnmashallJsonShardBlockError, err)
	}
	shardBlock.Body = blkBody
	return nil
}

// /*
// AddTransaction adds a new transaction into block
// */
// // #1 - tx
func (shardBlock *ShardBlock) AddTransaction(tx metadata.Transaction) error {
	if shardBlock.Body.Transactions == nil {
		return blockchain.NewBlockChainError(blockchain.UnExpectedError, errors.New("not init tx arrays"))
	}
	shardBlock.Body.Transactions = append(shardBlock.Body.Transactions, tx)
	return nil
}

// func (shardBlock *ShardBlock) VerifyBlockReward(blockchain *BlockChain) error {
// 	hasBlockReward := false
// 	txsFee := uint64(0)
// 	for _, tx := range shardBlock.Body.Transactions {
// 		if tx.GetMetadataType() == metadata.ShardBlockReward {
// 			if hasBlockReward {
// 				return errors.New("This block contains more than one coinbase transaction for shard block producer!")
// 			}
// 			hasBlockReward = true
// 		} else {
// 			txsFee += tx.GetTxFee()
// 		}
// 	}
// 	if !hasBlockReward {
// 		return errors.New("This block dont have coinbase tx for shard block producer")
// 	}
// 	numberOfTxs := len(shardBlock.Body.Transactions)
// 	if shardBlock.Body.Transactions[numberOfTxs-1].GetMetadataType() != metadata.ShardBlockReward {
// 		return errors.New("Coinbase transaction must be the last transaction")
// 	}

// 	receivers, values := shardBlock.Body.Transactions[numberOfTxs-1].GetReceivers()
// 	if len(receivers) != 1 {
// 		return errors.New("Wrong receiver")
// 	}
// 	if !common.ByteEqual(receivers[0], shardBlock.Header.ProducerAddress.Pk) {
// 		return errors.New("Wrong receiver")
// 	}
// 	reward := blockchain.getRewardAmount(shardBlock.Header.Height)
// 	reward += txsFee
// 	if reward != values[0] {
// 		return errors.New("Wrong reward value")
// 	}
// 	return nil
// }

// func (block *ShardBlock) getBlockRewardInst(blockHeight uint64) ([]string, error) {
// 	txsFee := uint64(0)

// 	for _, tx := range block.Body.Transactions {
// 		txsFee += tx.GetTxFee()
// 	}
// 	blkRewardInfo := metadata.NewBlockRewardInfo(txsFee, blockHeight)
// 	inst, err := blkRewardInfo.GetStringFormat()
// 	return inst, err
// }

func (block *ShardBlock) AddValidationField(validationData string) error {
	block.ConsensusHeader.ValidationData = validationData
	return nil
}

func (block ShardBlock) GetEpoch() uint64 {
	return block.Header.Epoch
}

func (block ShardBlock) GetProducer() string {
	return block.Header.Producer
}

func (block ShardBlock) GetValidationField() string {
	return block.ConsensusHeader.ValidationData
}

func (block ShardBlock) GetHeight() uint64 {
	return block.Header.Height
}

func (block ShardBlock) GetBeaconHeight() uint64 {
	return block.Header.BeaconHeight
}

func (block ShardBlock) GetRound() int {
	return block.Header.Round
}

func (block ShardBlock) GetRoundKey() string {
	return fmt.Sprint(block.Header.Height, "_", block.Header.Round)
}

func (block ShardBlock) GetInstructions() [][]string {
	return block.Body.Instructions
}

func (block ShardBlock) GetPreviousBlockHash() common.Hash {
	return block.Header.PreviousBlockHash
}

func (block ShardBlock) GetConsensusType() string {
	return block.Header.ConsensusType
}

func (shardBlock *ShardBlock) GetTimeslot() uint64 {
	return shardBlock.ConsensusHeader.TimeSlot
}

func (shardBlock *ShardBlock) GetCreateTimeslot() uint64 {
	return shardBlock.Header.TimeSlot
}

func (shardBlock *ShardBlock) GetBlockTimestamp() int64 {
	return shardBlock.Header.Timestamp
}

func (shardBlock *ShardBlock) GetBlockType() string {
	return "shard"
}

func (shardBlock *ShardBlock) CreateShardToBeaconBlock(bc *blockchain.BlockChain, log common.Logger) *ShardToBeaconBlock {
	if bc.IsTest {
		return &ShardToBeaconBlock{}
	}
	block := ShardToBeaconBlock{}

	block.ConsensusHeader = shardBlock.ConsensusHeader
	block.Header = shardBlock.Header
	blockInstructions := shardBlock.Body.Instructions
	instructions, err := blockchain.CreateShardInstructionsFromTransactionAndInstruction(shardBlock.Body.Transactions, bc, shardBlock.Header.ShardID)
	if err != nil {
		log.Error(err)
		return nil
	}

	block.Instructions = append(instructions, blockInstructions...)
	return &block
}

func (shardBlock *ShardBlock) CreateAllCrossShardBlock(activeShards int, log common.Logger) map[byte]blockinterface.CrossShardBlockInterface {
	allCrossShard := make(map[byte]blockinterface.CrossShardBlockInterface)
	if activeShards == 1 {
		return allCrossShard
	}
	for i := 0; i < activeShards; i++ {
		shardID := common.GetShardIDFromLastByte(byte(i))
		if shardID != shardBlock.Header.ShardID {
			crossShard, err := shardBlock.CreateCrossShardBlock(shardID)
			if crossShard != nil {
				log.Infof("Create CrossShardBlock from Shard %+v to Shard %+v: %+v \n", shardBlock.Header.ShardID, shardID, crossShard)
			}
			if crossShard != nil && err == nil {
				allCrossShard[byte(i)] = crossShard
			}
		}
	}
	return allCrossShard
}

func (shardBlock *ShardBlock) CreateCrossShardBlock(shardID byte) (blockinterface.CrossShardBlockInterface, error) {
	crossShard := &CrossShardBlock{}
	crossOutputCoin, crossCustomTokenPrivacyData := blockchain.GetCrossShardData(shardBlock.Body.Transactions, shardID)
	// Return nothing if nothing to cross
	if len(crossOutputCoin) == 0 && len(crossCustomTokenPrivacyData) == 0 {
		return nil, blockchain.NewBlockChainError(blockchain.CreateCrossShardBlockError, errors.New("No cross Outputcoin, Cross Custom Token, Cross Custom Token Privacy"))
	}
	merklePathShard, merkleShardRoot := blockchain.GetMerklePathCrossShard2(shardBlock.Body.Transactions, shardID)
	if merkleShardRoot != shardBlock.Header.ShardTxRoot {
		return crossShard, blockchain.NewBlockChainError(blockchain.VerifyCrossShardBlockShardTxRootError, fmt.Errorf("Expect Shard Tx Root To be %+v but get %+v", shardBlock.Header.ShardTxRoot, merkleShardRoot))
	}

	crossShard.ConsensusHeader = shardBlock.ConsensusHeader
	crossShard.Header = shardBlock.Header
	crossShard.MerklePathShard = merklePathShard
	crossShard.CrossOutputCoin = crossOutputCoin
	crossShard.CrossTxTokenPrivacyData = crossCustomTokenPrivacyData
	crossShard.ToShardID = shardID
	return crossShard, nil
}

func (shardBlock ShardBlock) GetShardHeader() blockinterface.ShardHeaderInterface {
	return shardBlock.Header
}
func (shardBlock ShardBlock) GetShardBody() blockinterface.ShardBodyInterface {
	return shardBlock.Body
}

func (block ShardBlock) GetBody() blockinterface.BlockBodyInterface {
	return block.Body
}
func (block ShardBlock) GetHeader() blockinterface.BlockHeaderInterface {
	return block.Header
}
