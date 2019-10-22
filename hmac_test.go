package main

import (
	"testing"
	"os"
)

func TestValidateAndComputeHash(t *testing.T) {

	os.Setenv("VSTS_SECRET", "sharedsecret1234")
	str := ComputeHash("teststring")
	
	check := ValidateHash("teststring",str)
	if (check == false){
		t.Errorf("Hmac validation failed")
	}

}
func TestValidateHashShouldReturnFalse(t *testing.T) {

	os.Setenv("VSTS_SECRET", "sharedsecret12345")
	str := ComputeHash("teststring")

	// changing secret value
	os.Setenv("VSTS_SECRET", "sharedsecret1234")
	
	check := ValidateHash("teststring",str)
	if (check == true){
		t.Errorf("Hmac validation failed")
	}

}