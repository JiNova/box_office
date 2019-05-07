package main

import (
	"os"
	"testing"
	"time"
)

func TestSellHandler_PresentMovies(t *testing.T) {
	var broker DataBroker
	broker.Init()
	seller := SellHandler{&broker}
	defer broker.Close()

	var tester TestHandler
	tester.setUpMockInput("2")
	defer tester.cleanUp()

	if choice, err := seller.PresentMovies(); err != nil {
		t.Error("Error while choosing movie", err)
	} else if choice.MovieID != 2 {
		t.Error("Wrong movie chosen, expected 2, got", choice.MovieID)
	}
}

func TestSellHandler_ChooseShow(t *testing.T) {
	var broker DataBroker
	broker.Init()

	seller := SellHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("3")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	movie := broker.GetMovieById(2)
	shows := broker.GetShowsByMovie(movie)
	if choice, time, err := seller.ChooseShow(movie, shows); err != nil {
		t.Error("Error while choosing show", err)
	} else if choice.ShowID != 10 {
		t.Error("Wrong ShowID, expected 8, got", choice.ShowID)
	} else if time.Weekday().String() != "Friday" || time.Hour() != 20 {
		t.Error("Wrong playtime, expected Friday 2pm, got", time.Format("Monday 3pm"))
	}
}

func TestSellHandler_ChooseTier(t *testing.T) {
	var broker DataBroker
	broker.Init()

	seller := SellHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("4")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	show := broker.GetShowById(11)
	loc, _ := time.LoadLocation("America/Chicago")
	date := time.Date(2019, time.May, 2, 14, 0, 0, 0, loc)
	avail := broker.GetAvailableTickets(&date, show)
	if tier, err := seller.ChooseTier(show, avail); err != nil {
		t.Error("Error while chosing tier,", err)
	} else if tier != 4 {
		t.Error("Resolved wrong tier, expected 4, got", tier)
	}
}

func TestSellHanlder_SellTickets(t *testing.T) {
	var broker DataBroker
	broker.Init()

	seller := SellHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("6")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	show := broker.GetShowById(11)
	loc, _ := time.LoadLocation("America/Chicago")
	date := time.Date(2019, time.May, 2, 14, 0, 0, 0, loc)
	tier := 1
	available := 10

	if serials, err := seller.SellTickets(show, &date, tier, available); err != nil {
		t.Error("Error while buying tickets,", err)
	} else if len(serials) != 6 {
		t.Error("Expected 6 ticket serials, got", len(serials))
	} else {
		broker.DeleteTicketsBySerial(serials)
	}
}
