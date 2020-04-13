package server

import (
	"bytes"
	"encoding/gob"
	"blockchain-tutorial/blockchain"
	"blockchain-tutorial/pkg/utils"
	"log"
)

type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}


func HandleVersion(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Version

	buff.Write(request[utils.CommandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		SendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		SendVersion(payload.AddrFrom, bc)
	}

	if !NodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

func SendVersion(addr string, bc *blockchain.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(Version{nodeVersion, bestHeight, NodeAddress})

	request := append(utils.CommandToBytes("version"), payload...)
	SendData(addr, request)
}