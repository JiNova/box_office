package main

import (
	"testing"
)

func TestUtil_ReadCmd(t *testing.T) {
	var tester TestHandler
	tester.setUpMockInput("this is ma jam!")
	defer tester.cleanUp()

	if input := readcmd("prompt"); input != "tis is ma jam!" {
		t.Error("Wrong user input received, expected \"this is ma jam\", got", input)
	}
}

func TestUtil_ChooseOption(t *testing.T) {
	var tester TestHandler
	tester.setUpMockInput("7")
	defer tester.cleanUp()

	if input, err := chooseMenuOption("prompt"); err != nil {
		t.Error("Error while receiving user choice, got", err)
	} else if input != 7 {
		t.Error("Wrong choice received, expected 7, got", input)
	}
}
