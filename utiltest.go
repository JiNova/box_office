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
