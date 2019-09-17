package main

import (
	"os"
	"hash"

	"crypto/sha512"
	"crypto/hmac"
	"encoding/hex"
	"log"
)

func ComputeHash(message string) [] byte {
	hashAlgorithm := GetHashAlgorithm()
	if(hashAlgorithm == nil) {
		return nil;
	}

	hashAlgorithm.Write([]byte(message))
	hashedMessage := hashAlgorithm.Sum(nil)

	return hashedMessage
}

func ValidateHash(message, inputHmac string) bool {
	hashAlgorithm := GetHashAlgorithm()
	if(hashAlgorithm == nil) {
		return true;
	}

	hashAlgorithm.Write([]byte(message))
	expectedMAC := hashAlgorithm.Sum(nil)
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