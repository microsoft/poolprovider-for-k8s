package main

import (
	"os"
	"hash"
    "encoding/hex"
	"crypto/sha512"
	"crypto/hmac"
	"log"
)

func ComputeHash(message string) string {
	hashAlgorithm := GetHashAlgorithm()
	if(hashAlgorithm == nil) {
		return "";
	}

	hashAlgorithm.Write([]byte(message))
	
	hashedMessage := hex.EncodeToString(hashAlgorithm.Sum(nil))

	return hashedMessage
}

func ValidateHash(message, inputHmac string) bool {
	hashAlgorithm := GetHashAlgorithm()
	if(hashAlgorithm == nil) {
		return false;
	}

	hashAlgorithm.Write([]byte(message))
	expectedMAC := hashAlgorithm.Sum(nil)

	//return hmac.Equal(inputHMACbytearr, expectedMAC)
	inputHMacArr, err := hex.DecodeString(inputHmac)
	if err != nil {
		log.Fatal(err)
	}
	return hmac.Equal(inputHMacArr, expectedMAC)
}

func GetHashAlgorithm() hash.Hash {
	hmacSecret := os.Getenv("VSTS_SECRET")
	if(hmacSecret == "") {
		return nil
	}

	// Setting up the hash algorithm
	mac := hmac.New(sha512.New, []byte(hmacSecret))
	return mac
}