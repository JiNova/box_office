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
