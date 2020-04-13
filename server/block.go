package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"blockchain-tutorial/blockchain"
	"blockchain-tutorial/pkg/utils"
	"log"
)

type Block struct {
	AddrFrom string
	Block    []byte
}

type GetBlocks struct {
	AddrFrom string
}

func HandleBlock(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Block

	buff.Write(utils.ExtractPayload(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)

	fmt.Println("Received a new block!")
	// TODO: verify block's validity
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(BlocksInTransit) > 0 {
		blockHash := BlocksInTransit[0]
		SendGetData(payload.AddrFrom, "block", blockHash)
		BlocksInTransit = BlocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{bc}
		// TODO: use UTXOSet.Update(block) instead of UTXOSet.Reindex()
		UTXOSet.Reindex()
	}
}

func HandleGetBlocks(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(utils.ExtractPayload(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	SendInv(payload. AddrFrom, "block", blocks)
}

func RequestBlocks() {
	for _, node := range KnownNodes {
		SendGetBlocks(node)
	}
}

func SendBlock(addr string, b *blockchain.Block) {
	data := Block{NodeAddress, b.Serialize()}
	payload := utils.GobEncode(data)
	request := append(utils.CommandToBytes("block"), payload...)

	SendData(addr, request)
}

func SendGetBlocks(address string) {
	payload := utils.GobEncode(GetBlocks{NodeAddress})
	request := append(utils.CommandToBytes("getblocks"), payload...)

	SendData(address, request)
}