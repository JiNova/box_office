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

func TestDataBroker_GetTicketDatesBySerials(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	serials := []string{"fd7d3b7b", "cb45c26a", "abcdefgh"}

	if dates := broker.GetTicketDatesBySerials(serials); dates == nil {
		t.Error("Could not load dates")
	} else if len(dates) != 2 {
		t.Error("Wronger number of dates, expected 2, got", dates)
	}
}

func TestDataBroker_GetTicketCountByDay(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	loc, _ := time.LoadLocation("America/Chicago")
	day := time.Date(2019, time.April, 28, 0, 0, 0, 0, loc)

	if count := broker.GetTicketCountByDay(&day); count != 19 {
		t.Error("Wrong number of tickets, expected 19, got", count)
	}
}

func TestDataBroker_GetSoldUnsoldTicketsByShow(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	loc, _ := time.LoadLocation("America/Chicago")
	date := time.Date(2019, time.April, 28, 14, 0, 0, 0, loc)
	showID := 27

	if sold, vacant := broker.GetSoldVacantTicketsByShow(&date, showID); sold != 9 || vacant != 31 {
		t.Error("Wrong ticket constellation, expected 9 sold, 31 vacant, got",
			sold, "sold,", vacant, "vacant")
	}
}

func TestDataBroker_GetAllMovies(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	movs := broker.GetAllMovies()
	if len(movs) != 13 {
		t.Error("Wrong number of movies, expected 13, got", len(movs))
	}
}

func TestDataBroker_GetShowsByMovie(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	movs := broker.GetAllMovies()
	if shows := broker.GetShowsByMovie(&(movs[0])); len(shows) != 7 {
		t.Error("Wrong number of shows, expected 7, got", len(shows))
	}
}
