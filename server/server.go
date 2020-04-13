package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"blockchain-tutorial/blockchain"
	"blockchain-tutorial/pkg/utils"
	"io/ioutil"
	"log"
	"net"
)

const (
	protocol = "tcp"
	nodeVersion = 1
)

type Addr struct {
	AddrList []string
}

var (
	NodeAddress 	string
	MiningAddress 	string
	KnownNodes = []string{"localhost:3000"}
	BlocksInTransit = [][]byte{}
	Mempool = make(map[string]blockchain.Transaction)
)

func StartServer(nodeID, minerAddress string) {
	NodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	MiningAddress = minerAddress
	ln, err := net.Listen(protocol, NodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := blockchain.NewBlockchain(nodeID)

	if NodeAddress != KnownNodes[0] {
		SendVersion(KnownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	command := utils.BytesToCommand(request[:utils.CommandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
	HandleAddr(request)
	case "block":
	HandleBlock(request, bc)
	case "inv":
	HandleInv(request, bc)
	case "getblocks":
	HandleGetBlocks(request, bc)
	case "getdata":
	HandleGetData(request, bc)
	case "tx":
	HandleTx(request, bc)
	case "version":
	HandleVersion(request, bc)
	default:
	fmt.Println("Unknown command!")
	}

	conn.Close()
}

func HandleAddr(request []byte) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(utils.ExtractPayload(request))
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(KnownNodes))
	RequestBlocks()
}

func SendAddr(address string) {
	nodes := Addr{KnownNodes}
	nodes.AddrList = append(nodes.AddrList, NodeAddress)
	payload := utils.GobEncode(nodes)
	request := append(utils.CommandToBytes("addr"), payload...)

	SendData(address, request)
}

func NodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}

	return false
}