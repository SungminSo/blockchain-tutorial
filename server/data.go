package server

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"blockchain-tutorial/blockchain"
	"blockchain-tutorial/pkg/utils"
	"io"
	"log"
	"net"
)

type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

func HandleGetData(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload GetData

	buff.Write(utils.ExtractPayload(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			log.Panic(err)
		}

		SendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := Mempool[txID]
		SendTx(payload.AddrFrom, &tx)
	}
}

func SendGetData(address, kind string, id []byte) {
	payload := utils.GobEncode(GetData{NodeAddress, kind, id})
	request := append(utils.CommandToBytes("getdata"), payload...)

	SendData(address, request)
}

func SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range KnownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		KnownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}