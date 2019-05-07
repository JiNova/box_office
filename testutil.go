package main

import (
	"io/ioutil"
	"os"
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
