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
	showID := uint(27)

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
	shows := broker.GetShowsByMovie(&(movs[0]))
	if len(shows) != 7 {
		t.Error("Wrong number of shows, expected 7, got", len(shows))
	}
}

func TestDataBroker_CreateTickets(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	loc, _ := time.LoadLocation("America/Chicago")
	date := time.Date(2019, time.April, 27, 20, 0, 0, 0, loc)
	show := broker.GetShowById(12)

	if serials := broker.CreateTickets(&date, show, 4, 3); len(serials) != 4 {
		t.Error("Wrong number of serials created, expected 4, got", len(serials))
	} else if len(show.Tickets) != 4 {
		t.Error("Wrong number of tickets created, expected 4, got", len(show.Tickets))
	} else { // Cleanup
		broker.DeleteTicketsBySerial(serials)
	}
}

func TestDataBroker_DeleteTicketsBySerial(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	loc, _ := time.LoadLocation("America/Chicago")
	date := time.Date(2019, time.April, 27, 20, 0, 0, 0, loc)
	show := broker.GetShowById(12)
	serials := broker.CreateTickets(&date, show, 4, 3)
	broker.DeleteTicketsBySerial(serials)

	if avail := broker.GetAvailableTickets(&date, show); avail[3] != 10 {
		t.Error("Deleting of tickets unsucessful!")
	}
}

func TestDataBroker_GetMovieTitlesByShows(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	show1 := broker.GetShowById(12)
	show2 := broker.GetShowById(28)
	show3 := broker.GetShowById(59)
	shows := []Show{*show1, *show2, *show3}
	expectedTitles := []string{"Interstellar", "Rogue One: A Star Wars Story", "March of the Penguins"}

	if titles := broker.GetMovieTitlesByShows(shows); len(titles) != 3 {
		t.Error("Received wrong number of titles, expected 3, got", len(titles))
	} else {
		for i, title := range titles {
			if title != expectedTitles[i] {
				t.Error("Show", i+1, "with ID", shows[i].ShowID, "has wrong title, expected", expectedTitles[i],
					"got", title)
			}
		}
	}
}

func TestDataBroker_GetTierPrice(t *testing.T) {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	if price := broker.GetTierPrice(2); price != 15.0 {
		t.Error("Received wrong price, expected 15, got", price)
	}
}
