package main

import (
	"encoding/base64"
	"fmt"

	"github.com/incognitochain/incognito-chain/common"
	"github.com/incognitochain/incognito-chain/privacy"
	zkp "github.com/incognitochain/incognito-chain/privacy/zeroknowledge"

	"github.com/incognitochain/incognito-chain/privacy/debugtool"
)

var _ = func() (_ struct{}) {
	fmt.Println("This runs before init()!")
	privacy.Logger.Init(common.NewBackend(nil).Logger("test", true))
	return
}()

func debugBP() {
	p := privacy.RandomPoint()
	fmt.Println(p)
	fmt.Println(privacy.Logger)
	// transactionHash := "8dabc7e212699c4e3272cc61114b9262b5dca09a7bca5521a88c679d6e026da6"
	// transactionHash := "2703b69b3845bf49c25863f0c1391a5ab7f8c53d0b439e6f0daece5a7bf045af"
	transactionHash := "93c5a85c0f3b99246d4ee3e295e903412adfbed1f108f268ee07fd49ffe8d2f2"
	// transactionHash := "7d4d44e2d75db45a81d6bb35b2ccc7f22feda884efb123b273826c2d19f8b443"
	proofBase64, err_network := debugtool.GetProofFromTransactionHash(transactionHash)
	if err_network != nil {
		fmt.Println("Get error")
		fmt.Println(err_network)
		return
	}
	proof, err_proof := base64.StdEncoding.DecodeString(proofBase64)
	if err_proof != nil {
		fmt.Println("Get error")
		fmt.Println(err_proof)
		return
	}
	//
	paymentProof := new(zkp.PaymentProof)
	paymentProof.SetBytes(proof)
	agg := paymentProof.GetAggregatedRangeProof()
	// fmt.Println(agg.ValidateSanity())
	v, _ := agg.Verify()
	fmt.Println("=====================")
	fmt.Println("Done aggregate test")
	fmt.Println(v)
}

func main() {
	// testDebugTool()
	debugBP()
}
