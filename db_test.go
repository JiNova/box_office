package main

import (
	"testing"
)

func TestData_FillModelById(t *testing.T) {
	var data Data
	data.Init()
	defer data.Close()

	var mov Movie
	err := data.FillModelById(&mov, 5)
	expectedTitle := "Rogue One: A Star Wars Story"
	expectedLen := uint(133)

	if err != nil {
		t.Error("Movie could not be loaded", err)
	} else if mov.Title != expectedTitle || mov.Length != expectedLen {
		t.Error("Expected", expectedTitle, "and", expectedLen,
				"got", mov.Title, "and", mov.Length)
	}
}

func TestData_FillModels(t *testing.T) {
	var data Data
	data.Init()
	defer data.Close()

	var movs []Movie
	err := data.FillModels(&movs)

	if movs == nil {
		t.Error("Movie models not loaded successfuly")
	} else if err != nil {
		t.Error("Could not load movies", err)
	} else if len(movs) != 13 {
		t.Error("Should be 13 movies, but we have", len(movs))
	}
}

func TestData_LoadAssociations(t *testing.T) {
	var data Data
	data.Init()
	defer data.Close()

	var mov Movie
	expectedShowNum := 5
	data.FillModelById(&mov, 5)
	err := data.LoadAssociations(&mov, "Shows")

	if err != nil {
		t.Error("Could not load shows")
	} else if mov.Shows == nil {
		t.Error("Show data not loaded successfully")
	} else if len(mov.Shows) != expectedShowNum {
		t.Error("Expected ", expectedShowNum, "got", len(mov.Shows))
	}
}
