package phonenumber

import (
	"encoding/base64"
	"errors"
	"math/rand"

	"github.com/dongri/phonenumber"
)

type ParsedPhoneNumber struct {
	Original          string
	Parsed            string
	WhatsappMessageID string
}

func Parse(input string, localAllowed bool) (*ParsedPhoneNumber, error) {
	if len(input) < 6 {
		return nil, errors.New("phone number too short")
	}
	if input[0] == '0' {
		if !localAllowed {
			return nil, errors.New("local phone number not allowed")
		}
	} else if input[0] != '+' {
		input = "+" + input
	}
	number := phonenumber.Parse(input, "NL")

	return &ParsedPhoneNumber{
		Original:          input,
		Parsed:            number,
		WhatsappMessageID: CreateWhatsappID(number),
	}, nil
}

func CreateWhatsappID(phoneNumber string) string {
	idRandomBytes := []byte{}
	for i := 0; i < 24; i++ {
		idRandomBytes = append(idRandomBytes, byte(rand.Intn(256)))
	}

	// wamid.HBgLMzE2MTExMzg3NTIVAgARGBIzRTQ0RDBCODQyMDFDQjkzQUYA
	id := []byte{0x1C, 0x18, 0x0B, 0x33}
	id = append(id, []byte(phoneNumber)...)
	id = append(id, idRandomBytes...)
	id = append(id, 0)
	return "wamid." + base64.StdEncoding.EncodeToString(id)
}
