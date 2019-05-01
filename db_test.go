package main

import (
	"testing"
)

func TestData_LoadModels(t *testing.T) {
	var data Data

	data.Init()
	defer data.Close()

	var movs []Movie
	err := data.LoadModels(&movs)

	if movs == nil {
		t.Error("Movie models not loaded successfuly")
	} else if err != nil {
		t.Error("Could not load movies", err)
	} else if len(movs) != 13 {
		t.Error("Should be 13 movies, but we have", len(movs))
	}
}

func TestData_LoadAssociations(t *testing.T) {

}
