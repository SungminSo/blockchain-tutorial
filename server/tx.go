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

type Tx struct {
	AddrFrom    string
	Transaction []byte
}

func HandleTx(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Tx

	buff.Write(utils.ExtractPayload(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := blockchain.DeserializeTransaction(txData)
	// TODO: verify tx before include in mempool
	Mempool[hex.EncodeToString(tx.ID)] = tx

	if NodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != NodeAddress && node != payload.AddrFrom {
				SendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(Mempool) >= 2 && len(MiningAddress) > 0 {
		MineTransactions:
			var txs []*blockchain.Transaction

			for id := range Mempool {
				tx := Mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := blockchain.NewCoinbaseTX(MiningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := blockchain.UTXOSet{bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(Mempool, txID)
			}

			for _, node := range KnownNodes {
				if node != NodeAddress {
					SendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(Mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func SendTx(addr string, tnx *blockchain.Transaction) {
	data := Tx{NodeAddress, tnx.Serialize()}
	payload := utils.GobEncode(data)
	request := append(utils.CommandToBytes("tx"), payload...)

	SendData(addr, request)
}