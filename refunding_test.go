package main

import (
	"os"
	"testing"
)

func TestRefundHandler_GetSerialsFromUser(t *testing.T) {
	var broker DataBroker
	broker.Init()
	refunder := RefundHandler{&broker}

	oldStdin := os.Stdin
	inputfile := emulateUserInput("aaabbb sdafdas asdfsdff adsffssdf")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	if serials, err := refunder.GetSerialsFromUser(); err != nil {
		t.Error("Error when trying to get serials, got", err)
	} else if len(serials) != 4 {
		t.Error("Wrong number of serials, expected 4, got", len(serials))
	}
}
