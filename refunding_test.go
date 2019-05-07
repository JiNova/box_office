package main

import (
	"os"
	"testing"
	"time"
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

func TestRefundHandler_FilterSerialsForValid(t *testing.T) {
	var broker DataBroker
	broker.Init()
	refunder := RefundHandler{&broker}
	defer broker.Close()

	show := broker.GetShowById(19)
	date := time.Now().AddDate(0, 0, 14)
	serials := broker.CreateTickets(&date, show, 4, 4)
	defer broker.DeleteTicketsBySerial(serials)
	serials = append(serials, "aaaaaa", "bcdefghh")

	if serials, err := refunder.FilterSerialsForValid(serials); err != nil {
		t.Error("Error while filtering tickets, got", err)
	} else if len(serials) != 4 {
		t.Error("Filtering returned wrong number of serials, expected 4, got", len(serials))
	}
}
