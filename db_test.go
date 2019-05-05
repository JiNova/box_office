package main

import (
	"strconv"
	"testing"
)

func TestData_FillModelById(t *testing.T) {
	var data DBHandler
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
	var data DBHandler
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

func TestData_LoadAssociationsSingle(t *testing.T) {
	var data DBHandler
	data.Init()
	defer data.Close()

	var mov Movie
	expectedShowNum := 5
	movieID := 5
	data.FillModelById(&mov, movieID)
	mov.Shows = []Show{}
	err := data.LoadAssociations(&mov, "Shows")

	if err != nil {
		t.Error("Could not load shows", err)
	} else if mov.Shows == nil {
		t.Error("Show data not loaded successfully")
	} else if len(mov.Shows) != expectedShowNum {
		t.Error("Expected ", expectedShowNum, "got", len(mov.Shows))
	}

	for i, show := range mov.Shows {
		if show.MovieID != uint(movieID) {
			t.Error("Show", i, "belongs to another movie!")
		}
	}
}

func TestData_LoadAssociationsMulti(t *testing.T) {
	var data DBHandler
	data.Init()
	defer data.Close()

	var movs []Movie
	expectedShowNum := 5
	verifyMovIndex := 4
	data.FillModels(&movs)
	for _, mov := range movs {
		mov.Shows = []Show{}
	}

	err := data.LoadAssociations(&movs, "Shows")

	if err != nil {
		t.Error("Could not load shows", err)
	} else if movs[verifyMovIndex].Shows == nil {
		t.Error("Show data not loaded successfully")
	} else if len(movs[verifyMovIndex].Shows) != expectedShowNum {
		t.Error("Expected ", expectedShowNum, "got", len(movs[verifyMovIndex].Shows))
	}

	for i, show := range movs[verifyMovIndex].Shows {
		if show.MovieID != uint(verifyMovIndex+1) {
			t.Error("Show", i, "belongs to another movie!")
		}
	}
}

func TestData_CountEntries(t *testing.T) {
	var data DBHandler
	data.Init()
	defer data.Close()

	if count, err := data.CountEntries(&Day{}); count != 7 {
		t.Error("Expected 7 days, got", count)
	} else if err != nil {
		t.Error("Could not count days, got", err)
	}
}

func TestData_QueryModelSingle(t *testing.T) {
	var data DBHandler
	data.Init()
	defer data.Close()

	var mov Movie
	expectedID := uint(5)

	if err := data.QueryModel(&mov, "title = ?", "Rogue One: A Star Wars Story"); err != nil {
		t.Error("Could not query, got", err)
	} else if expectedID != mov.MovieID {
		t.Error("Expected movie_id", expectedID, "got", mov.MovieID)
	}
}

func TestData_QueryModelMulti(t *testing.T) {
	var data DBHandler
	data.Init()
	defer data.Close()

	var shows []Show
	expectedID := uint(5)

	if err := data.QueryModel(&shows, "movie_id = ?", strconv.Itoa(int(expectedID))); err != nil {
		t.Error("Could not query, got", err)
	}

	for _, show := range shows {
		if expectedID != show.MovieID {
			t.Error("Expected movie_id", expectedID, "got", show.MovieID)
		}
	}
}
