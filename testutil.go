package main

import (
	"io/ioutil"
	"os"
)

type TestHandler struct {
	oldstdin  *os.File
	inputfile *os.File
}

func (handler *TestHandler) setUpMockInput(input string) {
	handler.oldstdin = os.Stdin
	handler.inputfile = emulateUserInput(input)
}

func (handler *TestHandler) cleanUp() {
	handler.inputfile.Close()
	os.Remove(handler.inputfile.Name())
	os.Stdin = handler.oldstdin
}

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
