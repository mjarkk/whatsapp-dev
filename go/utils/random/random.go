package random

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"strconv"
)

func Bytes(r *rand.Rand, size int) []byte {
	resp := []byte{}
	for {
		resp = binary.BigEndian.AppendUint64(resp, r.Uint64())
		if len(resp) > size {
			resp = resp[:size]
			break
		}
	}
	return resp
}

func Hex(r *rand.Rand, bytesSize int) string {
	return hex.EncodeToString(Bytes(r, bytesSize))
}

func Numbers(r *rand.Rand, size int) string {
	resp := ""
	for i := 0; i < size; i++ {
		resp += strconv.Itoa(r.Intn(10))
	}
	return resp
}

type RandomValues struct {
	PhoneNumber        string
	PhoneNumberID      string
	GraphToken         string
	AppSecret          string
	WebhookVerifyToken string
}

func GetRandomValuesForSetup(r *rand.Rand) RandomValues {
	return RandomValues{
		PhoneNumber:        "31" + Numbers(r, 8),
		PhoneNumberID:      Numbers(r, 15),
		GraphToken:         base64.StdEncoding.EncodeToString(Bytes(r, 172)),
		AppSecret:          Hex(r, 16),
		WebhookVerifyToken: Hex(r, 16),
	}
}
