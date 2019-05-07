package main

import (
	"os"
	"testing"
	"time"
)

func TestReportHandler_GetDateFromUser(t *testing.T) {
	var broker DataBroker
	broker.Init()

	reporter := ReportHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("04/27/2019")
	expectedDate := time.Date(2019, time.April, 27, 0, 0, 0, 0, time.UTC)

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	if date := reporter.GetDateFromUser(); !date.Equal(expectedDate) {
		t.Error("Wrong date returned, expected", expectedDate, "got", date)
	}
}

func TestReportHandler_GetTimeFromUser(t *testing.T) {
	var broker DataBroker
	broker.Init()

	reporter := ReportHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("7 pm")
	expectedTime, _ := time.Parse("3 pm", "7 pm")

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	if time := reporter.GetTimeFromUser(); !time.Equal(expectedTime) {
		t.Error("Wrong date returned, expected", expectedTime, "got", time)
	}
}

func TestReportHandler_GetSpecificShowFromUser(t *testing.T) {
	var broker DataBroker
	broker.Init()

	reporter := ReportHandler{&broker}
	oldStdin := os.Stdin
	inputfile := emulateUserInput("4")
	loc, _ := time.LoadLocation("America/Chicago")
	date := time.Date(2019, time.April, 28, 20, 0, 0, 0, loc)

	defer os.Remove(inputfile.Name())      // clean up
	defer func() { os.Stdin = oldStdin }() // Restore stdin
	defer broker.Close()

	if show, title, err := reporter.GetSpecificShowFromUser(&date); err != nil {
		t.Error("Error while getting show from user, got", err)
	} else if title != "Rogue One: A Star Wars Story" {
		t.Error("Wrong movie title, expected \"Rogue One: A Star Wars Story\", got", title)
	} else if show.ShowID != 28 {
		t.Error("Shhow with wrong ID received, expected 28, got", show.ShowID)
	}
}
