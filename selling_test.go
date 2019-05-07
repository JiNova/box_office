package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func emulateUserInput(input string) (inputfile *os.File) {
	content := []byte(input)

	inputfile, err := ioutil.TempFile("", "input")
	if err != nil {
		panic(err)
	}

	if _, err := inputfile.Write(content); err != nil {
		panic(err)
	}

	if _, err := inputfile.Seek(0, 0); err != nil {
		panic(err)
	}

	os.Stdin = inputfile
	return
}

func TestSellHandler_PresentMovies(t *testing.T) {
	var broker DataBroker
	broker.Init()

	seller := SellHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("2")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	if choice, err := seller.PresentMovies(); err != nil {
		t.Error("Error while choosing movie", err)
	} else if choice.MovieID != 2 {
		t.Error("Wrong movie chosen, expected 2, got", choice.MovieID)
	}

	if err := inputfile.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestSellHandler_ChooseShow(t *testing.T) {
	var broker DataBroker
	broker.Init()

	seller := SellHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("2")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	movie := broker.GetMovieById(2)
	shows := broker.GetShowsByMovie(movie)
	if choice, time, err := seller.ChooseShow(movie, shows); err != nil {
		t.Error("Error while choosing show", err)
	} else if choice.ShowID != 8 {
		t.Error("Wrong ShowID, expected 8, got", choice.ShowID)
	} else if time.Weekday().String() != "Tuesday" || time.Hour() != 14 {
		t.Error("Wrong playtime, expected Tuesday 2pm, got", time.Format("Monday 3pm"))
	}
}
