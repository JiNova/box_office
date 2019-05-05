package main

import (
	"testing"
	"time"
)

func TestDataBroker_GetAvailableTickets(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	now := time.Date(2019, 4, 28, 14, 0, 0, 0, time.Now().Location())
	var show Show
	broker.dbhandler.FillModelById(&show, 75)

	avail := broker.GetAvailableTickets(&now, &show)

	var sum int
	for _, amount := range avail {
		sum += amount
	}

	if sum != 35 {
		t.Error("Expected 40, got", sum)
	}
}

func TestDataBroker_GetDayIdByName(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	if dayID := broker.GetDayIdByName("Friday"); dayID != 5 {
		t.Error("Expected id 5, got", dayID)
	}
}

func TestDataBroker_GetShowsByPlaytime(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	if shows := broker.GetShowsByPlaytime("Tuesday", 14); shows == nil {
		t.Error("Could not load shows")
	} else if len(shows) != 5 {
		t.Error("Wrong number of shows, expected 5, got", len(shows))
	}
}

func TestDataBroker_GetTicketsBySerials(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	serials := []string{"fd7d3b7b", "cb45c26a", "abcdefgh"}

	if tickets := broker.GetTicketsBySerials(serials); tickets == nil {
		t.Error("Could not load tickets")
	} else if len(tickets) != 2 {
		t.Error("Wronger number of tickets, expected 2, got", len(tickets))
	}
}
