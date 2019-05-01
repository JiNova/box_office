package main

import (
	"testing"
	"time"
)

func TestDataBroker_GetAvailableTickets(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	now := time.Date(2019, 4, 28, 14, 0, 0, 0 ,time.Now().Location())
	var show Show
	broker.dbhandler.FillModelById(&show, 75)

	avail := broker.GetAvailableTickets(&now, &show)

	var sum int
	for _,amount := range avail {
		sum += amount
	}

	if sum != 35 {
		t.Error("Expected 40, got", sum)
	}
}
