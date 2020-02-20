package main

import (
	"fmt"

	"github.com/incognitochain/incognito-chain/privacy/debugtool"
)

func testDebugTool() {
	hash := "93c5a85c0f3b99246d4ee3e295e903412adfbed1f108f268ee07fd49ffe8d2f2"
	proof, err := debugtool.GetProofFromTransactionHash(hash)
	if err != nil {
		fmt.Println("Get error")
		fmt.Println(err)
		return
	}
	fmt.Println(proof)
}

func main() {
	testDebugTool()
}
