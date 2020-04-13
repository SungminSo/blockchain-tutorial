package server

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"blockchain-tutorial/blockchain"
	"blockchain-tutorial/pkg/utils"
	"log"
)

type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

func HandleInv(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Inv

	buff.Write(utils.ExtractPayload(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		BlocksInTransit = payload.Items

		blockHash := payload.Items[0]
		SendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range BlocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		BlocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]
		if Mempool[hex.EncodeToString(txID)].ID == nil {
			SendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func SendInv(address, kind string, items [][]byte) {
	inventory := Inv{NodeAddress, kind, items}
	payload := utils.GobEncode(inventory)
	request := append(utils.CommandToBytes("inv"), payload...)

	SendData(address, request)
}