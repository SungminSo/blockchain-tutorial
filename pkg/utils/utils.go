package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
)

const CommandLength = 12

func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func CommandToBytes(command string) []byte {
	var bytes [CommandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func BytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

func ExtractCommand(request []byte) []byte {
	return request[:CommandLength]
}

func ExtractPayload(request []byte) []byte {
	return request[CommandLength:]
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}