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
	inputfile := emulateUserInput("Tom")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	if choice, err := seller.PresentMovies(); err != nil {
		t.Error("Error while choosing movie", err)
	} else if choice != 2 {
		t.Error("Wrong movie chosen, expected 2, got", choice)
	}

	if err := inputfile.Close(); err != nil {
		log.Fatal(err)
	}
}
