package main

import (
	"io/ioutil"
	"os"
)

type TestHandler struct {
	oldstdin     *os.File
	testFilename string
}

func (handler *TestHandler) setUpMockInput(input string) {
	handler.oldstdin = os.Stdin
	inputfile := emulateUserInput("04/27/2019")
	handler.testFilename = inputfile.Name()
}

func (handler *TestHandler) cleanUp() {
	os.Remove(handler.testFilename)
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
