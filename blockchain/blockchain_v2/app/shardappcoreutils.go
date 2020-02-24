package app

import (
	"errors"
	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/metadata"
	"github.com/incognitochain/incognito-chain/privacy"
)

func (s *ShardCoreApp) buildWithDrawTransactionResponse(txRequest *metadata.Transaction, blkProducerPrivateKey *privacy.PrivateKey) (metadata.Transaction, error) {
	if (*txRequest).GetMetadataType() != metadata.WithDrawRewardRequestMeta {
		return nil, errors.New("Can not understand this request!")
	}
	requestDetail := (*txRequest).GetMetadata().(*metadata.WithDrawRewardRequest)
	amount, err := s.CreateState.bc.GetCommitteeReward(requestDetail.PaymentAddress.Pk, requestDetail.TokenID)
	if (amount == 0) || (err != nil) {
		return nil, errors.New("Not enough reward")
	}
	responseMeta, err := metadata.NewWithDrawRewardResponse((*txRequest).Hash())
	if err != nil {
		return nil, err
	}
	return s.CreateState.bc.InitTxSalaryByCoinID(
		&requestDetail.PaymentAddress,
		amount,
		blkProducerPrivateKey,
		responseMeta,
		requestDetail.TokenID,
		common.GetShardIDFromLastByte(requestDetail.PaymentAddress.Pk[common.PublicKeySize-1]))
}

func (s ShardCoreApp) buildReturnStakingAmountTx(swapPublicKey string,
	blkProducerPrivateKey *privacy.PrivateKey,
) (metadata.Transaction, error) {
	// addressBytes := blockGenerator.chain.config.UserKeySet.PaymentAddress.Pk
	//shardID := common.GetShardIDFromLastByte(addressBytes[len(addressBytes)-1])
	//publicKey, _ := blockGenerator.chain.config.ConsensusEngine.GetCurrentMiningPublicKey()
	//_, committeeShardID := blockGenerator.chain.BestState.Beacon.GetPubkeyRole(publicKey, 0)
	//
	//fmt.Println("SA: get tx for ", swapPublicKey, GetBestStateShard(committeeShardID).StakingTx, committeeShardID)
	//tx, ok := GetBestStateShard(committeeShardID).StakingTx[swapPublicKey]
	//if !ok {
	//	return nil, NewBlockChainError(GetStakingTransactionError, errors.New("No staking tx in best state"))
	//}
	//var txHash = &common.Hash{}
	//err := (&common.Hash{}).Decode(txHash, tx)
	//if err != nil {
	//	return nil, NewBlockChainError(DecodeHashError, err)
	//}
	//blockHash, index, err := blockGenerator.chain.config.DataBase.GetTransactionIndexById(*txHash)
	//if err != nil {
	//	return nil, NewBlockChainError(GetTransactionFromDatabaseError, err)
	//}
	//shardBlock, _, err := blockGenerator.chain.GetShardBlockByHash(blockHash)
	//if err != nil || shardBlock == nil {
	//	Logger.log.Error("ERROR", err, "NO Transaction in block with hash", blockHash, "and index", index, "contains", shardBlock.Body.Transactions[index])
	//	return nil, NewBlockChainError(FetchShardBlockError, err)
	//}
	//txData := shardBlock.Body.Transactions[index]
	//keyWallet, err := wallet.Base58CheckDeserialize(txData.GetMetadata().(*metadata.StakingMetadata).FunderPaymentAddress)
	//if err != nil {
	//	Logger.log.Error("SA: cannot get payment address", txData.GetMetadata().(*metadata.StakingMetadata), committeeShardID)
	//	return nil, blockchain.NewBlockChainError(blockchain.WalletKeySerializedError, err)
	//}
	//Logger.log.Info("SA: build salary tx", txData.GetMetadata().(*metadata.StakingMetadata).FunderPaymentAddress, committeeShardID)
	//paymentShardID := common.GetShardIDFromLastByte(keyWallet.KeySet.PaymentAddress.Pk[len(keyWallet.KeySet.PaymentAddress.Pk)-1])
	//if paymentShardID != committeeShardID {
	//	return nil, blockchain.NewBlockChainError(blockchain.WrongShardIDError, fmt.Errorf("Staking Payment Address ShardID %+v, Not From Current Shard %+v", paymentShardID, committeeShardID))
	//}
	//returnStakingMeta := metadata.NewReturnStaking(
	//	tx,
	//	keyWallet.KeySet.PaymentAddress,
	//	metadata.ReturnStakingMeta,
	//)
	//returnStakingTx := new(transaction.Tx)
	//err = returnStakingTx.InitTxSalary(
	//	txData.CalculateTxValue(),
	//	&keyWallet.KeySet.PaymentAddress,
	//	blkProducerPrivateKey,
	//	blockGenerator.chain.config.DataBase,
	//	returnStakingMeta,
	//)
	////modify the type of the salary transaction
	//returnStakingTx.Type = common.TxReturnStakingType
	//if err != nil {
	//	return nil, blockchain.NewBlockChainError(blockchain.InitSalaryTransactionError, err)
	//}
	//return returnStakingTx, nil
	return nil, nil
}
