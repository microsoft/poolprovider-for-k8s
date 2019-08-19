package main

import (
	"os"
	"hash"

	"crypto/sha512"
	"crypto/hmac"
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
	
	return hmac.Equal([]byte(inputHmac), expectedMAC)
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