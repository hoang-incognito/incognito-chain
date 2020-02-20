package debugtool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendRequest(query, url string) (*AutoTx, error) {
	var jsonStr = []byte(query)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		data := new(AutoTx)
		err := json.Unmarshal(body, data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func GetTransactionByHash(hash string) (*AutoTx, error) {
	query := fmt.Sprintf(`{  
			"jsonrpc":"1.0",
			"method":"gettransactionbyhash",
			"params":["%s"],
			"id":1
		}`, hash)
	url := `http://51.83.36.184:9334`
	return SendRequest(query, url)
}

func GetProofFromTransactionHash(hash string) (string, error) {
	tx, err := GetTransactionByHash(hash)
	if err != nil {
		return "", err
	}
	return tx.Result.Proof, nil
}
